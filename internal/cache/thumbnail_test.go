package cache_test

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/s02190058/reverse-proxy-cache/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestThumbnailCache_Get(t *testing.T) {
	t.Parallel()

	db, mock := redismock.NewClientMock()

	thumbnailCache := cache.NewThumbnailCache(db, 0)

	type mockBehavior func(clientMock redismock.ClientMock)

	testCases := []struct {
		name         string
		videoID      string
		videoType    string
		mockBehavior mockBehavior
		res          []byte
		err          error
	}{
		{
			name:      "cache hit",
			videoID:   "1",
			videoType: "type",
			mockBehavior: func(clientMock redismock.ClientMock) {
				clientMock.ExpectGet("thumbnail:1:type").SetVal("cache hit")
			},
			res: []byte("cache hit"),
			err: nil,
		},
		{
			name:      "cache miss",
			videoID:   "2",
			videoType: "type",
			mockBehavior: func(clientMock redismock.ClientMock) {
				clientMock.ExpectGet("thumbnail:2:type").RedisNil()
			},
			err: cache.ErrCacheMiss,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mock)
			res, err := thumbnailCache.Get(context.Background(), tc.videoID, tc.videoType)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.res, res)
		})
	}
}

func TestThumbnailCache_Set(t *testing.T) {
	t.Parallel()

	db, mock := redismock.NewClientMock()

	thumbnailCache := cache.NewThumbnailCache(db, 0)

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
				clientMock.ExpectSet("thumbnail:id:type", image, 0).SetVal("OK")
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.image)
			err := thumbnailCache.Set(context.Background(), tc.videoID, tc.videoType, tc.image)
			assert.Equal(t, tc.err, err)
		})
	}
}
