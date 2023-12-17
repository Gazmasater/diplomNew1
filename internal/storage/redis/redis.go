package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return &RedisClient{
		client: rdb,
	}
}

func (rc *RedisClient) Ping() (string, error) {
	ctx := context.Background()
	pong, err := rc.client.Ping(ctx).Result()
	if err != nil {
		return "", err
	}

	return pong, nil
}

func (rc *RedisClient) Set(key, value string) error {
	ctx := context.Background()
	err := rc.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, err
}
