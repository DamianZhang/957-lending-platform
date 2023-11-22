// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package db

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error)
	UpdateUserByEmail(ctx context.Context, arg UpdateUserByEmailParams) (User, error)
}

var _ Querier = (*Queries)(nil)
