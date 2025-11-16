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

func runInTx[T any](db *sql.DB, fn func(tx *sql.Tx) (T, error)) (T, error) {
	var zero T

	tx, err := db.Begin()
	if err != nil {
		return zero, err
	}

	object, err := fn(tx)
	if err != nil {
		tx.Rollback()
		return zero, err
	}

	if err := tx.Commit(); err != nil {
		return zero, err
	}

	return object, nil
}

type Postgres struct {
	db  *sql.DB
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
		db:  db,
		cnf: cnf,
	}, nil
}

func (s *Postgres) UpdateToy(newToy *models.Toy) (*models.Toy, error) {
	const op = "Postgres.UpdateToy"

	stmt, err := s.db.Prepare(kUpdateToy)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

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
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if description.Valid {
		dbToy.Description = &description.String
	}

	if photoUrl.Valid {
		dbToy.PhotoUrl = &photoUrl.String
	}

	return &dbToy, nil
}

func (s *Postgres) InsertToy(newToy *models.Toy) (*models.Toy, error) {
	const op = "Postgres.InsertToy"

	stmt, err := s.db.Prepare(kInsertToy)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

	var dbToy models.Toy
	var description, photoUrl sql.NullString
	err = stmt.QueryRowContext(
		ctx,
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
	)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if description.Valid {
		dbToy.Description = &description.String
	}

	if photoUrl.Valid {
		dbToy.PhotoUrl = &photoUrl.String
	}

	return &dbToy, nil
}

func (s *Postgres) UpdateToyStatus(toyId string, userId string, status models.ToyStatus) (*models.Toy, error) {
	const op = "Postgres.UpdateToyStatus"

	stmt, err := s.db.Prepare(kUpdateToyStatus)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

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
	)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if description.Valid {
		dbToy.Description = &description.String
	}

	if photoUrl.Valid {
		dbToy.PhotoUrl = &photoUrl.String
	}

	return &dbToy, nil
}

func (s *Postgres) SelectToyById(toyId string) (*models.Toy, error) {
	const op = "Postgres.SelectToyById"

	stmt, err := s.db.Prepare(kSelectToyById)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

	var dbToy models.Toy
	var description, photoUrl sql.NullString
	err = stmt.QueryRowContext(
		ctx,
		toyId,
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
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if description.Valid {
		dbToy.Description = &description.String
	}

	if photoUrl.Valid {
		dbToy.PhotoUrl = &photoUrl.String
	}

	return &dbToy, nil
}

func (s *Postgres) SelectToyByUserId(toyId string, userId string) (*models.Toy, error) {
	const op = "Postgres.SelectToyByUserId"

	stmt, err := s.db.Prepare(kSelectToyByUserId)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

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
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if description.Valid {
		dbToy.Description = &description.String
	}

	if photoUrl.Valid {
		dbToy.PhotoUrl = &photoUrl.String
	}

	return &dbToy, nil
}

func (s *Postgres) SelectToyByToken(token string) (*models.Toy, error) {
	const op = "Postgres.SelectToyByToken"

	stmt, err := s.db.Prepare(kSelectToyByToken)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

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
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if description.Valid {
		dbToy.Description = &description.String
	}

	if photoUrl.Valid {
		dbToy.PhotoUrl = &photoUrl.String
	}

	return &dbToy, nil
}

func (s *Postgres) SelectToysList(query *models.QueryToys, cursor *string, limit int64) ([]models.Toy, *string, error) {
	const op = "Postgres.SelectToysList"

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
	queryParams = append(queryParams, limit+1)

	sqlQuery := fmt.Sprintf("%s%s", kSelectToysList, strings.Join(whereClauses, "\n"))

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		sqlQuery,
		queryParams...,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("2 %s, %w", op, err)
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

	var nextCursor *string = nil
	if int64(len(dbToys)) == limit+1 {
		nextCursor = &dbToys[len(dbToys)-1].ToyId
		dbToys = dbToys[:len(dbToys)-1]
	}

	return dbToys, nextCursor, nil
}

func (s *Postgres) insertExchange(exchange *models.Exchange) (*models.Exchange, error) {
	const op = "Postgres.insertExchange"

	stmt, err := s.db.Prepare(kInsertExchange)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

	var dbExchange models.Exchange
	err = stmt.QueryRowContext(
		ctx,
		exchange.SrcToyId,
		exchange.DstToyId,
		exchange.IdempotencyToken,
	).Scan(
		&dbExchange.ExchangeId,
		&dbExchange.SrcToyId,
		&dbExchange.DstToyId,
		&dbExchange.Status,
		&dbExchange.IdempotencyToken,
		&dbExchange.CreatedAt,
		&dbExchange.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return &dbExchange, nil
}

func (s *Postgres) insertExchangeDetails(exchangeDetails *models.ExchangeDetails) (*models.ExchangeDetails, error) {
	const op = "Postgres.insertExchangeDetails"
	stmt, err := s.db.Prepare(kInsertExchangeDetails)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

	var dbExchangeDetails models.ExchangeDetails
	err = stmt.QueryRowContext(
		ctx,
		exchangeDetails.ExchangeId,
		exchangeDetails.ToyId,
		exchangeDetails.UserId,
	).Scan(
		&dbExchangeDetails.ExchangeId,
		&dbExchangeDetails.ToyId,
		&dbExchangeDetails.UserId,
		&dbExchangeDetails.Status,
		&dbExchangeDetails.CreatedAt,
		&dbExchangeDetails.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return &dbExchangeDetails, nil
}

func (s *Postgres) InsertExchange(exchange *models.Exchange, exchangeDetails []models.ExchangeDetails) (*models.Exchange, error) {
	return runInTx(s.db, func(tx *sql.Tx) (*models.Exchange, error) {

		dbExchange, err := s.insertExchange(exchange)
		if err != nil {
			return nil, err
		}

		exchangeDetails[0].ExchangeId = dbExchange.ExchangeId
		exchangeDetails[1].ExchangeId = dbExchange.ExchangeId

		for _, details := range exchangeDetails {
			_, err = s.insertExchangeDetails(&details)
			if err != nil {
				return nil, err
			}
		}

		return dbExchange, nil
	})
}

func getExchangeParticipant(rows *sql.Rows) (*models.ExchangeParticipant, error) {
	var p models.ExchangeParticipant
	var toyDesc, toyPhoto, middleName sql.NullString

	err := rows.Scan(
		&p.ExchangeId,
		&p.ExchangeStatus,
		&p.IdempotencyToken,
		&p.ExchangeCreatedAt,
		&p.ExchangeUpdatedAt,

		&p.ToyId,
		&p.ToyName,
		&toyDesc,
		&toyPhoto,

		&p.UserId,
		&p.FirstName,
		&middleName,
		&p.LastName,

		&p.UserExchangeStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if toyDesc.Valid {
		p.ToyDescription = &toyDesc.String
	}
	if toyPhoto.Valid {
		p.ToyPhotoURL = &toyPhoto.String
	}
	if middleName.Valid {
		p.MiddleName = &middleName.String
	}

	return &p, nil
}

func (s *Postgres) SelectExchangeWithParticipants(exchangeId string) ([]models.ExchangeParticipant, error) {
	const op = "Postgres.SelectExchangeWithParticipants"

	stmt, err := s.db.Prepare(kSelectExchangeWithParticipants)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

	rows, err := stmt.QueryContext(ctx, exchangeId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	participants := make([]models.ExchangeParticipant, 0)
	for rows.Next() {
		p, err := getExchangeParticipant(rows)
		if err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}

		participants = append(participants, *p)
	}

	return participants, nil
}

func (s *Postgres) UpdateExchangeWithParticipants(exchangeId string, userId string, status models.ExchangeDetailsStatus) ([]models.ExchangeParticipant, error) {
	const op = "Postgres.UpdateExchangeWithParticipants"
	
	return runInTx(s.db, func(tx *sql.Tx) ([]models.ExchangeParticipant, error) {

		stmt, err := s.db.Prepare(kUpdateExchangeStatus)

		if err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		defer stmt.Close()

		ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
		defer cancel()

		_, err = stmt.ExecContext(ctx, exchangeId, userId, status)
		if err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}

		return s.SelectExchangeWithParticipants(exchangeId)
	})
}

func (s *Postgres) SelectExchangeList(query *models.QueryExchanges, userId string, cursor *string, limit int64) ([]models.ExchangeParticipant, *string, error) {
	const op = "Postgres.SelectExchangeList"

	var (
		whereClauses []string
		queryParams  []interface{}
		paramIndex   = 2
	)

	queryParams = append(queryParams, userId)

	if query.Statuses != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("AND e.status = ANY($%d)", paramIndex))
		queryParams = append(queryParams, pq.Array(query.Statuses))
		paramIndex++
	}

	if cursor != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("AND e.exchange_id >= $%d", paramIndex))
		queryParams = append(queryParams, *cursor)
		paramIndex++
	}

	whereClauses = append(whereClauses, "ORDER BY e.exchange_id, e.updated_at DESC")
	whereClauses = append(whereClauses, fmt.Sprintf("LIMIT $%d", paramIndex))
	queryParams = append(queryParams, limit+1)

	sqlQuery := fmt.Sprintf("%s%s", kSelectExchangeIdList, strings.Join(whereClauses, "\n"))

	ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		sqlQuery,
		queryParams...,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("2 %s, %w", op, err)
	}
	defer rows.Close()

	exchangeIds := make([]string, 0)
	for rows.Next() {
		var exchangeId string
		err := rows.Scan(&exchangeId)
		if err != nil {
			return nil, nil, fmt.Errorf("3 %s: %w", op, err)
		}

		exchangeIds = append(exchangeIds, exchangeId)
	}

	fmt.Println(exchangeIds)

	var nextCursor *string = nil
	if int64(len(exchangeIds)) == limit+1 {
		nextCursor = &exchangeIds[len(exchangeIds)-1]
		exchangeIds = exchangeIds[:len(exchangeIds)-1]
	}

	stmt, err := s.db.Prepare(kSelectExchangeList)

	if err != nil {
		return nil, nil, fmt.Errorf("%s, %w", op, err)
	}
	defer stmt.Close()

	ctx2, cancel2 := context.WithTimeout(context.Background(), s.cnf.Timeout)
	defer cancel2()

	rows2, err := stmt.QueryContext(ctx2, pq.Array(exchangeIds))
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows2.Close()

	participants := make([]models.ExchangeParticipant, 0)
	for rows2.Next() {
		fmt.Println("fffff")
		p, err := getExchangeParticipant(rows2)
		if err != nil {
			return nil, nil, fmt.Errorf("%s, %w", op, err)
		}

		participants = append(participants, *p)
	}

	fmt.Println(participants)

	return participants, nextCursor, nil
}

func (s *Postgres) CreateUser(user *models.User) (*models.User, error) {
    const op = "Postgres.CreateUser"

    stmt, err := s.db.Prepare(kInsertUser)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer stmt.Close()

    ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
    defer cancel()

    var dbUser models.User
    var middleName sql.NullString
    
    err = stmt.QueryRowContext(
        ctx,
        user.UserName.FirstName,
        user.UserName.MiddleName,
        user.UserName.LastName,
        user.Email,
        user.HashPassword,
    ).Scan(
        &dbUser.UserId,
        &dbUser.UserName.FirstName,
        &middleName,
        &dbUser.UserName.LastName,
        &dbUser.Email,
		&dbUser.HashPassword,
        &dbUser.CreatedAt,
        &dbUser.UpdatedAt,
    )

    if err == sql.ErrNoRows {
		return nil, nil
    }

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if middleName.Valid {
		dbUser.UserName.MiddleName = &middleName.String
	}

    return &dbUser, nil
}

func (s *Postgres) SelectUserByEmail(user *models.User) (*models.User, error) {
    const op = "Postgres.SelectUserByEmail"

    stmt, err := s.db.Prepare(kSelectUserByEmail)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer stmt.Close()

    ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
    defer cancel()

    var dbUser models.User
    var middleName sql.NullString
    
    err = stmt.QueryRowContext(
        ctx,
        user.Email,
    ).Scan(
        &dbUser.UserId,
        &dbUser.UserName.FirstName,
        &middleName,
        &dbUser.UserName.LastName,
        &dbUser.Email,
		&dbUser.HashPassword,
        &dbUser.CreatedAt,
        &dbUser.UpdatedAt,
    )

    if err == sql.ErrNoRows {
		return nil, nil
    }

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if middleName.Valid {
		dbUser.UserName.MiddleName = &middleName.String
	}

    return &dbUser, nil
}

func (s *Postgres) SelectUserById(user *models.User) (*models.User, error) {
    const op = "Postgres.SelectUserById"

    stmt, err := s.db.Prepare(kSelectUserById)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer stmt.Close()

    ctx, cancel := context.WithTimeout(context.Background(), s.cnf.Timeout)
    defer cancel()

    var dbUser models.User
    var middleName sql.NullString
    
    err = stmt.QueryRowContext(
        ctx,
        user.UserId,
    ).Scan(
        &dbUser.UserId,
        &dbUser.UserName.FirstName,
        &middleName,
        &dbUser.UserName.LastName,
        &dbUser.Email,
		&dbUser.HashPassword,
        &dbUser.CreatedAt,
        &dbUser.UpdatedAt,
    )

    if err == sql.ErrNoRows {
		return nil, nil
    }

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if middleName.Valid {
		dbUser.UserName.MiddleName = &middleName.String
	}

    return &dbUser, nil
}