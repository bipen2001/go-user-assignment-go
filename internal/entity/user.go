package entity

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Id        string     `json:"id" validate:"omitempty,uuid"`
	FirstName string     `json:"firstName" validate:"required" `
	LastName  string     `json:"lastName" validate:"required"`
	Email     *string    `json:"email" validate:"required,email"`
	Dob       *time.Time `json:"dob" validate:"required"`
	Password  string     `json:"password,omitempty" validate:"required"`
}
type UpdateUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
type Response struct {
	Id        string     `json:"id"`
	FirstName *string    `json:"firstName,omitempty"`
	LastName  *string    `json:"lastName,omitempty"`
	Email     string     `json:"email,omitempty"`
	Dob       *time.Time `json:"dob,omitempty"`
	Age       int        `json:"age,omitempty"`
	Password  string     `json:"password,omitempty"`
}

func (user *UpdateUser) Update() string {

	query := ""

	if user.FirstName != "" {
		query += createQuery(" firstName ", user.FirstName)

	}
	if user.LastName != "" {
		query += createQuery(" lastName ", user.LastName)

	}

	query = strings.TrimSuffix(query, ",")

	return query

}
func createQuery(tag, value string) string {
	return tag + `= ` + "'" + value + "'" + " ,"

}

type QueryParams struct {
	FirstName string
	LastName  string
	Email     string
	Archived  string
	Sort      string
	Order     string
}

type Creds struct {
	Email    string
	Password string
}

type Claims struct {
	Email string `json:"email"`
	Id    string `json:"id"`
	jwt.StandardClaims
}

type JwtToken struct {
	Expires time.Time
	Token   string
}
