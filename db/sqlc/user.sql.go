// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: user.sql

package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "users" (
  "email",
  "hashed_password",
  "line_id",
  "nickname",
  "role"
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id, created_at, updated_at, email, hashed_password, line_id, nickname, is_email_verified, role
`

type CreateUserParams struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	LineID         string `json:"line_id"`
	Nickname       string `json:"nickname"`
	Role           string `json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.HashedPassword,
		arg.LineID,
		arg.Nickname,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.LineID,
		&i.Nickname,
		&i.IsEmailVerified,
		&i.Role,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, line_id, nickname, is_email_verified, role FROM "users"
WHERE "email" = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.LineID,
		&i.Nickname,
		&i.IsEmailVerified,
		&i.Role,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, created_at, updated_at, email, hashed_password, line_id, nickname, is_email_verified, role FROM "users"
LIMIT $1
OFFSET $2
`

type GetUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Email,
			&i.HashedPassword,
			&i.LineID,
			&i.Nickname,
			&i.IsEmailVerified,
			&i.Role,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserByEmail = `-- name: UpdateUserByEmail :one
UPDATE "users"
SET
  "hashed_password" = COALESCE($1, "hashed_password"),
  "line_id" = COALESCE($2, "line_id"),
  "nickname" = COALESCE($3, "nickname"),
  "is_email_verified" = COALESCE($4, "is_email_verified"),
  "updated_at" = COALESCE($5, "updated_at")
WHERE
  "email" = $6
RETURNING id, created_at, updated_at, email, hashed_password, line_id, nickname, is_email_verified, role
`

type UpdateUserByEmailParams struct {
	HashedPassword  pgtype.Text `json:"hashed_password"`
	LineID          pgtype.Text `json:"line_id"`
	Nickname        pgtype.Text `json:"nickname"`
	IsEmailVerified pgtype.Bool `json:"is_email_verified"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Email           string      `json:"email"`
}

func (q *Queries) UpdateUserByEmail(ctx context.Context, arg UpdateUserByEmailParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserByEmail,
		arg.HashedPassword,
		arg.LineID,
		arg.Nickname,
		arg.IsEmailVerified,
		arg.UpdatedAt,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.LineID,
		&i.Nickname,
		&i.IsEmailVerified,
		&i.Role,
	)
	return i, err
}
