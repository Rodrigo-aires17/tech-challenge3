package internal

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implementa Cache usando o ElastiCache Redis provisionado via Terraform.
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client: client, ttl: ttl}
}

func (c *RedisCache) Get(ctx context.Context, key string) (bool, bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	return val == "1", true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value bool) error {
	v := "0"
	if value {
		v = "1"
	}
	return c.client.Set(ctx, key, v, c.ttl).Err()
}
