package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedis(
	url string,
) *RedisCache {
	opt, _ := redis.ParseURL(url)

	client := redis.NewClient(opt)

	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) Set(
	ctx context.Context,
	key string,
	value any,
	ttl int,
) error {
	return r.client.Set(
		ctx,
		key,
		value,
		time.Duration(ttl)*time.Second,
	).Err()
}

func (r *RedisCache) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(
	ctx context.Context,
	key string,
) error {
	return r.client.Del(ctx, key).Err()
}
