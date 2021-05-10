package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func New(addr, password string) *Redis {
	r := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &Redis{client: r, ctx: context.Background()}
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *Redis) SetEx(key, value string, expiration int64) error {
	return r.client.SetEX(r.ctx, key, value, time.Duration(expiration*int64(time.Second))).Err()
}

func (r *Redis) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
