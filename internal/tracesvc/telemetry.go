package tracesvc

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer        = otel.Tracer("trace-svc")
	meter         = otel.Meter("trace-svc")
	traceRequests metric.Int64Counter
	initOnce      sync.Once
	initErr       error
)

func initTelemetry() error {
	initOnce.Do(func() {
		var err error
		traceRequests, err = meter.Int64Counter(
			"trace_requests_total",
			metric.WithDescription("Total trace requests handled by trace-svc"),
		)
		if err != nil {
			initErr = err
		}
	})
	return initErr
}
