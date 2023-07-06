package app

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/s02190058/reverse-proxy-cache/internal/cache"
	"github.com/s02190058/reverse-proxy-cache/internal/config"
	"github.com/s02190058/reverse-proxy-cache/internal/service"
	v1 "github.com/s02190058/reverse-proxy-cache/internal/transport/grpc/v1"
	"github.com/s02190058/reverse-proxy-cache/pkg/grpc"
	"github.com/s02190058/reverse-proxy-cache/pkg/logger"
	"github.com/s02190058/reverse-proxy-cache/pkg/redis"
	"github.com/s02190058/reverse-proxy-cache/pkg/youtube"
	"golang.org/x/exp/slog"
)

// Run creates dependencies and launches the service.
func Run(cfg *config.Config) {
	l, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatal(err)
	}
	l.Info("logger initialized", slog.String("env", cfg.Env))

	rdb, err := redis.NewClient(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.DialTimeout,
		cfg.Redis.ConnTimeout,
	)
	if err != nil {
		l.Error("redis.NewClient", logger.Err(err))
		return
	}
	l.Info(
		"redis client initialized",
		slog.String("host", cfg.Redis.Host),
		slog.String("port", cfg.Redis.Port),
	)

	thumbnailCache := cache.NewThumbnailCache(rdb, cfg.Redis.TTL)

	u := url.URL{
		Scheme: cfg.Youtube.Scheme,
		Host:   cfg.Youtube.Host,
	}
	youtubeClient, err := youtube.NewClient(u.String(), cfg.Youtube.Timeout)
	if err != nil {
		l.Error("youtube.NewClient", logger.Err(err))
		return
	}
	l.Info("youtube client initialized")

	thumbnailService := service.NewThumbnailService(thumbnailCache, youtubeClient)

	grpcServer := grpc.NewServer(l)

	v1.RegisterThumbnailHandlers(grpcServer, l, thumbnailService)

	if err = grpcServer.Start(cfg.GRPC.Port); err != nil {
		l.Error("grpcServer.Start", logger.Err(err))
		return
	}
	l.Info("starting grpc server", slog.String("port", cfg.GRPC.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		l.Info("server interrupt", slog.Any("signal", sig))
	case err = <-grpcServer.Notify():
		l.Error("error occurred since grpc server started", logger.Err(err))
	}

	grpcServer.Stop()
}
