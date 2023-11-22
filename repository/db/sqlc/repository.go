package db

import "github.com/jackc/pgx/v5/pgxpool"

// Repository defines all functions to execute SQL queries and transactions
type Repository interface {
	Querier
}

// PostgresRepository provides all functions to execute SQL queries and transactions
type PostgresRepository struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewRepository creates a new Repository
func NewRepository(connPool *pgxpool.Pool) Repository {
	return &PostgresRepository{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
