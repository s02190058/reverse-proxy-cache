package v1

import (
	"context"
	"errors"

	thumbnailpb "github.com/s02190058/reverse-proxy-cache/gen/go/thumbnail/v1"
	"github.com/s02190058/reverse-proxy-cache/internal/service"
	grpcserver "github.com/s02190058/reverse-proxy-cache/pkg/grpc"
	"github.com/s02190058/reverse-proxy-cache/pkg/logger"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ThumbnailService for mocking.
type ThumbnailService interface {
	Download(ctx context.Context, videoID string, typ string) ([]byte, error)
}

// ThumbnailHandler is an implementation of thumbnailpb.ThumbnailServiceServer.
type ThumbnailHandler struct {
	thumbnailpb.UnimplementedThumbnailServiceServer

	l       *slog.Logger
	service ThumbnailService
}

// RegisterThumbnailHandlers adds thumbnail handlers to the gRPC server.
func RegisterThumbnailHandlers(
	server *grpcserver.Server,
	l *slog.Logger,
	service ThumbnailService,
) {
	server.RegisterHandlers(func(s *grpc.Server) {
		thumbnailpb.RegisterThumbnailServiceServer(s, &ThumbnailHandler{
			l:       l,
			service: service,
		})
	})
}

// Download handler implementation.
func (h *ThumbnailHandler) Download(
	ctx context.Context,
	req *thumbnailpb.DownloadThumbnailRequest,
) (*thumbnailpb.DownloadThumbnailResponse, error) {
	image, err := h.service.Download(ctx, req.VideoID, req.Type.String())
	if err != nil {
		h.l.Error("ThumbnailService.Download", logger.Err(err))
		switch {
		// image wasn't cache, but response is valid
		case errors.Is(err, service.ErrBadCache):
			break
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	return &thumbnailpb.DownloadThumbnailResponse{
		Image: image,
	}, nil
}
