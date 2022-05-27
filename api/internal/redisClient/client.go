package vredis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type client struct {
	rdb *redis.Client
}

func New(ctx context.Context, endpoint string, port string, password string) (Client, error) {
	c := client{}
	redisAddress := fmt.Sprintf("%s:%s", endpoint, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: password,
	})

	pingResp := rdb.Ping(ctx)
	if pingResp.Err() != nil {
		return nil, pingResp.Err()
	}

	c.rdb = rdb

	return &c, nil
}

func (c *client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *client) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.rdb.Get(ctx, key)
}

func (c *client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.rdb.Set(ctx, key, value, expiration)
}
