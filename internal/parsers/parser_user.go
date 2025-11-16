package parsers

import (
	"service/internal/models"
	"service/internal/service"

	"github.com/gofiber/fiber/v2"
)

const (
	kConfirmPassword = "confirm_password"
	kPassword = "password"
	kFirstName = "first_name"
	kLastName = "last_name"
	kMiddleName = "middle_name"
	kEmail = "email"
)

func ParseRegister(req *models.RequestRegister, app *service.Application, context *fiber.Ctx) (error) {
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

func ParseLogin(req *models.RequestLogin, app *service.Application, context *fiber.Ctx) (error) {
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