package handlers

import (
	"service/internal/parsers"
	"service/internal/models"
	"service/internal/service"
	"service/internal/utils"

	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	kCreated = "created"
)

// Тестовые данные
var kPhotoUrl string = "/100days.png"
var kDescripton string = "hello"

var iToken string

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

func saveFile(req *models.RequestToyPost, app *service.Application, context *fiber.Ctx) (*string, error) {
	if req.File == nil {
		return nil, nil
	}

	filename := fmt.Sprintf(
		"%s_%s",
		uuid.New().String(), 
		strings.ReplaceAll(req.File.Filename, " ", "_"),
	)

	path := fmt.Sprintf(
		"%s/%s", 
		app.Cnf.Server.Prefix_upload, 
		filename,
	)

	photo_url := fmt.Sprintf(
		"%s/%s", 
		app.Cnf.Server.PhotoUrl, 
		filename,
	)

	if err := context.SaveFile(req.File, path); err != nil {
		return nil, err
	}

	app.Log.Info(fmt.Sprintf("Created file %s", path))

	return &photo_url, nil
}

//TODO: доделать взаимодействие с БД + проверка на idempotency Update Returning + фотка
func CreateToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyPost

		if err := parsers.ParseToyPost(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error(),
				},
			)
		}

		app.Log.Info("Start POST v1/toys", slog.Any("request", req))

		dbToy, err := app.Storage.SelectToyByToken(req.IdempotencyToken)
		if err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KCreateToyError,
					Message: err.Error(),
				},
			)
		}

		if dbToy != nil {
			return context.Status(fiber.StatusCreated).
				JSON(models.ReponseToyPost{Toy: *dbToy})
		}
		
		new_id := uuid.New().String()
		photo_url, err := saveFile(&req, app, context)
		if err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KErrorSaveFile,
					Message: err.Error(),
				},
			)
		}

		toy := models.Toy{
			ToyId: new_id,
			Name: req.Toy.Name,
			IdempotencyToken: req.IdempotencyToken,
			Description: req.Toy.Description,
			PhotoUrl: photo_url,
			UserId: req.UserId,
			Status: kCreated,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		dbToy, err = app.Storage.InsertToy(&toy)
		if err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KCreateToyError,
					Message: err.Error(),
			})
		}

		// for _, toy := range toys {
		// 	app.Log.Info("", slog.Any("toy", toy))
		// }

		return context.Status(fiber.StatusCreated).
				JSON(models.ReponseToyPost{Toy: *dbToy})
	}
}

//TODO: доделать взаимодействие с БД + фотка
func UpdateToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyPut

		if err := parsers.ParseToyPut(&req, app, context); err != nil {
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
					Toy: *toy,
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

		if err := parsers.ParseToysList(&req, app, context); err != nil {
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

		if err := parsers.ParseToyPatch(&req, app, context); err != nil {
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

		if err := parsers.ParseToyDelete(&req, app, context); err != nil {
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

		if err := parsers.ParseToyGet(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error(),
				},
			)
		}

		app.Log.Info("Start GET v1/toys", slog.Any("request", req))

		// for _, toy := range toys {
		// 	app.Log.Info("", slog.Any("toy", toy))
		// }

		toy, err := app.Storage.SelectToyById(req.ToyId, req.UserId)

		if err != nil {
			return context.Status(fiber.StatusBadRequest).JSON( // 400 - заменить
				models.ResponseError{
					Code: models.KToyNotFound,
					Message: err.Error(),
				},
			)
		}

		if toy == nil {
			return context.Status(fiber.StatusNotFound).JSON(
				models.ResponseError{
					Code: models.KToyNotFound,
					Message: err.Error(),
				},
			)
		}

		return context.Status(fiber.StatusOK).JSON(
			models.ReponseToyGet{Toy: *toy})
	}
}
