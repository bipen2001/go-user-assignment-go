package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bipen2001/go-user-assignment-go/internal/entity"
	"github.com/dgrijalva/jwt-go"
)

type ErrorResponse struct {
	Status       int
	ErrorMessage string
}

func JsonResponse(w http.ResponseWriter, status int, resp interface{}) {

	response, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "127.0.0.1:5500")
	// w.Header().Set("Access-Control-Allow-Credentials","true")

	w.WriteHeader(status)
	w.Write(response)
}

func SanitizeRequest(r *http.Request, req interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(req); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func CreateJWT(user entity.Response) (*entity.JwtToken, error) {

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &entity.Claims{
		Email: user.Email,
		Id:    user.Id,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &entity.JwtToken{
		Expires: expirationTime,
		Token:   tokenString,
	}, nil

}
