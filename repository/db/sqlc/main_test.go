package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testRepository Repository

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	testRepository = NewRepository(connPool)

	os.Exit(m.Run())
}
