// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    hashed_password,
    host
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, username, email, hashed_password, host, created_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	Host           bool   `json:"host"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.HashedPassword,
		arg.Host,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.HashedPassword,
		&i.Host,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
select id, username, email, hashed_password, host, created_at from users
where email = $1 limit 1
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.HashedPassword,
		&i.Host,
		&i.CreatedAt,
	)
	return i, err
}
