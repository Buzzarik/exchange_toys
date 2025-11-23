package handlers

import (
	"service/internal/models"
	"service/internal/parsers"
	"service/internal/service"
	"service/internal/utils"
	"service/internal/service/clients"

	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func isValidToyUser(app *service.Application, toyUser *models.UserIdToyId) (bool) {
	toy, err := app.Storage.SelectToyByUserId(toyUser.ToyId, toyUser.UserId)

	if err != nil || toy == nil {
		app.Log.Info("toy is not exist")

		return false
	}

	return true
}

func hasUserExchange(userId string, expectedUserId string) (bool) {
	return userId == expectedUserId
}

func isOtherUsers(userId1 string, userId2 string) (bool) {
	return userId1 != userId2
}

func getDetailsInfo(detail *models.ExchangeParticipant) (models.ExchangeDetailsInfo) {
	user := models.UserName{
		FirstName: detail.FirstName,
		LastName: detail.LastName,
		MiddleName: detail.MiddleName,
	}

	toy := models.ToyInfo{
		ToyId: detail.ToyId,
		UserId: detail.UserId,
		Name: detail.ToyName,
		Description: detail.ToyDescription,
		PhotoUrl: detail.ToyPhotoURL,
	}

	return models.ExchangeDetailsInfo{
		User: user,
		Toy: toy,
		Status: detail.UserExchangeStatus,
	}
}

func getExchange(exchange []models.ExchangeParticipant) (models.ExchangeInfo) {
	details := make([]models.ExchangeDetailsInfo, 0, len(exchange))
	for _, dbDetails := range(exchange) {
		details = append(details, getDetailsInfo(&dbDetails))
	}

	return models.ExchangeInfo{
		ExchangeId: exchange[0].ExchangeId,
		Details: details,
		IdempotencyToken: exchange[0].IdempotencyToken,
		Status: exchange[0].ExchangeStatus,
		CreatedAt: exchange[0].ExchangeCreatedAt,
		UpdatedAt: exchange[0].ExchangeUpdatedAt,
	}
}

func CreateExchange(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestExchangePost

		if err := parsers.ParseExchangePost(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error()})
		}

		app.Log.Info("Start POST v1/exchange", slog.Any("request", req))

        isValid := isValidToyUser(app, &req.Body.UserToy1) && 
                   isValidToyUser(app, &req.Body.UserToy2) &&
                   (hasUserExchange(req.UserId, req.Body.UserToy1.UserId) || 
                    hasUserExchange(req.UserId, req.Body.UserToy2.UserId)) &&
                   isOtherUsers(req.Body.UserToy1.UserId, req.Body.UserToy2.UserId)

        if !isValid {
            return context.Status(fiber.StatusBadRequest).JSON(
                models.ResponseError{
                    Code: models.KInvalidArgument,
                    Message: "toy is not exist or user is not initial exchange or user do not exchange with self"})
        }

		exchangeDetails := []models.ExchangeDetails {
			models.ExchangeDetails{
				ToyId: req.Body.UserToy1.ToyId,
				UserId: req.Body.UserToy1.UserId,
			},
			models.ExchangeDetails{
				ToyId: req.Body.UserToy2.ToyId,
				UserId: req.Body.UserToy2.UserId,
			},
		}

		exchange := models.Exchange{
			SrcToyId: req.Body.UserToy1.ToyId,
			DstToyId: req.Body.UserToy2.ToyId,
			IdempotencyToken: req.IdempotencyToken,
		}

		dbExchange, err := app.Storage.InsertExchange(&exchange, exchangeDetails)
		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidUpdateToy,
					Message: err.Error()})
		}

		return context.Status(fiber.StatusCreated).
				JSON(models.ResponseExchangePost{Exchange: *dbExchange})
	}
}

func GetExchange(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestExchangeGet

		if err := parsers.ParseExchangeGet(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error()})
		}

		app.Log.Info("Start GET v1/exchange", slog.Any("request", req))

		dbExchange, err := app.Storage.SelectExchangeWithParticipants(req.ExchangeId)
		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidGetExchange,
					Message: err.Error()})
		}

		if dbExchange == nil {
			return context.Status(fiber.StatusNotFound).JSON(
				models.ResponseError{
					Code: models.KExchangeNotFound,
					Message: "exchange not found"})
		}

		return context.Status(fiber.StatusOK).JSON(
			models.ResponseExchangeGet{
				Exchange: getExchange(dbExchange)})
	}
}

// TODO: в конце проверить статус, если будет confirm, то разослать сообщения
func PatchExchange(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestExchangePatch

		if err := parsers.ParseExchangePatch(&req, app, context); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(
				models.ResponseError{
					Code: models.KInvalidArgument,
					Message: err.Error()})
		}

		app.Log.Info("Start PATCH v1/exchange", slog.Any("request", req))

		dbExchange, err := app.Storage.UpdateExchangeWithParticipants(req.ExchangeId, req.UserId, req.Body.Status)
		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidUpdateExchangeStatus,
					Message: err.Error()})
		}

		if dbExchange == nil {
			return context.Status(fiber.StatusNotFound).JSON(
				models.ResponseError{
					Code: models.KExchangeNotFound,
					Message: "exchange not found"})
		}

		details := make([]models.ExchangeDetailsInfo, 0, len(dbExchange))
		for _, dbDetails := range(dbExchange) {
			details = append(details, getDetailsInfo(&dbDetails))
		}

		exchange := models.ExchangeInfo{
			ExchangeId: dbExchange[0].ExchangeId,
			Details: details,
			IdempotencyToken: dbExchange[0].IdempotencyToken,
			Status: dbExchange[0].ExchangeStatus,
			CreatedAt: dbExchange[0].ExchangeCreatedAt,
			UpdatedAt: dbExchange[0].ExchangeUpdatedAt,
		}

		if exchange.Status == models.KConfirmExchangeStatus {
			go clients.SendEmailToSingleParticipant(app, exchange.Details[0].Toy.UserId, exchange.Details[1].Toy.UserId)
			go clients.SendEmailToSingleParticipant(app, exchange.Details[1].Toy.UserId, exchange.Details[0].Toy.UserId)
		}
		// по логике надо тут транзакцию закрывать, когда будет уверенность, то корутины запущены

		return context.Status(fiber.StatusOK).JSON(
			models.ResponseExchangePatch{
				Exchange: exchange})
	}
}

func GetExchangeList(app *service.Application) fiber.Handler {
	return func(context *fiber.Ctx) error {
		var req models.RequestExchangeList

		if err := parsers.ParseExchangeList(&req, app, context); err != nil {
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

		app.Log.Info("Start POST v1/exchange/list", slog.Any("request", req), slog.Any("dd", req.Body.Query.Statuses))

		dbExchanges, cursor, err := app.Storage.SelectExchangeList(&req.Body.Query, req.UserId, cursor, *req.Body.Limit)
		if err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(
				models.ResponseError{
					Code: models.KInvalidExchangeList,
					Message: err.Error()})
		}

		if cursor != nil {
			cursor = utils.Encode(cursor)
		}

		detailsByExhangeId := make(map[string]([]models.ExchangeParticipant), len(dbExchanges))
		for _, dbDetails := range(dbExchanges) {
			if _, ok := detailsByExhangeId[dbDetails.ExchangeId]; !ok {
				detailsByExhangeId[dbDetails.ExchangeId] = make([]models.ExchangeParticipant, 0)
			}
			detailsByExhangeId[dbDetails.ExchangeId] = append(detailsByExhangeId[dbDetails.ExchangeId], dbDetails)
		}

		exchanges := make([]models.ExchangeInfo, 0, len(detailsByExhangeId))
		for _, exchange := range(detailsByExhangeId) {
			exchanges = append(exchanges, getExchange(exchange))
		}

		return context.Status(fiber.StatusOK).JSON(
			models.ResponseExchangeList{
				Exchanges: exchanges,
				Cursor: cursor})
	}
}