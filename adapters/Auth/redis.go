package auth

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func NewService() *Redis {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &Redis{client: r, ctx: context.Background()}
}

func (r *Redis) Get(key string) string {
	return r.client.Get(r.ctx, key).String()
}

func (r *Redis) SetEx(key, value string, timeout int) error {
	return r.client.SetEX(r.ctx, key, value, time.Duration(timeout)*time.Second).Err()
}
