package postgres

const (
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
    		status,
    		created_at,
    		updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
    		updated_at;
	`
)