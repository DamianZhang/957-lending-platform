package cache

import (
	"log"
	"os"
	"testing"

	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/redis/go-redis/v9"
)

var testCacher Cacher

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       0,
	})

	testCacher = NewCacher(redis)

	os.Exit(m.Run())
}
