package postgres

const (
// TOY
	kSelectToysList = 
	`
		SELECT 
			toy_id,
			user_id,
			name,
			description,
			idempotency_token,
			photo_url,
			status,
			created_at,
			updated_at
		FROM toys
		WHERE true
	`

	kSelectToyById = 
	`
		SELECT 
			toy_id,
			user_id,
			name,
			description,
			idempotency_token,
			photo_url,
			status,
			created_at,
			updated_at
		FROM toys
		WHERE true
			AND toy_id = $1
			AND status != 'removed'
		;
	`

	kSelectToyByUserId = 
	`
		SELECT 
			toy_id,
			user_id,
			name,
			description,
			idempotency_token,
			photo_url,
			status,
			created_at,
			updated_at
		FROM toys
		WHERE true
			AND toy_id = $1
			AND user_id = $2
			AND status != 'removed'
		;
	`

	kSelectToyByToken = 
	`
		SELECT 
			toy_id,
			user_id,
			name,
			description,
			idempotency_token,
			photo_url,
			status,
			created_at,
			updated_at
		FROM toys
		WHERE true
			AND idempotency_token = $1
		;
	`

	kInsertToy =
	`
		INSERT INTO toys (
    		user_id, 
    		name,
    		description,
    		idempotency_token,
			photo_url,
    		status
		) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (idempotency_token)
		DO UPDATE SET
        	idempotency_token = EXCLUDED.idempotency_token
		RETURNING 
			toy_id,
    		user_id, 
    		name,
    		description,
    		idempotency_token,
			photo_url,
    		status,
    		created_at,
    		updated_at
		;
	`

	kUpdateToyStatus = 
	`
		UPDATE toys 
		SET 
    		status = $3,
    		updated_at = NOW()
		WHERE true
			AND toy_id = $1
    		AND user_id = $2
			AND status != 'removed'
		RETURNING 
			toy_id,
    		user_id, 
    		name,
    		description,
    		idempotency_token,
			photo_url,
    		status,
    		created_at,
    		updated_at
		;
	`

	kUpdateToy =
	`
		UPDATE toys 
		SET 
			name = $3,
			description = $4,
			photo_url = COALESCE($5, photo_url),
			updated_at = NOW()
		WHERE true
			AND toy_id = $1 
			AND user_id = $2
			AND status != 'removed'
		RETURNING 
			toy_id,
    		user_id, 
    		name,
    		description,
    		idempotency_token,
			photo_url,
    		status,
    		created_at,
    		updated_at
		;
	`

// EXCHANGE
	kInsertExchange = 
	`
		INSERT INTO exchange 
			(src_toy_id, dst_toy_id, idempotency_token)
		VALUES 
		($1, $2, $3)
		ON CONFLICT (idempotency_token)
		DO UPDATE SET
        	idempotency_token = EXCLUDED.idempotency_token
		RETURNING 
			exchange_id, 
			src_toy_id, 
			dst_toy_id, 
			status, 
			idempotency_token, 
			created_at,
			updated_at
		;
	`

	kInsertExchangeDetails = 
	`
		INSERT INTO exchange_details 
			(exchange_id, toy_id, user_id)
		VALUES 
			($1, $2, $3)
		ON CONFLICT (exchange_id, toy_id, user_id)
		DO UPDATE SET
        	exchange_id = EXCLUDED.exchange_id
		RETURNING 
			exchange_id, 
			toy_id, 
			user_id, 
			status, 
			created_at,
			updated_at
		;
	`

	kSelectExchangeWithParticipants = 
	`
        SELECT 
            e.exchange_id,
            e.status AS exchange_status,
            e.idempotency_token,
            e.created_at AS exchange_created_at,
            e.updated_at AS exchange_updated_at,
            
            t.toy_id,
            t.name AS toy_name,
            t.description AS toy_description,
            t.photo_url AS toy_photo_url,
            
            u.user_id,
            u.first_name,
            u.middle_name,
            u.last_name,
            
            ed.status AS user_exchange_status

        FROM exchange e
        INNER JOIN exchange_details ed ON e.exchange_id = ed.exchange_id
        INNER JOIN toys t ON ed.toy_id = t.toy_id
        INNER JOIN users u ON ed.user_id = u.user_id
        WHERE e.exchange_id = $1
        ORDER BY e.exchange_id, u.user_id
	`

	kUpdateExchangeStatus = 
	`
		UPDATE exchange_details
		SET 
			status = $3
		WHERE true
			AND user_id = $2
			AND exchange_id = $1;
	`

	kSelectExchangeIdList = 
	`
		SELECT 
			e.exchange_id
		FROM exchange e
        INNER JOIN exchange_details ed ON e.exchange_id = ed.exchange_id
		WHERE true
			AND ed.user_id = $1
	`

	kSelectExchangeList = 
	`
        SELECT 
            e.exchange_id,
            e.status AS exchange_status,
            e.idempotency_token,
            e.created_at AS exchange_created_at,
            e.updated_at AS exchange_updated_at,
            
            t.toy_id,
            t.name AS toy_name,
            t.description AS toy_description,
            t.photo_url AS toy_photo_url,
            
            u.user_id,
            u.first_name,
            u.middle_name,
            u.last_name,
            
            ed.status AS user_exchange_status

        FROM exchange e
        INNER JOIN exchange_details ed ON e.exchange_id = ed.exchange_id
        INNER JOIN toys t ON ed.toy_id = t.toy_id
        INNER JOIN users u ON ed.user_id = u.user_id
        WHERE e.exchange_id = ANY($1);
	`

	kInsertUser = 
	`    
		INSERT INTO users 
			(first_name, middle_name, last_name, email, password_hash)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (email) DO NOTHING
        RETURNING user_id, first_name, middle_name, last_name, email, password_hash, created_at, updated_at
	`

	kSelectUserByEmail = 
	`
		SELECT user_id, first_name, middle_name, last_name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	kSelectUserById = 
	`
		SELECT user_id, first_name, middle_name, last_name, email, password_hash, created_at, updated_at
		FROM users
		WHERE user_id = $1
	`
)