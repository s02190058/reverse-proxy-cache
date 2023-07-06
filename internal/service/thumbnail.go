package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/s02190058/reverse-proxy-cache/internal/cache"
)

var (
	ErrBadCache = errors.New("cannot cache a thumbnail")
	ErrInternal = errors.New("internal server error")
)

// ThumbnailCache for mocking.
type ThumbnailCache interface {
	Get(ctx context.Context, videoID string, typ string) ([]byte, error)
	Set(ctx context.Context, videoID string, typ string, image []byte) error
}

// YoutubeClient for mocking.
type YoutubeClient interface {
	VideoThumbnail(ctx context.Context, videoID string, typ string) ([]byte, error)
}

// ThumbnailService implements core logic of service.
type ThumbnailService struct {
	cache         ThumbnailCache
	youtubeClient YoutubeClient
}

// NewThumbnailService creates a thumbnail service.
func NewThumbnailService(cache ThumbnailCache, youtubeClient YoutubeClient) *ThumbnailService {
	return &ThumbnailService{
		cache:         cache,
		youtubeClient: youtubeClient,
	}
}

// Download gets a thumbnail for the video with a certain id.
func (s *ThumbnailService) Download(ctx context.Context, videoID string, typ string) ([]byte, error) {
	image, err := s.cache.Get(ctx, videoID, typ)
	// cache hit
	if err == nil {
		return image, nil
	}

	// internal error
	if !errors.Is(err, cache.ErrCacheMiss) {
		return nil, withInternalError(err)
	}

	//cache miss
	image, err = s.youtubeClient.VideoThumbnail(ctx, videoID, typ)
	if err != nil {
		return nil, withInternalError(err)
	}

	if err = s.cache.Set(ctx, videoID, typ, image); err != nil {
		return image, withBadCacheError(err)
	}

	return image, nil
}

// withBadCacheError wraps an internal error.
func withBadCacheError(err error) error {
	return fmt.Errorf("%w: %w", ErrBadCache, err)
}

// withInternalError wraps an internal error.
func withInternalError(err error) error {
	return fmt.Errorf("%w: %w", ErrInternal, err)
}
