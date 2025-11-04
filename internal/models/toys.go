package models

import(
	"service/internal/service"

	"time"
	"mime/multipart"
	_ "net/http"

	"github.com/gofiber/fiber/v2"
)

const(
	kToyId = "toy_id"
	kUserId = "user_id"
	kDescripton = "description"
	kName = "name"
	kStatus = "status"
	kIdempotencyToken = "idempotency_token"
	kLimit int64 = 40
)

type Toy struct {
	ToyId 		string 		`json:"toy_id"`
	UserId 		string 		`json:"user_id"`
	Name 		string 		`json:"name"`
	IdempotencyToken string `json:"idempotency_token"`
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

func (req *RequestToyGet) Parse(app *service.Application, context *fiber.Ctx) error {
	req.ToyId = context.Params(kToyId)
	req.UserId = context.Get(kUserId)

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	return nil
}

type RequestToyDelete struct {
	ToyId string `json:"toy_id" validate:"required,min=1"`
	UserId string `json:"user_id" validate:"required,min=1"`
}

func (req *RequestToyDelete) Parse(app *service.Application, context *fiber.Ctx) error {
	req.ToyId = context.Params(kToyId)
	req.UserId = context.Get(kUserId)

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	return nil
}

type RequestToyPatchBody struct {
	Status string `json:"status" validate:"required,oneof=created exchanging"`
}

type RequestToyPatch struct {
	ToyId string `json:"toy_id" validate:"required,min=1"`
	UserId string `json:"user_id" validate:"required,min=1"`
	Body RequestToyPatchBody `json:"body" validate:"required"`
}

func (req *RequestToyPatch) Parse(app *service.Application, context *fiber.Ctx) (error) {
	req.ToyId = context.Params(kToyId)
	req.UserId = context.Get(kUserId)

	if err := context.BodyParser(&req.Body); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	return nil
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

func (req *RequestToysList) Parse(app *service.Application, context *fiber.Ctx) (error) {
	req.UserId = context.Get(kUserId)

	if err := context.BodyParser(&req.Body); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	if req.Body.Limit == nil {
		limit := kLimit
		req.Body.Limit = &limit
	}

	return nil
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

func (req *RequestToyPut) Parse(app *service.Application, context *fiber.Ctx) (error) {
	req.UserId = context.Get(kUserId)
	req.Toy.ToyId = context.Get(kToyId)
	req.Toy.Name = context.FormValue(kName)

	if description := context.FormValue(kDescripton); description != "" {
		req.Toy.Description = &description
	} 

	if file, err := context.FormFile("file"); err != nil {
		app.Log.Info("File not added")
	} else {
		req.File = file
	}

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	return nil
}

type RequestToyPostBody struct {
	Name 		string 		`json:"name" validate:"required,min=1"`
	Description *string 	`json:"description,omitempty" validate:"omitempty"`
	IdempotencyToken string `json:"idempotency_token" validate:"required,min=1"`
}

type RequestToyPost struct {
	UserId string `json:"user_id" validate:"required,min=1"`
	Toy RequestToyPostBody `json:"toy" validate:"required"`
	File *multipart.FileHeader	   `json:"file,omitempty" validate:"omitempty"`
}

func (req *RequestToyPost) Parse(app *service.Application, context *fiber.Ctx) (error) {
	req.UserId = context.Get(kUserId)
	req.Toy.IdempotencyToken = context.Get(kIdempotencyToken)
	req.Toy.Name = context.FormValue(kName)

	if description := context.FormValue(kDescripton); description != "" {
		req.Toy.Description = &description
	} 

	if file, err := context.FormFile("file"); err != nil {
		app.Log.Info("File not added")
	} else {
		req.File = file
	}

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	return nil
}

// Response 
type ReponseToyGet struct {
	Toy *Toy `json:"toy" validate:"required"`
}

type ResponseToysList struct {
	Toys []Toy `json:"toys" validate:"required"`
	Cursor *string `json:"cursor,omitempty" validate:"omitempty,min=1"`
}

type ReponseToyPut struct {
	Toy *Toy `json:"toy" validate:"required"`
}
