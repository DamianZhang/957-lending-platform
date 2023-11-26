package main

import (
	"context"
	"log"

	"github.com/DamianZhang/957-lending-platform/api"
	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
	service "github.com/DamianZhang/957-lending-platform/service/impl"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not connect to DB:", err)
	}

	// store
	store := db.NewStore(connPool)

	// service
	borrowerService := service.NewBorrowerServiceImpl(store)

	// create server
	server, err := api.NewServer(borrowerService)
	if err != nil {
		log.Fatal("can not create server:", err)
	}

	// start server
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("can not start server:", err)
	}
}
