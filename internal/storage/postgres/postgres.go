package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"service/internal/config"
	"service/internal/models"

	"github.com/lib/pq"
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

func (s *Postgres) UpdateToy(newToy *models.Toy) (*models.Toy, error) {
	const op = "Postgres.UpdateToy";

	stmt, err := s.db.Prepare(kUpdateToy);

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
		newToy.PhotoUrl,
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

func (s *Postgres) UpdateToyStatus(toyId string, userId string, status models.ToyStatus) (*models.Toy, error) {
	const op = "Postgres.UpdateToyStatus";
	
	stmt, err := s.db.Prepare(kUpdateToyStatus);

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
		status,
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


func (s *Postgres) SelectToysList(query *models.QueryToys, cursor *string, limit int64) ([]models.Toy, *string, error) {
	const op = "Postgres.SelectToysList";

	var (
		whereClauses []string
		queryParams  []interface{} 
		paramIndex   = 1  
	)

	if query.Statuses != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("AND status = ANY($%d)", paramIndex))
		queryParams = append(queryParams, pq.Array(query.Statuses))
		paramIndex++
	}

	if query.ExcludeUserIds != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("AND user_id != ALL($%d)", paramIndex))
		queryParams = append(queryParams, pq.Array(query.ExcludeUserIds))
		paramIndex++
	}

	if query.UserIds != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("AND user_id = ANY($%d)", paramIndex))
		queryParams = append(queryParams, pq.Array(query.UserIds))
		paramIndex++
	}

	if cursor != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("AND toy_id >= $%d", paramIndex))
		queryParams = append(queryParams, *cursor)
		paramIndex++
	}

	whereClauses = append(whereClauses, "ORDER BY toy_id, updated_at DESC")
	whereClauses = append(whereClauses, fmt.Sprintf("LIMIT $%d", paramIndex))
	queryParams = append(queryParams, limit + 1)

	sqlQuery := fmt.Sprintf("%s%s", kSelectToysList, strings.Join(whereClauses, "\n"))

	fmt.Println(sqlQuery)

	// stmt, err := s.db.Prepare(kSelectToysList);
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("1 %s, %w", op, err);
	// }

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout);
	defer cancel();

	rows, err := s.db.QueryContext(
		ctx,
		sqlQuery,
		queryParams...
	)

	if err != nil {
		return nil, nil, fmt.Errorf("2 %s, %w", op, err);
	}
	defer rows.Close()

	dbToys := make([]models.Toy, 0)

	for rows.Next() {
		var toy models.Toy
        var description, photoUrl sql.NullString
        
        err := rows.Scan(
			&toy.ToyId,
			&toy.UserId,
			&toy.Name,
			&description,
			&toy.IdempotencyToken,
			&photoUrl,
			&toy.Status,
			&toy.CreatedAt,
			&toy.UpdatedAt,
        )
        if err != nil {
            return nil, nil, fmt.Errorf("3 %s: %w", op, err)
        }
        
        if description.Valid {
            toy.Description = &description.String
        }
        if photoUrl.Valid {
            toy.PhotoUrl = &photoUrl.String
        }
        
        dbToys = append(dbToys, toy)
	}

	var next_cursor *string = nil
	if int64(len(dbToys)) == limit + 1 {
		next_cursor = &dbToys[len(dbToys) - 1].ToyId
		dbToys = dbToys[:len(dbToys) - 1]
	}

	return dbToys, next_cursor, nil;
}