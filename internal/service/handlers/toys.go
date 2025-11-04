package handlers

import (
	"service/internal/models"
	"service/internal/service"
	"service/internal/utils"

	"log/slog"
	"time"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Тестовые данные
var kPhotoUrl string = "/100days.png"
var kDescripton string = "hello"
var kId int64 = 3

var toys = map[string]*models.Toy {
		"toy_id_1": &models.Toy{
			ToyId: 			"toy_id_1",
			UserId: 		"user_id_1",
			Name:			"toy_1",
			IdempotencyToken: "token_1",
			Description: 	&kDescripton,
			PhotoUrl:		&kPhotoUrl,
			Status: 		"created",
			CreatedAt: 		time.Now(),
			UpdatedAt: 		time.Now(),
		},
		"toy_id_2": &models.Toy{
			ToyId: 			"toy_id_2",
			UserId: 		"user_id_2",
			Name:			"toy_2",
			IdempotencyToken: "token_2",
			Description: 	nil,
			PhotoUrl:		nil,
			Status: 		"exchanging",
			CreatedAt: 		time.Now(),
			UpdatedAt: 		time.Now(),
		},
}

//TODO: доделать взаимодействие с БД + проверка на idempotency Update Returning + фотка
func CreateToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyPost

		if err := req.Parse(app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error(),
				},
			)
		}

		app.Log.Info("Start POST v1/toys", slog.Any("request", req))

		new_id := "toy_id_" + strconv.Itoa(int(kId))
		kId++
		toys[new_id] = &models.Toy{
			ToyId: new_id,
			Name: req.Toy.Name,
			Description: req.Toy.Description,
			IdempotencyToken: req.Toy.IdempotencyToken,
			UserId: req.UserId,
			Status: "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		return context.SendStatus(fiber.StatusCreated)
	}
}

//TODO: доделать взаимодействие с БД + фотка
func UpdateToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyPut

		if err := req.Parse(app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error(),
				},
			)
		}

		app.Log.Info("Start PUT v1/toys", slog.Any("request", req))

		if toy, ok := toys[req.Toy.ToyId]; ok {
			toy.Name = req.Toy.Name
			toy.Description = req.Toy.Description

			return context.Status(fiber.StatusOK).JSON(
				models.ReponseToyPut{
					Toy: toy,
				},
			)
		}

		return context.Status(fiber.StatusNotFound).JSON(
			models.ResponseError{
				Code: models.KToyNotFound,
				Message: "Toy is not exist",
			},
		)
	}
}

//TODO: доделать взаимодействие с БД
func GetToysList(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToysList

		if err := req.Parse(app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error(),
				},
			)
		}

		cursor, err := utils.Decode(req.Body.Cursor)
		if err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidCursor,
					Message: err.Error(),
				},
			)
		}

		app.Log.Info("Start POST v1/toys/list", slog.Any("request", req))

		// тест
		var resp models.ResponseToysList
		resp.Toys = make([]models.Toy, 0, len(toys)) // cap 
		for _, toy := range(toys) {
			resp.Toys = append(resp.Toys, *toy)
		}
		resp.Cursor = cursor

		return context.Status(fiber.StatusOK).JSON(resp)

	}
}

//TODO: доделать взаимодействие с БД + как надо фейлить сделки по этой игрушке, если будет статус created
func UpdateStatusToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyPatch

		if err := req.Parse(app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error(),
				},
			)
		}

		app.Log.Info("Start PATCH v1/toys", slog.Any("request", req))

		if toy, ok := toys[req.ToyId]; ok {
			toy.Status = req.Body.Status
		}

		return context.SendStatus(fiber.StatusOK)

	}
}

//TODO: доделать взаимодействие с БД + как надо фейлить сделки по этой игрушке
func DeleteToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyDelete

		if err := req.Parse(app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error(),
				},
			)
		}

		app.Log.Info("Start DELETE v1/toys", slog.Any("request", req))

		if toy, ok := toys[req.ToyId]; ok {
			toy.Status = "removed"	
		}

		return context.SendStatus(fiber.StatusOK)

	}
}

//TODO: доделать взаимодействие с БД
func GetToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyGet

		if err := req.Parse(app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error(),
				},
			)
		}

		app.Log.Info("Start GET v1/toys", slog.Any("request", req))

		if toy, ok := toys[req.ToyId]; ok {
			
			return context.Status(fiber.StatusOK).JSON(
				models.ReponseToyGet{
					Toy: toy,
				},
			)
		}

		return context.Status(fiber.StatusNotFound).JSON(
			models.ResponseError{
				Code: models.KToyNotFound,
				Message: "Toy is not exist",
			},
		)

	}
}
