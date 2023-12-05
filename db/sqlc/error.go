package db

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
)

var (
	ErrConnDone       = sql.ErrConnDone
	ErrRecordNotFound = pgx.ErrNoRows
)
