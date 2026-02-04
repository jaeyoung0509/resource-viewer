package tracesvc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jaeyoung050/resource-checker/internal/config"
	"github.com/jaeyoung050/resource-checker/pkg/health"
)

type traceResponse struct {
	Message string `json:"message"`
	TraceID string `json:"traceId"`
}

func Run(ctx context.Context, cfg config.TraceSvcConfig) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Handler)
	mux.HandleFunc("/api/trace", handleTrace)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		return err
	}
}

func handleTrace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx, span := tracer.Start(r.Context(), "trace.handle")
	defer span.End()

	traceRequests.Add(ctx, 1)

	traceID := span.SpanContext().TraceID().String()
	writeJSON(w, traceResponse{
		Message: "trace generated",
		TraceID: traceID,
	})
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(payload)
}
