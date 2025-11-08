package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"service/internal/config"
	"service/internal/models"

	_"github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
	cnf *config.ConfigPostgres
}

func New(cnf *config.ConfigPostgres) (*Postgres, error) {
	const op = "StoragePostgres.New"

	conn_str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cnf.Host,
		cnf.Port,
		cnf.Username,
		cnf.Password,
		cnf.DbName,
		cnf.Sslmode)

	db, err := sql.Open(cnf.Driver, conn_str)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(int(cnf.MaxOpenConns))
	db.SetMaxIdleConns(int(cnf.MaxIdleConns))
	db.SetConnMaxIdleTime(cnf.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), cnf.Timeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{
		db: db,
		cnf: cnf,
	}, nil
}


func (s *Postgres) InsertToy(newToy *models.Toy) (*models.Toy, error) {
	const op = "Postgres.InsertToy";

	stmt, err := s.db.Prepare(kInsertToy);

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err);
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout);
	defer cancel();

	var dbToy models.Toy
	var description, photoUrl sql.NullString
	err = stmt.QueryRowContext(
		ctx, 
		newToy.ToyId,
		newToy.UserId,
		newToy.Name,
		newToy.Description,
		newToy.IdempotencyToken,
		newToy.PhotoUrl,
		newToy.Status,
		newToy.CreatedAt,
		newToy.UpdatedAt,
	).Scan(
		&dbToy.ToyId,
		&dbToy.UserId,
		&dbToy.Name,
		&description,
		&newToy.IdempotencyToken,
		&photoUrl,
		&dbToy.Status,
		&dbToy.CreatedAt,
		&dbToy.UpdatedAt,
	);

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err);
	}

	if description.Valid {
		dbToy.Description = &description.String
	}

	if photoUrl.Valid {
		dbToy.PhotoUrl = &photoUrl.String
	}

	return &dbToy, nil;
}

func (s *Postgres) SelectToyById(toyId string, userId string) (*models.Toy, error) {
	const op = "Postgres.SelectToyById";
	
	stmt, err := s.db.Prepare(kSelectToyById);

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err);
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout);
	defer cancel();

	var dbToy models.Toy
	var description, photoUrl sql.NullString
	err = stmt.QueryRowContext(
		ctx,
		toyId,
		userId,
	).Scan(
		&dbToy.ToyId,
		&dbToy.UserId,
		&dbToy.Name,
		&description,
		&dbToy.IdempotencyToken,
		&photoUrl,
		&dbToy.Status,
		&dbToy.CreatedAt,
		&dbToy.UpdatedAt,
	);

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err);
	}


	if description.Valid {
		dbToy.Description = &description.String
	}

	if photoUrl.Valid {
		dbToy.PhotoUrl = &photoUrl.String
	}

	return &dbToy, nil;
}

func (s *Postgres) SelectToyByToken(token string) (*models.Toy, error) {
	const op = "Postgres.SelectToyByToken";
	
	stmt, err := s.db.Prepare(kSelectToyByToken);

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err);
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout);
	defer cancel();

	var dbToy models.Toy
	var description, photoUrl sql.NullString
	err = stmt.QueryRowContext(
		ctx,
		token,
	).Scan(
		&dbToy.ToyId,
		&dbToy.UserId,
		&dbToy.Name,
		&description,
		&dbToy.IdempotencyToken,
		&photoUrl,
		&dbToy.Status,
		&dbToy.CreatedAt,
		&dbToy.UpdatedAt,
	);

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err);
	}


	if description.Valid {
		dbToy.Description = &description.String
	}

	if photoUrl.Valid {
		dbToy.PhotoUrl = &photoUrl.String
	}

	return &dbToy, nil;
}
