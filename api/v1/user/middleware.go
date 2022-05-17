package userApi

import (
	"net/http"
	"os"

	"github.com/bipen2001/go-user-assignment-go/internal/entity"
	"github.com/bipen2001/go-user-assignment-go/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func Authenticate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var jwt_sec = []byte(os.Getenv("JWT_SECRET"))

		token, err := req.Cookie("token")

		if err != nil {
			utils.JsonResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Status: http.StatusUnauthorized, ErrorMessage: UNAUTHORIZED})

			return
		}
		tknStr := token.Value

		claims := &entity.Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwt_sec, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.JsonResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Status: http.StatusUnauthorized, ErrorMessage: UNAUTHORIZED})

				return
			}
			utils.JsonResponse(w, http.StatusBadRequest, utils.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: err.Error()})

			return
		}
		if !tkn.Valid {
			utils.JsonResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Status: http.StatusUnauthorized, ErrorMessage: UNAUTHORIZED})

			return
		}
		if req.Method == "PATCH" || req.Method == "DELETE" {
			if claims.Id != mux.Vars(req)["id"] {
				utils.JsonResponse(w, http.StatusForbidden, utils.ErrorResponse{Status: http.StatusForbidden, ErrorMessage: FORBIDDEN})

				return
			}
		}

		f.ServeHTTP(w, req)

	}
}
