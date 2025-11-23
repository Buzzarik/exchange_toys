package models

import (
	"strings"
	"time"
)

type User struct {
	UserId string `json:"user_id" validate:"required,min=1"`
	UserName UserName `json:"user_name" validate:"required"`
	HashPassword string `json:"password" validate:"required,min=1"`
	Email string `json:"email" validate:"required,email"`
	CreatedAt 	time.Time  	`json:"created_at"`
	UpdatedAt 	time.Time  	`json:"updated_at"`
}

func (u *User) FullName() string {
    parts := []string{u.UserName.LastName, u.UserName.FirstName}
    
    if u.UserName.MiddleName != nil {
        parts = append(parts, *u.UserName.MiddleName)
    }
    
    return strings.Join(parts, " ")
}

type UserName struct {
	FirstName string 	`json:"first_name" validate:"required,min=1"`
	LastName string 	`json:"last_name" validate:"required,min=1"`
	MiddleName *string 	`json:"middle_name,omitempty" validate:"omitempty"`
}

type RequestRegisterBody struct {
	UserName UserName `json:"user_name"`
	Password string   `json:"password" validate:"required,min=1"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email"`
}

type RequestRegister struct {
	Body RequestRegisterBody `json:"body"`
}

type RequestLoginBody struct {
	Password string   `json:"password" validate:"required,min=1"`
	Email           string `json:"email" validate:"required,email"`
}

type RequestLogin struct {
	Body RequestLoginBody `json:"body"`
}

type ResponseRegister struct {
	UserId string `json:"user_id" validate:"required,min=1"`
}

type ResponseLogin struct {
	UserId string `json:"user_id" validate:"required,min=1"`
}