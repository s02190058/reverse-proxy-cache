package youtube

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	video = "vi"

	defaultType  = "DEFAULT"
	mediumType   = "MEDIUM"
	highType     = "HIGH"
	standardType = "STANDARD"
	maxresType   = "MAXRES"

	shortDefaultType  = ""
	shortMediumType   = "mq"
	shortHighType     = "hq"
	shortStandardType = "sd"
	shortMaxresType   = "maxres"

	defaultJPG = "default.jpg"
)

var (
	ErrInvalidBaseURL  = errors.New("invalid base url")
	ErrInvalidType     = errors.New("invalid thumbnail type")
	ErrTimeoutExceeded = errors.New("timeout exceeded")
	ErrInternal        = errors.New("youtube client internal error")
)

// Client for downloading images from YouTube.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a thumbnail management object
func NewClient(baseURL string, timeout time.Duration) (*Client, error) {
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, ErrInvalidBaseURL
	}

	client := Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
	return &client, nil
}

// VideoThumbnail downloads a video thumbnail of the specified type.
// Although all video ids consists of exactly 11 symbols of the form -, 0-9, A-Z, _, a-z
// this is not documented in their API, so we can't validate it.
func (c *Client) VideoThumbnail(ctx context.Context, videoID string, typ string) ([]byte, error) {
	shortType, err := shortVideoThumbnailType(typ)
	if err != nil {
		return nil, err
	}

	u, err := videoThumbnailURL(c.baseURL, videoID, shortType)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, wrapWithInternalError(err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return nil, ErrTimeoutExceeded
		}
		return nil, wrapWithInternalError(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, wrapWithInternalError(err)
	}
	return image, nil
}

// shortVideoThumbnailType matches human-readable type with the short type.
func shortVideoThumbnailType(typ string) (string, error) {
	switch typ {
	case defaultType:
		return shortDefaultType, nil
	case mediumType:
		return shortMediumType, nil
	case highType:
		return shortHighType, nil
	case standardType:
		return shortStandardType, nil
	case maxresType:
		return shortMaxresType, nil
	default:
		return "", ErrInvalidType
	}
}

// videoThumbnailPath generates a path along which the thumbnail is located
func videoThumbnailURL(baseURL string, videoID string, shortType string) (string, error) {
	path, err := url.JoinPath(baseURL, video, videoID, shortType+defaultJPG)
	if err != nil {
		return "", wrapWithInternalError(err)
	}
	return path, nil
}

func wrapWithInternalError(err error) error {
	return fmt.Errorf("%w: %w", ErrInternal, err)
}
