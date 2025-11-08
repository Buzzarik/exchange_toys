package parsers

import (
	"service/internal/service"
	"service/internal/models"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

const(
	kToyId = "toy_id"
	kXUserId = "x_user_id"
	kDescripton = "description"
	kName = "name"
	kStatus = "status"
	kXIdempotencyToken = "x_idempotency_token"
	kLimit int64 = 40
)

// чтобы пофиксить баг с сохранением больше чем одного пользовательского заголовка
func getHeader(context *fiber.Ctx, key string) string {
	return fmt.Sprintf("%s", context.Get(key))
}

func ParseToyGet(req *models.RequestToyGet, app *service.Application, context *fiber.Ctx) error {
	req.ToyId = context.Params(kToyId)
	req.UserId = getHeader(context, kXUserId)

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	return nil
}

func ParseToyDelete(req *models.RequestToyDelete, app *service.Application, context *fiber.Ctx) (error) {
	req.ToyId = context.Params(kToyId)
	req.UserId = getHeader(context, kXUserId)

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	return nil
}

func ParseToyPatch(req *models.RequestToyPatch, app *service.Application, context *fiber.Ctx) (error) {
	req.ToyId = context.Params(kToyId)
	req.UserId = getHeader(context, kXUserId)

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

func ParseToysList(req *models.RequestToysList, app *service.Application, context *fiber.Ctx) (error) {
	req.UserId = getHeader(context, kXUserId)

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

func ParseToyPut(req *models.RequestToyPut, app *service.Application, context *fiber.Ctx) (error) {
	req.UserId = getHeader(context, kXUserId)
	req.Toy.ToyId = getHeader(context, kToyId)
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

func ParseToyPost(req *models.RequestToyPost, app *service.Application, context *fiber.Ctx) (error) {
	req.UserId = getHeader(context, kXUserId)
	req.IdempotencyToken = getHeader(context, kXIdempotencyToken)
	req.Toy.Name = context.FormValue(kName)

	if description := context.FormValue(kDescripton); description != "" {
		req.Toy.Description = &description
	} 

	// также можно добавить проверку типов jpg, png и тд
	if file, err := context.FormFile("file"); err != nil {
		app.Log.Info("File not added")
		req.File = nil
	} else {
		req.File = file
	}

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	return nil
}