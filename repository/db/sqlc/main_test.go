package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dbSource = "postgresql://admin:admin@localhost:5432/957-lending-platform?sslmode=disable"

var testRepository Repository

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	testRepository = NewRepository(connPool)

	os.Exit(m.Run())
}
