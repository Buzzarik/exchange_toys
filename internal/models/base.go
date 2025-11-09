package models

const (
	KToyNotFound = "Toy not found"
	KInvalidArgument = "Invalid argument"
	KInvalidCursor = "Invalid cursor"
	KErrorSaveFile = "File is not save"
	KInvalidCreateToy = "Invalid create toy"
	KInvalidUpdateToyStatus = "Invalid update toy status"
	KInvalidGetToy = "Invalid get toy"
	KInvalidUpdateToy = "Invalid update toy"
	KInvalidToysList = "Invalid toys list"
)

type ResponseError struct { // может быть в общий компонент перенести потому что встроенный как-то не о чем
	Code string `json:"code" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type QueryToys struct {
	Statuses []string `json:"statuses,omitempty" validate:"omitempty,min=1,dive,oneof=created exchanging removed"`
	UserIds []string  `json:"user_ids,omitempty" validate:"omitempty,min=1,dive,min=1"`
	ExcludeUserIds []string `json:"exclude_user_ids,omitempty" validate:"omitempty,min=1,dive,min=1"`
}
