package v1

import (
	"context"

	thumbnailpb "github.com/s02190058/reverse-proxy-cache/gen/go/thumbnail/v1"
)

// ThumbnailService for mocking.
type ThumbnailService interface {
	Download(ctx context.Context, videoID string, typ string) ([]byte, error)
}

// ThumbnailHandler is an implementation of thumbnailpb.ThumbnailServiceServer.
type ThumbnailHandler struct {
	thumbnailpb.UnimplementedThumbnailServiceServer

	service ThumbnailService
}

// Download handler implementation.
func (h *ThumbnailHandler) Download(
	ctx context.Context,
	req *thumbnailpb.DownloadThumbnailRequest,
) (*thumbnailpb.DownloadThumbnailResponse, error) {
	image, err := h.service.Download(ctx, req.VideoID, req.Type.String())
	if err != nil {
		//
	}

	return &thumbnailpb.DownloadThumbnailResponse{
		Image: image,
	}, nil
}
