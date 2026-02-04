package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaeyoung050/resource-checker/internal/config"
	"github.com/jaeyoung050/resource-checker/internal/observability"
	"github.com/jaeyoung050/resource-checker/internal/tracesvc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := config.LoadTraceSvcConfig()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serviceName := getenv("SERVICE_NAME", "trace-svc")
	shutdown, err := observability.Init(ctx, serviceName, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatalf("otel init failed: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutdownCtx)
	}()

	if err := tracesvc.Run(ctx, cfg); err != nil {
		log.Fatalf("trace-svc stopped: %v", err)
	}
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
