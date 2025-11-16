package service

import (
	"service/internal/config"
	"service/internal/models"

	"log/slog"

	"github.com/go-playground/validator/v10"
)

type Storage interface {
	// TOY
	InsertToy(newToy *models.Toy) (*models.Toy, error)
	SelectToyById(toyId string) (*models.Toy, error)
	SelectToyByUserId(toyId string, userId string) (*models.Toy, error)
	SelectToyByToken(token string) (*models.Toy, error)
	UpdateToyStatus(toyId string, userId string, status models.ToyStatus) (*models.Toy, error)
	UpdateToy(newToy *models.Toy) (*models.Toy, error)
	SelectToysList(query *models.QueryToys, cursor *string, limit int64) ([]models.Toy, *string, error)

	// EXCHANGE
	InsertExchange(exchange *models.Exchange, exchangeDetails []models.ExchangeDetails) (*models.Exchange, error)
	SelectExchangeWithParticipants(exchangeId string) ([]models.ExchangeParticipant, error)
	UpdateExchangeWithParticipants(exchangeId string, userId string, status models.ExchangeDetailsStatus) ([]models.ExchangeParticipant, error)
	SelectExchangeList(query *models.QueryExchanges, userId string, cursor *string, limit int64) ([]models.ExchangeParticipant, *string, error)

	// USER
	SelectUserById(user *models.User) (*models.User, error)
	SelectUserByEmail(user *models.User) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
}

type Application struct {
	Cnf *config.Config
	Storage Storage
	Log *slog.Logger
	Validator *validator.Validate
	//wg     sync.WaitGroup //updated
}