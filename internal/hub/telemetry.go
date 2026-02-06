package hub

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer        = otel.Tracer("resource-hub")
	meter         = otel.Meter("resource-hub")
	wsConnections metric.Int64Counter
	wsMessages    metric.Int64Counter
	apiRequests   metric.Int64Counter
	initOnce      sync.Once
	initErr       error
)

func initTelemetry() error {
	initOnce.Do(func() {
		var err error
		wsConnections, err = meter.Int64Counter(
			"hub_ws_connections",
			metric.WithDescription("WebSocket connections accepted by the hub"),
		)
		if err != nil {
			initErr = err
			return
		}
		wsMessages, err = meter.Int64Counter(
			"hub_ws_messages",
			metric.WithDescription("Messages broadcast to WebSocket clients"),
		)
		if err != nil {
			initErr = err
			return
		}
		apiRequests, err = meter.Int64Counter(
			"hub_api_requests",
			metric.WithDescription("HTTP API requests handled by the hub"),
		)
		if err != nil {
			initErr = err
		}
	})
	return initErr
}

func apiAttrs(route string) []metric.AddOption {
	return []metric.AddOption{
		metric.WithAttributes(attribute.String("route", route)),
	}
}
