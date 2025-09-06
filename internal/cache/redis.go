package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

var ErrKeyNotExist = fmt.Errorf("cache key not exists")

type RedisConfig struct {
	DB redis.UniversalClient `validate:"required"`
}

type Redis struct {
	client redis.UniversalClient
}

func NewRedis(conf RedisConfig) (*Redis, error) {
	err := validator.New().Struct(conf)
	if err != nil {
		return nil, fmt.Errorf("error validate cache redis: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conf.DB.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("redis is up and running")

	return &Redis{client: conf.DB}, nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, string(jsonValue), expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// GetAs method for type-safe retrieval
func (r *Redis) GetAs(ctx context.Context, key string, out interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("%w: key %s", ErrKeyNotExist, key)
	}

	if err != nil {
		return fmt.Errorf("error occurred on redis get: %w", err)
	}

	return json.Unmarshal([]byte(val), out)
}

// SetExp method for your preference
func (r *Redis) SetExp(ctx context.Context, key string, inValue interface{}, expireDur time.Duration) error {
	val, err := json.Marshal(inValue)
	if err != nil {
		return fmt.Errorf("cannot marshal json value: %w", err)
	}

	return r.client.Set(ctx, key, val, expireDur).Err()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeletePattern deletes all keys matching the pattern (from your peer's implementation)
func (r *Redis) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("error getting keys with pattern %s: %w", pattern, err)
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	return nil
}
