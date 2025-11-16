package models

import (
	"time"
)

type ExchangeStatus string
type ExchangeDetailsStatus string

const (
	KCreatedExchangeStatus ExchangeStatus = "created"
	KConfirmExchangeStatus ExchangeStatus = "confirm"
	KSuccessExchangeStatus ExchangeStatus = "success"
	KFailedExchangeStatus ExchangeStatus = "failed"

	KCreatedExchangeDetailsStatus ExchangeDetailsStatus = "created"
	KConfirm1ExchangeDetailsStatus ExchangeDetailsStatus = "confirm_1"
	KConfirm2ExchangeDetailsStatus ExchangeDetailsStatus = "confirm_2"
	KSuccessExchangeDetailsStatus ExchangeDetailsStatus = "success"
	KFailedExchangeDetailsStatus ExchangeDetailsStatus = "failed"
)

type Exchange struct {
	ExchangeId 	string 		`json:"exchange_id"`
	SrcToyId 	string 		`json:"src_toy_id"`
	DstToyId	string 		`json:"dst_toy_id"`
	IdempotencyToken string `json:"idempotency_token" validate:"required,min=1"`
	Status 		ExchangeStatus `json:"status"`
	CreatedAt 	time.Time  	`json:"created_at"`
	UpdatedAt 	time.Time  	`json:"updated_at"`
}

type ExchangeDetails struct {
	ExchangeId 	string 		`json:"exchange_id"`
	ToyId 		string 			`json:"toy_id"`
	UserId 		string			`json:"user_id"`
	Status 		ExchangeDetailsStatus `json:"status"`
	CreatedAt 	time.Time  	`json:"created_at"`
	UpdatedAt 	time.Time  	`json:"updated_at"`
}

type ExchangeParticipant struct {
    ExchangeId         string    `json:"exchange_id"`
    ExchangeStatus     ExchangeStatus    `json:"exchange_status"`
    IdempotencyToken   string    `json:"idempotency_token"`
    ExchangeCreatedAt  time.Time `json:"exchange_created_at"`
    ExchangeUpdatedAt  time.Time `json:"exchange_updated_at"`
    
    ToyId              string    `json:"toy_id"`
    ToyName            string    `json:"toy_name"`
    ToyDescription     *string   `json:"toy_description,omitempty" validate:"omitempty"`
    ToyPhotoURL        *string   `json:"toy_photo_url,omitempty" validate:"omitempty"`
    
    UserId             string    `json:"user_id"`
    FirstName          string    `json:"first_name"`
    MiddleName         *string   `json:"middle_name,omitempty" validate:"omitempty"`
    LastName           string    `json:"last_name"`
    
    UserExchangeStatus ExchangeDetailsStatus    `json:"user_exchange_status"`
}

type ExchangeDetailsInfo struct {
	Toy 		ToyInfo 	`json:"toy"`
	User	 	UserName	`json:"user"`
	Status 		ExchangeDetailsStatus `json:"status"`
}

type ExchangeInfo struct {
	ExchangeId 	string 		`json:"exchange_id"`
	Details []ExchangeDetailsInfo `json:"exchange_details"`
	IdempotencyToken string `json:"idempotency_token" validate:"required,min=1"`
	Status 		ExchangeStatus `json:"status"`
	CreatedAt 	time.Time  	`json:"created_at"`
	UpdatedAt 	time.Time  	`json:"updated_at"`
}

type UserIdToyId struct {
	UserId 	string `json:"user_id" validate:"required,min=1"`
	ToyId 	string `json:"toy_id" validate:"required,min=1"`
}

type RequestExchangePostBody struct {
	UserToy1 UserIdToyId `json:"user_toy_1" validate:"required"`
	UserToy2 UserIdToyId `json:"user_toy_2" validate:"required"`
}

type RequestExchangePost struct {
	UserId	string		 			`json:"user_id" validate:"required,min=1"`
	IdempotencyToken string 		`json:"idempotency_token" validate:"required,min=1"`
	Body 	RequestExchangePostBody `json:"body" validate:"required"`
}

type RequestExchangeGet struct {
	UserId string `json:"user_id" validate:"required,min=1"`
	ExchangeId string `json:"exchange_id" validate:"required,min=1"`
}

type RequestExchangePatchBody struct {
	Status ExchangeDetailsStatus `json:"status" validate:"required,oneof=confirm_1 confirm_2 failed"`
}

type RequestExchangePatch struct {
	UserId string `json:"user_id" validate:"required,min=1"`
	ExchangeId string `json:"exchange_id" validate:"required,min=1"`
	Body RequestExchangePatchBody `json:"body"`
}

type RequestExchangeListBody struct {
	Query QueryExchanges `json:"query" validate:"required"`
	Limit *int64 `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Cursor *string `json:"cursor,omitempty" validate:"omitempty,min=1"`
}

type RequestExchangeList struct {
	UserId string `json:"user_id" validate:"required,min=1"`
	Body RequestExchangeListBody `json:"body" validate:"required"`
}

//response
type ResponseExchangeList struct {
	Exchanges []ExchangeInfo `json:"exchanges" validate:"required"`
	Cursor *string `json:"cursor,omitempty" validate:"omitempty,min=1"`
}

type ResponseExchangePatch struct {
	Exchange ExchangeInfo `json:"exchange"`
}

type ResponseExchangeGet struct {
	Exchange ExchangeInfo `json:"exchange"`
}

type ResponseExchangePost struct {
	Exchange Exchange `json:"exchange"`
}
