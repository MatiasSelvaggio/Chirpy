// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, hashed_password
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const deleteUsers = `-- name: DeleteUsers :exec
DELETE FROM USERS
`

func (q *Queries) DeleteUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteUsers)
	return err
}

const getUsersByEmail = `-- name: GetUsersByEmail :one
SELECT id, created_at, updated_at, email, hashed_password FROM users
WHERE email = $1
`

func (q *Queries) GetUsersByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUsersByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const updateUsers = `-- name: UpdateUsers :one
UPDATE users
SET 
email = $1,
hashed_password = $2,
updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email, hashed_password
`

type UpdateUsersParams struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
}

func (q *Queries) UpdateUsers(ctx context.Context, arg UpdateUsersParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUsers, arg.Email, arg.HashedPassword, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}
