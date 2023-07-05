package youtube_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/s02190058/reverse-proxy-cache/pkg/youtube"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		baseURL string
		err     error
	}{
		{
			name:    "valid base url",
			baseURL: "https://golang.org",
			err:     nil,
		},
		{
			name:    "invalid base url",
			baseURL: "htp///golang,org",
			err:     youtube.ErrInvalidBaseURL,
		},
	}

	for _, tc := range testCases {
		// to avoid races
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := youtube.NewClient(tc.baseURL, 0)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestClient_VideoThumbnail(t *testing.T) {
	t.Parallel()

	videoID := "un6ZyFkqFKo"

	// mock external dependency
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp string
		switch r.URL.Path {
		case "/vi/" + videoID + "/default.jpg":
			resp = "default"
		case "/vi/" + videoID + "/mqdefault.jpg":
			resp = "medium"
		case "/vi/" + videoID + "/hqdefault.jpg":
			resp = "high"
		case "/vi/" + videoID + "/sddefault.jpg":
			resp = "standard"
		case "/vi/" + videoID + "/maxresdefault.jpg":
			resp = "maxres"
		default:
			// endpoint for timeout exceeded error
			time.Sleep(2 * time.Second)
		}

		if _, err := w.Write([]byte(resp)); err != nil {
			t.Fatalf("ResponseWriter.Write: %v", err)
		}
	}))

	client, err := youtube.NewClient(httpServer.URL, time.Second)
	if err != nil {
		t.Fatalf("youtube.NewClient: %v", err)
	}

	testCases := []struct {
		name    string
		videoID string
		typ     string
		res     []byte
		err     error
	}{
		{
			name:    "default type",
			videoID: videoID,
			typ:     "DEFAULT",
			res:     []byte("default"),
			err:     nil,
		},
		{
			name:    "medium type",
			videoID: videoID,
			typ:     "MEDIUM",
			res:     []byte("medium"),
			err:     nil,
		},
		{
			name:    "high type",
			videoID: videoID,
			typ:     "HIGH",
			res:     []byte("high"),
			err:     nil,
		},
		{
			name:    "standard type",
			videoID: videoID,
			typ:     "STANDARD",
			res:     []byte("standard"),
			err:     nil,
		},
		{
			name:    "maxres type",
			videoID: videoID,
			typ:     "MAXRES",
			res:     []byte("maxres"),
			err:     nil,
		},
		{
			name:    "invalid type",
			videoID: videoID,
			typ:     "INVALID",
			res:     nil,
			err:     youtube.ErrInvalidType,
		},
		{
			name:    "timeout exceeded error",
			videoID: "",
			typ:     "DEFAULT",
			res:     nil,
			err:     youtube.ErrTimeoutExceeded,
		},
	}

	for _, tc := range testCases {
		// to avoid races
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := client.VideoThumbnail(tc.videoID, tc.typ)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.res, res)
		})
	}
}
