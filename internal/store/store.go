package store

import (
	"context"
	"fmt"
	"time"

	"github.com/eddietindame/gorssagg/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_HOST,
		Username: config.REDIS_USERNAME,
		Password: config.REDIS_PASSWORD,
		DB:       0,
	})
}

func StoreToken(ctx context.Context, token, email string, ttl time.Duration) error {
	err := RedisClient.Set(ctx, token, email, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to store token in Redis: %w", err)
	}
	return nil
}

func DeleteToken(ctx context.Context, token string) error {
	_, err := RedisClient.Del(ctx, token).Result()
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}

func GetEmailFromToken(ctx context.Context, token string) (string, error) {
	email, err := RedisClient.Get(ctx, token).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("token not found")
	} else if err != nil {
		return "", fmt.Errorf("error getting token from Redis: %w", err)
	}
	return email, nil
}
