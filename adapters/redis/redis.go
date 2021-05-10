package redis

import (
	"context"
	"fmt"
	"icfs_pg/env"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func New(host string, port int, password string) (*Redis, error) {
	ctx := context.Background()
	r := redis.NewClient(&redis.Options{
		Addr:     getAddr(host, port),
		Password: password,
		DB:       0,
	})
	err := r.Ping(ctx).Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to redis")
	}
	return &Redis{client: r, ctx: ctx}, nil
}

func getAddr(host string, port int) string {
	if env.DockerEnabled() {
		host = "datastore"
	}
	return fmt.Sprintf("%s:%d", host, port)
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
