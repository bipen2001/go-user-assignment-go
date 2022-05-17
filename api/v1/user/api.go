package userApi

import (
	"fmt"
	"net/http"

	"github.com/bipen2001/go-user-assignment-go/internal/entity"
	"github.com/bipen2001/go-user-assignment-go/internal/service/user/model"
	"github.com/bipen2001/go-user-assignment-go/utils"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
)

type resource struct {
	userService model.Service
}

var validate = validator.New()

func RegisterHandlers(r *mux.Router, service model.Service) {

	res := resource{service}

	r.HandleFunc("/login", res.login).Methods("POST")
	r.HandleFunc("/signup", res.CreateUser).Methods("POST")

	r.HandleFunc("/user", Authenticate(res.GetUsers)).Methods("GET")
	r.HandleFunc("/user/{id}", Authenticate(res.GetUserById)).Methods("GET")
	r.HandleFunc("/user/{id}", Authenticate(res.UpdateUser)).Methods("PATCH")

	r.HandleFunc("/user/{id}", Authenticate(res.DeleteUser)).Methods("DELETE")

}

func (r *resource) login(w http.ResponseWriter, req *http.Request) {

	var cred entity.Creds

	err := utils.SanitizeRequest(req, &cred)

	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: "could not parse request body"})
		return
	}

	user, err := r.userService.Get(req.Context(), entity.QueryParams{
		Email: cred.Email,
	}, true)

	if err != nil {
		fmt.Println(err)
		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: "User with that email do not exist"})
		return
	}

	if len(user) > 1 {
		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: "Multiple User with this email exists"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user[0].Password), []byte(cred.Password))

	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Status: http.StatusUnauthorized, ErrorMessage: "Unauthorized!! Enter Correct Email or Password"})

		return
	}

	token, err := utils.CreateJWT(user[0])
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, utils.ErrorResponse{Status: http.StatusInternalServerError, ErrorMessage: err.Error()})

		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   token.Token,
			Expires: token.Expires,
		})
	user[0].Password = ""
	utils.JsonResponse(w, http.StatusOK, user)

}

func (r *resource) GetUserById(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	if id == "" {
		utils.JsonResponse(w, http.StatusNotFound, utils.ErrorResponse{Status: http.StatusNotFound, ErrorMessage: "Please Provide a valid id"})

		return
	}
	user, err := r.userService.GetById(req.Context(), id)
	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, utils.ErrorResponse{Status: http.StatusNotFound, ErrorMessage: err.Error()})

		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func (r *resource) GetUsers(w http.ResponseWriter, req *http.Request) {

	var queryParams entity.QueryParams

	queryParams.FirstName = req.URL.Query().Get("firstName")
	queryParams.LastName = req.URL.Query().Get("lastName")
	queryParams.Email = req.URL.Query().Get("email")
	queryParams.Archived = req.URL.Query().Get("archived")
	queryParams.Sort = req.URL.Query().Get("sort")
	queryParams.Order = req.URL.Query().Get("order")

	user, err := r.userService.Get(req.Context(), queryParams, false)
	if err != nil {

		utils.JsonResponse(w, http.StatusNotFound, utils.ErrorResponse{Status: http.StatusNotFound, ErrorMessage: err.Error()})

		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func (r *resource) CreateUser(w http.ResponseWriter, req *http.Request) {

	var user *entity.User

	err := utils.SanitizeRequest(req, &user)

	if err != nil {

		utils.JsonResponse(
			w,
			http.StatusBadRequest,
			utils.ErrorResponse{
				Status:       http.StatusBadRequest,
				ErrorMessage: "Cannot parse request body",
			},
		)
		return

	}

	err = validate.Struct(user)
	if err != nil {

		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: err.Error()})

		return
	}

	resUser, err := r.userService.Create(req.Context(), *user)

	if err != nil {

		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: err.Error()})

		return
	}
	token, err := utils.CreateJWT(*resUser)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, utils.ErrorResponse{Status: http.StatusInternalServerError, ErrorMessage: err.Error()})

		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   token.Token,
			Expires: token.Expires,
		})

	utils.JsonResponse(w, http.StatusOK, resUser)
}

func (r *resource) UpdateUser(w http.ResponseWriter, req *http.Request) {

	var usr entity.UpdateUser
	id := mux.Vars(req)["id"]
	if id == "" {

		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: "Please provide a valid id"})

		return
	}
	err := utils.SanitizeRequest(req, &usr)

	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: "Please provide a valid user" + err.Error()})

		return
	}
	err = validate.Struct(usr)
	if err != nil {

		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: err.Error()})

		return
	}

	user, err := r.userService.Update(req.Context(), id, usr)

	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: err.Error()})

		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func (r *resource) DeleteUser(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	if id == "" {
		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: "Please provide a valid id"})

		return
	}

	err := r.userService.Delete(req.Context(), id)

	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: "Could not find the user in db" + err.Error()})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
