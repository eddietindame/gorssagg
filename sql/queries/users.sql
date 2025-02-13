-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, password, api_key)
VALUES ($1, $2, $3, $4, $5,
  encode(sha256(random()::text::bytea), 'hex')
)
RETURNING *;

-- name: GetUserByApiKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: GetUserPassword :one
SELECT password FROM users WHERE username = $1;
