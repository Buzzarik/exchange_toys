package models

import (
	"mime/multipart"
	"time"
)

type Toy struct {
	ToyId 		string 		`json:"toy_id"`
	UserId 		string 		`json:"user_id"`
	Name 		string 		`json:"name"`
	IdempotencyToken string `json:"idempotency_token" validate:"required,min=1"`
	Description *string 	`json:"description,omitempty" validate:"omitempty"`
	PhotoUrl 	*string 	`json:"photo_url,omitempty" validate:"omitempty"`
	Status 		string 		`json:"status"`
	CreatedAt 	time.Time  	`json:"created_at"`
	UpdatedAt 	time.Time  	`json:"updated_at"`
}


// Request
type RequestToyGet struct {
	ToyId string `json:"toy_id" validate:"required,min=1"`
	UserId string `json:"user_id" validate:"required,min=1"`
}

type RequestToyDelete struct {
	ToyId string `json:"toy_id" validate:"required,min=1"`
	UserId string `json:"user_id" validate:"required,min=1"`
}

type RequestToyPatchBody struct {
	Status string `json:"status" validate:"required,oneof=created exchanging"`
}

type RequestToyPatch struct {
	ToyId string `json:"toy_id" validate:"required,min=1"`
	UserId string `json:"user_id" validate:"required,min=1"`
	Body RequestToyPatchBody `json:"body" validate:"required"`
}

type RequestToysListBody struct {
	Query QueryToys `json:"query" validate:"required"`
	Limit *int64 `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Cursor *string `json:"cursor,omitempty" validate:"omitempty,min=1"`
}

type RequestToysList struct {
	UserId string `json:"user_id" validate:"required,min=1"`
	Body RequestToysListBody `json:"body" validate:"required"`
}

type RequestToyPutBody struct {
	ToyId 		string 		`json:"toy_id" validate:"required,min=1"`
	Name 		string 		`json:"name" validate:"required,min=1"`
	Description *string 	`json:"description,omitempty" validate:"omitempty"`
}

type RequestToyPut struct {
	UserId string `json:"user_id" validate:"required,min=1"`
	Toy RequestToyPutBody `json:"toy" validate:"required"`
	File *multipart.FileHeader	   `json:"file,omitempty" validate:"omitempty"`
}

type RequestToyPostBody struct {
	Name 		string 		`json:"name" validate:"min=1"`
	Description *string 	`json:"description,omitempty" validate:"omitempty"`
}

type RequestToyPost struct {
	UserId string `json:"user_id" validate:"required,min=1"`
	Toy RequestToyPostBody `json:"toy"`
	IdempotencyToken string `json:"idempotency_token" validate:"min=1"`
	File *multipart.FileHeader	   `json:"file,omitempty" validate:"omitempty"`
}

// Response 
type ReponseToyGet struct {
	Toy Toy `json:"toy" validate:"required"`
}

type ReponseToyPost struct {
	Toy Toy `json:"toy" validate:"required"`
}

type ResponseToysList struct {
	Toys []Toy `json:"toys" validate:"required"`
	Cursor *string `json:"cursor,omitempty" validate:"omitempty,min=1"`
}

type ReponseToyPut struct {
	Toy Toy `json:"toy" validate:"required"`
}
