package entity

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Id        string    `json:"id" validate:"omitempty,uuid"`
	FirstName string    `json:"firstName" validate:"required,gte=2,lte=20" `
	LastName  string    `json:"lastName" validate:"required,gte=0,lte=20"`
	Email     string    `json:"email" validate:"required,email,gte=3,lte=20"`
	Dob       time.Time `json:"dob" validate:"required"`
	Password  string    `json:"password,omitempty" validate:"required,gte=8,lte=20"`
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
type UserResponse struct {
	Status int        `json:"status"`
	Data   []Response `json:"data"`
	Count  int        `json:"count"`
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
	Limit     int
	Page      int
}

type Creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
