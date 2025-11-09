package postgres

const (
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
    		toy_id,
    		user_id, 
    		name,
    		description,
    		idempotency_token,
			photo_url,
    		status
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (idempotency_token) DO NOTHING
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
)

		// WHERE true
		// 	// AND (status = ANY($1))
		// 	// AND (user_id = ANY($2))
		// 	// AND (user_id != ALL($3))
		// 	// AND (toy_id >= $4)
		// ORDER BY toy_id, updated_at DESC
		// LIMIT $5;