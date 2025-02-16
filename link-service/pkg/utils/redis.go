package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

var Ctx = context.Background()

var Rdb *redis.Client

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) *RedisAdapter {
	return &RedisAdapter{client: client}
}

func (r *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisAdapter) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func CreateRedisConn() *redis.Client {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	})

	return Rdb
}

func GetRedisConn() (context.Context, *redis.Client) {
	return Ctx, Rdb
}
