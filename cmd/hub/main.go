package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/jaeyoung050/resource-checker/internal/config"
	"github.com/jaeyoung050/resource-checker/internal/hub"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := config.LoadHubConfig()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := hub.Run(ctx, cfg); err != nil {
		log.Fatalf("hub stopped: %v", err)
	}
}
