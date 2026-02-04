package hub

import (
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
)

func init() {
	wsConnections, _ = meter.Int64Counter(
		"hub_ws_connections",
		metric.WithDescription("WebSocket connections accepted by the hub"),
	)
	wsMessages, _ = meter.Int64Counter(
		"hub_ws_messages",
		metric.WithDescription("Messages broadcast to WebSocket clients"),
	)
	apiRequests, _ = meter.Int64Counter(
		"hub_api_requests",
		metric.WithDescription("HTTP API requests handled by the hub"),
	)
}

func apiAttrs(route string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("route", route),
	}
}
