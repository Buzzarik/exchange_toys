package middlewares

import (
	"service/internal/models"
	"service/internal/service"

	"github.com/gofiber/fiber/v2"
)

const (
	kXUserId = "x_user_id"
)

func AuthMiddleware(app *service.Application) fiber.Handler {
    return func(c *fiber.Ctx) error {

		app.Log.Info("Checking authorization")

        userId := c.Get(kXUserId)
        if userId == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(models.ResponseError{
                Code:    models.KUnauthorized,
                Message: "user_id header is required",
            })
        }
        

        user, err := app.Storage.SelectUserById(&models.User{UserId: userId})
        if err != nil || user == nil {
            return c.Status(fiber.StatusUnauthorized).JSON(models.ResponseError{
                Code:    models.KUnauthorized,
                Message: "invalid user_id",
            })
        }
        
        return c.Next()
    }
}