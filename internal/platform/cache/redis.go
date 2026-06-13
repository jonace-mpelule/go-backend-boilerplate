package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

type NoopCache struct{}

func NewRedis(ctx context.Context, url string) (Cache, error) {
	if url == "" {
		return NewNoop(), nil
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

func NewNoop() Cache {
	return &NoopCache{}
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (n *NoopCache) Set(context.Context, string, []byte, time.Duration) error {
	return ErrCacheDisabled
}

func (n *NoopCache) Get(context.Context, string) ([]byte, error) {
	return nil, ErrCacheDisabled
}

func (n *NoopCache) Delete(context.Context, string) error {
	return ErrCacheDisabled
}

func (n *NoopCache) Ping(context.Context) error {
	return ErrCacheDisabled
}

func (n *NoopCache) Close() error {
	return nil
}
