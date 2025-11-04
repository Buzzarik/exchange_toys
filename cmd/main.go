package main

import (
	"service/internal/config"
	"service/internal/service/handlers"
	"service/internal/service"

	"fmt"
	"log/slog"
	"os"

    "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/go-playground/validator/v10"
)

type Response struct{
	User string `json:"name"`
	Message string `json:"message"`
}

type Request struct{
	User string `json:"user"`
}

const uploadDir = "./uploads" // Директория, где хранятся загруженные файлы (тогда в конфиг)


//NOTE: CONFIG_PATH=./config/local_config.yaml go run ./cmd/main.go
func main() {
    // Создаем новый экземпляр Fiber
    app := fiber.New()

	// считываем с конфига
	cnf := config.New();

	application := &service.Application{
		Cnf: cnf,
		Log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Validator: validator.New(),
	}

	// настройка middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}:${port} ${method} ${path} ${status} ${reqHeader:User-Agent}\n",
	}))
	app.Use(recover.New())
	//TODO: добавить проверку пользователя

	app.Static("/upload", cnf.Server.Prefix_upload)

	toysV1Group := app.Group("/v1/toys");
	{
		toysV1Group.Post("/", handlers.CreateToy(application))
		toysV1Group.Put("/", handlers.UpdateToy(application))
		toysV1Group.Patch("/:toy_id", handlers.UpdateStatusToy(application))
		toysV1Group.Delete("/:toy_id", handlers.DeleteToy(application))
		toysV1Group.Get("/:toy_id", handlers.GetToy(application))
		toysV1Group.Post("/list", handlers.GetToysList(application))
	}

	// exchangeV1Group := app.Group("/v1/exchange")
	// {
	// 	exchangeV1Group.Get("/:exchange_id", )
	// 	exchangeV1Group.Patch("/:exchange_id", )
	// 	exchangeV1Group.Post("/list", )
	// }

	// shopV1Group := app.Group("/v1/shop")
	// {
	// 	shopV1Group.Get("/:toy_id", )
	// 	shopV1Group.Post("/", )
	// 	shopV1Group.Post("/list", )
	// }

	// {
	// 	app.Post("/v1/register/form", )
	// 	app.Post("v1/login/form", )
	// }


    // Определяем обработчик для корневого маршрута
    // app.Get("/", func(c *fiber.Ctx) error {
	// 	var req Request
		
	// 	req.User = c.Query("user", "unknown")

	// 	respBody := &Response{
	// 		User: req.User,
	// 		Message: "Пользователь получен",
	// 	}

	// 	c.SaveFile() // сразу сохраняет

	// 	// c.FormFile("key") // содержимое файла
	// 	// c.FormValue("key") // содержимое полей

    //     return c.Status(fiber.StatusOK).JSON(respBody) 
    // })

    // Запускаем сервер на порту 3000
    app.Listen(
		fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port), 
		//":3000"
	)

}