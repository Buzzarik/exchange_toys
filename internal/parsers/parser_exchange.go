package parsers

import (
	"service/internal/service"
	"service/internal/models"

	"github.com/gofiber/fiber/v2"
)

func ParseExchangePost(req *models.RequestExchangePost, app *service.Application, context *fiber.Ctx) error {
	req.UserId = getHeader(context, kXUserId)
	req.IdempotencyToken = getHeader(context, kXIdempotencyToken)

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

func ParseExchangeGet(req *models.RequestExchangeGet, app *service.Application, context *fiber.Ctx) error {
	req.UserId = getHeader(context, kXUserId)
	req.ExchangeId = context.Params(kExchangeId)

	if err := app.Validator.Struct(req); err != nil {
		app.Log.Warn(err.Error())

		return err
	}

	return nil
}

func ParseExchangePatch(req *models.RequestExchangePatch, app *service.Application, context *fiber.Ctx) error {
	req.UserId = getHeader(context, kXUserId)
	req.ExchangeId = context.Params(kExchangeId)

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

func ParseExchangeList(req *models.RequestExchangeList, app *service.Application, context *fiber.Ctx) (error) {
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