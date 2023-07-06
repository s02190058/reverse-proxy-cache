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
	"github.com/s02190058/reverse-proxy-cache/pkg/redis"
	"github.com/s02190058/reverse-proxy-cache/pkg/youtube"
)

// Run creates dependencies and launches the service.
func Run(cfg *config.Config) {
	rdb, err := redis.NewClient(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.DialTimeout,
		cfg.Redis.ConnTimeout,
	)
	if err != nil {
		log.Fatal(err)
	}

	thumbnailCache := cache.NewThumbnailCache(rdb, cfg.Redis.TTL)

	u := url.URL{
		Scheme: cfg.Youtube.Scheme,
		Host:   cfg.Youtube.Host,
	}
	yotubeClient, err := youtube.NewClient(u.String(), cfg.Youtube.Timeout)
	if err != nil {
		log.Fatal(err)
	}

	thumbnailService := service.NewThumbnailService(thumbnailCache, yotubeClient)

	grpcServer := grpc.NewServer()

	v1.RegisterThumbnailHandlers(grpcServer, thumbnailService)

	if err = grpcServer.Start(cfg.GRPC.Port); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		_ = sig
	case err = <-grpcServer.Notify():
		//
	}

	grpcServer.Stop()
}
