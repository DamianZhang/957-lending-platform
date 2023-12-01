package main

import (
	"context"
	"log"

	"github.com/DamianZhang/957-lending-platform/api"
	"github.com/DamianZhang/957-lending-platform/cache"
	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
	service "github.com/DamianZhang/957-lending-platform/service/impl"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

	redis := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       0,
	})

	// store
	store := db.NewStore(connPool)

	// cacher
	cacher := cache.NewCacher(redis)

	// service
	borrowerService := service.NewBorrowerServiceImpl(store, cacher)

	// create server
	server := api.NewServer(config, borrowerService)

	// start server
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("can not start server:", err)
	}
}
