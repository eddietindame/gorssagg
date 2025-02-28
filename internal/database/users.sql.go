// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, email, password, api_key)
VALUES ($1, $2, $3, $4, $5, $6,
  encode(sha256(random()::text::bytea), 'hex')
)
RETURNING id, created_at, updated_at, api_key, password, username, email
`

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Email     string
	Password  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Username,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ApiKey,
		&i.Password,
		&i.Username,
		&i.Email,
	)
	return i, err
}

const getUserByApiKey = `-- name: GetUserByApiKey :one
SELECT id, created_at, updated_at, api_key, password, username, email FROM users WHERE api_key = $1
`

func (q *Queries) GetUserByApiKey(ctx context.Context, apiKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByApiKey, apiKey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ApiKey,
		&i.Password,
		&i.Username,
		&i.Email,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, api_key, password, username, email FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ApiKey,
		&i.Password,
		&i.Username,
		&i.Email,
	)
	return i, err
}

const getUserPassword = `-- name: GetUserPassword :one
SELECT password FROM users WHERE username = $1
`

func (q *Queries) GetUserPassword(ctx context.Context, username string) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserPassword, username)
	var password string
	err := row.Scan(&password)
	return password, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users SET password = $1 where email = $2
`

type UpdateUserPasswordParams struct {
	Password string
	Email    string
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.Password, arg.Email)
	return err
}
