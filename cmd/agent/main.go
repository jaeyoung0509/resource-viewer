package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaeyoung050/resource-checker/internal/agent"
	"github.com/jaeyoung050/resource-checker/internal/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	nodeName, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v", err)
	}

	cfg := config.LoadAgentConfig()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := agent.Run(ctx, cfg, nodeName); err != nil {
		log.Fatalf("agent stopped: %v", err)
	}
}
