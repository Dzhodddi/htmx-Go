package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"project/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewRedisStorage(cacheRedis *redis.Client) Storage {
	return Storage{
		Users: &UserStore{cacheRedis: cacheRedis},
	}
}
