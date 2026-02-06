package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaeyoung050/resource-checker/internal/config"
	"github.com/jaeyoung050/resource-checker/internal/hub"
	"github.com/jaeyoung050/resource-checker/internal/observability"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := config.LoadHubConfig()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serviceName := config.EnvOrDefault("SERVICE_NAME", "resource-hub")
	otlpEndpoint := config.EnvOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	otlpInsecure := config.EnvBoolOrDefault("OTEL_EXPORTER_OTLP_INSECURE", false)
	shutdown, err := observability.Init(ctx, serviceName, otlpEndpoint, otlpInsecure)
	if err != nil {
		log.Fatalf("otel init failed: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutdownCtx)
	}()

	if err := hub.Run(ctx, cfg); err != nil {
		log.Fatalf("hub stopped: %v", err)
	}
}
