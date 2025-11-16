package handlers

import (
	"service/internal/models"
	"service/internal/parsers"
	"service/internal/service"
	"service/internal/utils"

	"fmt"
	"log/slog"
	"strings"
	"time"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RequestFile interface {
	GetFile() *multipart.FileHeader
}

func saveFile[T RequestFile](req T, app *service.Application, context *fiber.Ctx) (*string, error) {
	if req.GetFile() == nil {
		return nil, nil
	}

	filename := fmt.Sprintf(
		"%s_%s",
		uuid.New().String(), 
		strings.ReplaceAll(req.GetFile().Filename, " ", "_"),
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

	if err := context.SaveFile(req.GetFile(), path); err != nil {
		return nil, err
	}

	app.Log.Info(fmt.Sprintf("Created file %s", path))

	return &photo_url, nil
}

func CreateToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyPost

		if err := parsers.ParseToyPost(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error()})
		}

		app.Log.Info("Start POST v1/toys", slog.Any("request", req))

		dbToy, err := app.Storage.SelectToyByToken(req.IdempotencyToken)
		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidCreateToy,
					Message: err.Error()})
		}

		if dbToy != nil {
			return context.Status(fiber.StatusCreated).
				JSON(models.ReponseToyPost{Toy: *dbToy})
		}
		
		photoUrl, err := saveFile(&req, app, context)
		if err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KErrorSaveFile,
					Message: err.Error()})
		}

		toy := models.Toy{
			Name: req.Toy.Name,
			IdempotencyToken: req.IdempotencyToken,
			Description: req.Toy.Description,
			PhotoUrl: photoUrl,
			UserId: req.UserId,
			Status: models.KCreatedToyStatus,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		dbToy, err = app.Storage.InsertToy(&toy)
		if err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidCreateToy,
					Message: err.Error()})
		}

		return context.Status(fiber.StatusCreated).
				JSON(models.ReponseToyPost{Toy: *dbToy})
	}
}

func UpdateToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyPut

		if err := parsers.ParseToyPut(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error()})
		}

		app.Log.Info("Start PUT v1/toys", slog.Any("request", req))
		
		photoUrl, err := saveFile(&req, app, context)
		if err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KErrorSaveFile,
					Message: err.Error()})
		}

		toy := models.Toy{
			ToyId: req.Toy.ToyId,
			UserId: req.UserId,
			Description: req.Toy.Description,
			Name: req.Toy.Name,
		}
		if photoUrl != nil {
			toy.PhotoUrl = photoUrl
		}

		dbToy, err := app.Storage.UpdateToy(&toy)
		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidUpdateToy,
					Message: err.Error()})
		}

		if dbToy == nil {
			return context.Status(fiber.StatusNotFound).JSON(
				models.ResponseError{
					Code: models.KToyNotFound,
					Message: "Toy is not exist"})
		}

		return context.Status(fiber.StatusOK).
				JSON(models.ResponseToyPut{Toy: *dbToy})
	}
}

func GetToysList(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToysList

		if err := parsers.ParseToysList(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error()})
		}

		cursor, err := utils.Decode(req.Body.Cursor)
		if err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidCursor,
					Message: err.Error()})
		}

		app.Log.Info("Start POST v1/toys/list", slog.Any("request", req))

		dbToys, cursor, err := app.Storage.SelectToysList(&req.Body.Query, cursor, *req.Body.Limit)
		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidToysList,
					Message: err.Error()})
		}

		if cursor != nil {
			cursor = utils.Encode(cursor)
		}

		return context.Status(fiber.StatusOK).JSON(
			models.ResponseToysList{
				Toys: dbToys,
				Cursor: cursor})
	}
}

func UpdateToyStatus(app *service.Application) fiber.Handler {
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

		dbToy, err := app.Storage.UpdateToyStatus(req.ToyId, req.UserId, models.ToyStatus(req.Body.Status))

		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidUpdateToyStatus,
					Message: err.Error(),
				},
			)
		}

		if dbToy == nil {
			return context.Status(fiber.StatusNotFound).JSON(
				models.ResponseError{
					Code: models.KToyNotFound,
					Message: "Error change status toy"})
		}

		return context.Status(fiber.StatusOK).
			JSON(models.ReponseToyPatch{Toy: *dbToy})
	}
}

func DeleteToy(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestToyDelete

		if err := parsers.ParseToyDelete(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error()})
		}

		app.Log.Info("Start DELETE v1/toys", slog.Any("request", req))

		_, err := app.Storage.UpdateToyStatus(req.ToyId, req.UserId, models.KRemovedToyStatus)

		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidUpdateToyStatus,
					Message: err.Error(),
				},
			)
		}

		return context.SendStatus(fiber.StatusOK)

	}
}

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
		toy, err := app.Storage.SelectToyById(req.ToyId)

		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidGetToy,
					Message: err.Error(),
				},
			)
		}

		if toy == nil {
			return context.Status(fiber.StatusNotFound).JSON(
				models.ResponseError{
					Code: models.KToyNotFound,
					Message: "toy not found"})
		}

		return context.Status(fiber.StatusOK).JSON(
			models.ReponseToyGet{Toy: *toy})
	}
}
