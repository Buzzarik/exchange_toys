package service

import (
	"service/internal/config"
	"service/internal/models"
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io"
	"log/slog"
	// "net/http"
	_ "sync"
	// "time"

	"github.com/go-playground/validator/v10"
)

type Storage interface {
	// InsertUser(user *models.User) (*models.User error)
	// SelectUserByEmail(email string) (*models.User, error)
	// SelectUserById(user_id string) (*models.User, error)

	InsertToy(newToy *models.Toy) (*models.Toy, error)
	SelectToyById(toyId string, userId string) (*models.Toy, error)
	SelectToyByToken(token string) (*models.Toy, error)
	UpdateToyStatus(toyId string, userId string, status models.ToyStatus) (*models.Toy, error)
	UpdateToy(newToy *models.Toy) (*models.Toy, error)
	SelectToysList(query *models.QueryToys, cursor *string, limit int64) ([]models.Toy, *string, error)
	// DeleteToy(
	// 	toy_id string, 
	// 	user_id string,
	// ) (error)

	// SelectExchangeById(exchange_id string) (*models.Exchange, error)
	// SelectExchangesByUserId(
	// 	user_id string, 
	// 	cursor *string, 
	// 	limit int64,
	// ) ([]*models.Exchange, *string, error)
	// UpdateExchangeStatus(
	// 	exchange_id string, 
	// 	user_id string, 
	// 	status string,
	// )
}

type Application struct {
	Cnf *config.Config
	Storage Storage
	Log *slog.Logger
	Validator *validator.Validate
	//wg     sync.WaitGroup //updated
}