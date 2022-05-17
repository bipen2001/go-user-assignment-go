package userApi

import (
	"net/http"
	"os"

	"github.com/bipen2001/go-user-assignment-go/internal/entity"
	"github.com/bipen2001/go-user-assignment-go/utils"
	"github.com/dgrijalva/jwt-go"
)

var jwt_sec = []byte(os.Getenv("JWT_SECRET"))

func Authenticate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("token")

		if err != nil {
			utils.JsonResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Status: http.StatusUnauthorized, ErrorMessage: "Unauthorized"})

			return
		}
		tknStr := token.Value

		claims := &entity.Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwt_sec, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.JsonResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Status: http.StatusUnauthorized, ErrorMessage: "Unauthorized"})

				return
			}
			utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: err.Error()})

			return
		}
		if !tkn.Valid {
			utils.JsonResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Status: http.StatusUnauthorized, ErrorMessage: "Unauthorized"})

			return
		}

		f.ServeHTTP(w, req)

	}
}
