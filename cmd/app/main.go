package main

import (
	"flag"
	"log"

	"github.com/s02190058/reverse-proxy-cache/internal/app"
	"github.com/s02190058/reverse-proxy-cache/internal/config"
)

var configPath = flag.String("config", "./configs/main.yml", "path to config file")

// main is a service entrypoint.
func main() {
	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatalf("unable to read config: %v", err)
	}

	app.Run(cfg)
}
