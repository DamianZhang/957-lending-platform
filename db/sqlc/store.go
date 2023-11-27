package db

import "github.com/jackc/pgx/v5/pgxpool"

// Store defines all functions to execute SQL queries and transactions
type Store interface {
	Querier
}

// PostgresStore provides all functions to execute SQL queries and transactions
type PostgresStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new Store
func NewStore(connPool *pgxpool.Pool) Store {
	return &PostgresStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
