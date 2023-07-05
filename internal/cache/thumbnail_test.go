package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/s02190058/reverse-proxy-cache/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestThumbnailCache_Set(t *testing.T) {
	t.Parallel()

	db, mock := redismock.NewClientMock()

	ttl := 1 * time.Second

	thumbnailCache := cache.NewThumbnailCache(db, ttl)

	type mockBehavior func(clientMock redismock.ClientMock, image []byte)

	testCases := []struct {
		name         string
		videoID      string
		videoType    string
		image        []byte
		mockBehavior mockBehavior
		err          error
	}{
		{
			name:      "valid",
			videoID:   "id",
			videoType: "type",
			image:     []byte("valid"),
			mockBehavior: func(clientMock redismock.ClientMock, image []byte) {
				clientMock.ExpectSet("thumbnail:id:type", image, ttl).SetVal("OK")
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		// to avoid races
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.mockBehavior(mock, tc.image)
			err := thumbnailCache.Set(context.Background(), tc.videoID, tc.videoType, tc.image)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestThumbnailCache_Get(t *testing.T) {
	t.Parallel()

	db, mock := redismock.NewClientMock()

	ttl := 1 * time.Second

	thumbnailCache := cache.NewThumbnailCache(db, ttl)

	type mockBehavior func(clientMock redismock.ClientMock)

	testCases := []struct {
		name         string
		videoID      string
		videoType    string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name:      "valid",
			videoID:   "id",
			videoType: "type",
			mockBehavior: func(clientMock redismock.ClientMock) {
				clientMock.ExpectGet("thumbnail:id:type").RedisNil()
			},
			err: cache.ErrCacheMiss,
		},
	}

	for _, tc := range testCases {
		// to avoid races
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.mockBehavior(mock)
			_, err := thumbnailCache.Get(context.Background(), tc.videoID, tc.videoType)
			assert.Equal(t, tc.err, err)
		})
	}
}
