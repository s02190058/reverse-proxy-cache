package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrBadConnection = errors.New("connection cannot be established")
)

// NewClient creates a client for Redis.
func NewClient(
	host, port string,
	password string,
	db int,
	dialTimeout time.Duration,
	connTimeout time.Duration,
) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:                  addr(host, port),
		Password:              password,
		DB:                    db,
		DialTimeout:           dialTimeout,
		ContextTimeoutEnabled: true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBadConnection, err)
	}

	return client, nil
}

// addr returns a host:port Redis address.
func addr(host, port string) string {
	return host + ":" + port
}
