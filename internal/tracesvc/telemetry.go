package tracesvc

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer         = otel.Tracer("trace-svc")
	meter          = otel.Meter("trace-svc")
	traceRequests  metric.Int64Counter
)

func init() {
	traceRequests, _ = meter.Int64Counter(
		"trace_requests_total",
		metric.WithDescription("Total trace requests handled by trace-svc"),
	)
}
