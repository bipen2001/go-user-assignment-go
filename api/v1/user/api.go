package user

import (
	"net/http"

	"github.com/bipen2001/go-user-assignment-go/internal/entity"
	"github.com/bipen2001/go-user-assignment-go/internal/service/user/model"
	"github.com/bipen2001/go-user-assignment-go/utils"
	"github.com/gorilla/mux"
)

type resource struct {
	userService model.Service
}

func RegisterHandlers(r *mux.Router, service model.Service) {

	res := resource{service}

	r.HandleFunc("/user", res.GetUsers).Methods("GET")
	r.HandleFunc("/user/{id}", res.GetUserById).Methods("GET")
	r.HandleFunc("/user/{id}", res.CreateUser).Methods("POST")
	r.HandleFunc("/user/{id}", res.UpdateUser).Methods("PATCH")

	r.HandleFunc("/user/{id}", res.DeleteUser).Methods("DELETE")

}

func (r *resource) GetUserById(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "Please Provide a Id", http.StatusBadRequest)
		return
	}
	user, err := r.userService.GetById(req.Context(), id)
	if err != nil {
		http.Error(w, "Could not find the user in db", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func (r *resource) GetUsers(w http.ResponseWriter, req *http.Request) {

	user, err := r.userService.Get(req.Context())
	if err != nil {
		http.Error(w, "Could not find the user in db", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func (r *resource) CreateUser(w http.ResponseWriter, req *http.Request) {

	var user *entity.User
	err := utils.SanitizeRequest(req, &user)
	if err != nil {
		http.Error(w, "Please Provide Proper User", http.StatusBadRequest)
		return
	}

	user, err = r.userService.Create(req.Context(), *user)

	if err != nil {
		http.Error(w, "Could not Create the user", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func (r *resource) UpdateUser(w http.ResponseWriter, req *http.Request) {

	var user *entity.User
	err := utils.SanitizeRequest(req, &user)
	if err != nil {
		http.Error(w, "Please Provide a valid user", http.StatusBadRequest)
		return
	}

	user, err = r.userService.Update(req.Context(), *user)

	if err != nil {
		http.Error(w, "Could not find the user in db", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func (r *resource) DeleteUser(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "Please Provide a Id", http.StatusBadRequest)
		return
	}

	err := r.userService.Delete(req.Context(), id)

	if err != nil {
		http.Error(w, "Could not find the user in db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
