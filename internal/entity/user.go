package entity

import "time"

type User struct {
	Id        int        `json:"id"`
	FirstName *string    `json:"firstName"`
	LastName  *string    `json:"lastName"`
	Email     *string    `json:"email"`
	Dob       *time.Time `json:"dob"`
	Password  string     `json:"password"`
}
