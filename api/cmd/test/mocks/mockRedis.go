package mocks

import (
	"context"
	"time"

	vredis "github.com/VamaSingapore/vama-api/internal/redisClient"
	"github.com/go-redis/redis/v8"
)

type MockRedisClient struct{}

func NewRedisClient() vredis.Client {
	return &MockRedisClient{}
}

func (muc *MockRedisClient) Ping(ctx context.Context) error {
	return nil
}

func (muc *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return nil
}

func (muc *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return nil
}
