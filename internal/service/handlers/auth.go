package handlers

import (
	"service/internal/models"
	"service/internal/parsers"
	"service/internal/service"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

//TODO: v1/register/form POST создание пользователя

//TODO: v2/login/form POST заход на аккаунт
func Register(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestRegister

		if err := parsers.ParseRegister(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error()})
		}

		app.Log.Info("Start POST v1/register", slog.Any("request", req))

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Body.Password), bcrypt.DefaultCost)
		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidRegister,
					Message: err.Error()})
		}

		user := models.User{
			UserName: models.UserName{
				FirstName: req.Body.UserName.FirstName,
				LastName: req.Body.UserName.LastName,
				MiddleName: req.Body.UserName.MiddleName,
			},
			HashPassword: string(hashedPassword),
			Email: req.Body.Email,
		}

		dbUser, err := app.Storage.CreateUser(&user)

		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidRegister,
					Message: err.Error(),
				},
			)
		}

		if dbUser == nil {
			return context.Status(fiber.StatusConflict).JSON(
				models.ResponseError{
					Code: models.KExistUser,
					Message: "user already exist",
				},
			)
		}

		return context.Status(fiber.StatusCreated).JSON(
			models.ResponseRegister{UserId: dbUser.UserId})
	}
}

func Login(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestLogin

		if err := parsers.ParseLogin(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error()})
		}

		app.Log.Info("Start POST v1/login", slog.Any("request", req))

		user := models.User{
			Email: req.Body.Email,
		}

		dbUser, err := app.Storage.SelectUserByEmail(&user)

		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidLogin,
					Message: err.Error(),
				},
			)
		}

		if dbUser != nil && bcrypt.CompareHashAndPassword([]byte(dbUser.HashPassword), []byte(req.Body.Password)) == nil {
			return context.Status(fiber.StatusOK).JSON(
				models.ResponseLogin{UserId: dbUser.UserId})
		}

		return context.Status(fiber.StatusNotFound).JSON(
			models.ResponseError{
				Code: models.KInvalidVerify,
				Message: "invalid username or password"})

	}
}