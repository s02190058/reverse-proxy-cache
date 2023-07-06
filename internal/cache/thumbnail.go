package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInternal  = errors.New("redis cache internal error")
	ErrCacheMiss = errors.New("cache miss")
)

// ThumbnailCache is a redis cache for thumbnails.
type ThumbnailCache struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewThumbnailCache creates a redis cache.
func NewThumbnailCache(rdb *redis.Client, ttl time.Duration) *ThumbnailCache {
	return &ThumbnailCache{
		rdb: rdb,
		ttl: ttl,
	}
}

// Get checks whether the image bytes are stored in cache.
func (c *ThumbnailCache) Get(ctx context.Context, videoID string, typ string) ([]byte, error) {
	key := thumbnailKey(videoID, typ)
	image, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, withInternalError(err)
	}

	return image, nil
}

// Set saves an image bytes in cache for a certain time.
func (c *ThumbnailCache) Set(ctx context.Context, videoID string, typ string, image []byte) error {
	key := thumbnailKey(videoID, typ)
	if err := c.rdb.Set(ctx, key, image, c.ttl).Err(); err != nil {
		return withInternalError(err)
	}

	return nil
}

// generateThumbnailKey generates a "thumbnail:<video_id>:<type>" key.
func thumbnailKey(videoID string, typ string) string {
	return fmt.Sprintf("%s:%s:%s",
		"thumbnail",
		videoID,
		strings.ToLower(typ),
	)
}

// withInternalError wraps an internal error.
func withInternalError(err error) error {
	return fmt.Errorf("%w: %w", ErrInternal, err)
}
