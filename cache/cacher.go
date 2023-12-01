package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cacher interface {
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	GetSessionByID(ctx context.Context, id string) (Session, error)
}

type RedisCacher struct {
	*redis.Client
}

func NewCacher(redis *redis.Client) Cacher {
	return &RedisCacher{
		Client: redis,
	}
}
