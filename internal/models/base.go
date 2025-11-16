package models

const (
	KToyNotFound = "Toy not found"
	KExchangeNotFound = "Exchange not found"
	KInvalidArgument = "Invalid argument"
	KInvalidCursor = "Invalid cursor"
	KErrorSaveFile = "File is not save"
	KInvalidCreateToy = "Invalid create toy"
	KInvalidUpdateToyStatus = "Invalid update toy status"
	KInvalidRegister = "Invalid register"
	KInvalidLogin = "Invalid login"
	KInvalidGetToy = "Invalid get toy"
	KInvalidUpdateToy = "Invalid update toy"
	KInvalidToysList = "Invalid toys list"
	KInvalidExchangeList = "Invalid exchange list"
	KInvalidCreateExchange = "Invalid create exchange"
	KInvalidGetExchange = "Invalid get exchange"
	KInvalidUpdateExchangeStatus = "Invalid update exchange status"
	KInvalidVerify = "Invalid verify"
	KUnauthorized = "Unauthorized"
	KExistUser = "User is exist"
)

type ResponseError struct {
	Code string `json:"code" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type QueryToys struct {
	Statuses []string `json:"statuses,omitempty" validate:"omitempty,min=1,dive,oneof=created exchanging removed"`
	UserIds []string  `json:"user_ids,omitempty" validate:"omitempty,min=1,dive,min=1"`
	ExcludeUserIds []string `json:"exclude_user_ids,omitempty" validate:"omitempty,min=1,dive,min=1"`
}

type QueryExchanges struct {
	Statuses []string `json:"statuses,omitempty" validate:"omitempty,min=1,dive,oneof=created confirm success failed"`
}
