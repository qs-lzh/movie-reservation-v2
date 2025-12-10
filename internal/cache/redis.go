package cache

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(url string) *RedisCache {
	client := redis.NewClient(
		&redis.Options{
			Addr:     url,
			Password: "",
			DB:       0,
		},
	)
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisCache) Get(key string, dest any) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *RedisCache) SetBool(key string, value bool) error {
	strValue := "false"
	if value {
		strValue = "true"
	}
	return r.client.Set(ctx, key, strValue, 5*time.Minute).Err()
}

func (r *RedisCache) GetBool(key string) (value bool, err error) {
	value, err = r.client.Get(ctx, key).Bool()
	if err != nil {
		return false, err
	}
	return value, nil
}
