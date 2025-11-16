package main

import (
	"service/internal/config"
	"service/internal/service"
	"service/internal/service/handlers"
	"service/internal/service/middlewares"
	"service/internal/storage/postgres"

	"fmt"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	storage, err := postgres.New(&cnf.Postgres)
	if err != nil {
		panic(err.Error())
	}

	application := &service.Application{
		Cnf: cnf,
		Log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Storage: storage,
		Validator: validator.New(),
	}

	// настройка middleware
	app.Use(logger.New(logger.Config{
		Format: "${blue}[${time}]${reset} ${cyan}${ip}:${port}${reset} ${method} ${green}${path}${reset} ${status} ${magenta}${latency}${reset} ${white}${reqHeader:User-Agent}${reset}\n",
	}))
	app.Use(recover.New())

	app.Static("/upload", cnf.Server.Prefix_upload)

	toysV1Group := app.Group("/v1/toys");
	toysV1Group.Use(middlewares.AuthMiddleware(application))
	{
		toysV1Group.Post("/", handlers.CreateToy(application))
		toysV1Group.Put("/", handlers.UpdateToy(application))
		toysV1Group.Post("/list", handlers.GetToysList(application))
		toysV1Group.Patch("/:toy_id", handlers.UpdateToyStatus(application))
		toysV1Group.Delete("/:toy_id", handlers.DeleteToy(application))
		toysV1Group.Get("/:toy_id", handlers.GetToy(application))
	}

	exchangeV1Group := app.Group("/v1/exchange")
	exchangeV1Group.Use(middlewares.AuthMiddleware(application))
	{
		exchangeV1Group.Post("/", handlers.CreateExchange(application))
		exchangeV1Group.Get("/:exchange_id", handlers.GetExchange(application))
		exchangeV1Group.Patch("/:exchange_id", handlers.PatchExchange(application))
		exchangeV1Group.Post("/list", handlers.GetExchangeList(application))
	}

	{
		app.Post("/v1/register", handlers.Register((application)))
		app.Post("v1/login", handlers.Login(application))
	}

    // Запускаем сервер на порту 3000
    app.Listen(fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port))

}