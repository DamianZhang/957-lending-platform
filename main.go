package main

import (
	"context"
	"log"

	"github.com/DamianZhang/957-lending-platform/api"
	db "github.com/DamianZhang/957-lending-platform/repository/db/sqlc"
	service "github.com/DamianZhang/957-lending-platform/service/impl"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DBSource          = "postgresql://root:secret@localhost:5432/957-lending-platform?sslmode=disable"
	HTTPServerAddress = "localhost:8000"
)

func main() {
	connPool, err := pgxpool.New(context.Background(), DBSource)
	if err != nil {
		log.Fatal("can not connect to DB:", err)
	}

	// repository
	repository := db.NewRepository(connPool)

	// service
	borrowerService := service.NewBorrowerServiceImpl(repository)

	// create server
	server, err := api.NewServer(borrowerService)
	if err != nil {
		log.Fatal("can not create server:", err)
	}

	// start server
	err = server.Start(HTTPServerAddress)
	if err != nil {
		log.Fatal("can not start server:", err)
	}
}
