package hub

import (
	"context"
	"net/http"
	"time"
)

func (m *metricsAPI) handlePodMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx, span := tracer.Start(r.Context(), "hub.metrics.pods")
	defer span.End()
	apiRequests.Add(ctx, 1, apiAttrs("/api/metrics/pods")...)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	items, err := m.collectPodMetrics(ctx)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, err)
		return
	}

	writeJSON(w, items)
}

func (m *metricsAPI) handleDeploymentMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx, span := tracer.Start(r.Context(), "hub.metrics.deployments")
	defer span.End()
	apiRequests.Add(ctx, 1, apiAttrs("/api/metrics/deployments")...)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	items, err := m.collectDeploymentMetrics(ctx)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, err)
		return
	}

	writeJSON(w, items)
}
