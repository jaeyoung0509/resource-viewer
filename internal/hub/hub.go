package hub

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"

	"github.com/jaeyoung050/resource-checker/internal/config"
	"github.com/jaeyoung050/resource-checker/pkg/health"
)

type client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

type Hub struct {
	mu       sync.RWMutex
	clients  map[*client]struct{}
	upgrader websocket.Upgrader
}

func New() *Hub {
	return &Hub{
		clients: make(map[*client]struct{}),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool { return true },
		},
	}
}

func Run(ctx context.Context, cfg config.HubConfig) error {
	if _, err := os.Stat(cfg.TLSCertPath); err != nil {
		return err
	}
	if _, err := os.Stat(cfg.TLSKeyPath); err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	pubsub := redisClient.PSubscribe(ctx, cfg.Pattern)
	defer pubsub.Close()

	hub := New()
	scaleAPI, err := newScaler(cfg)
	if err != nil {
		return err
	}
	metricsAPI, err := newMetricsAPI(cfg)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Handler)
	if scaleAPI != nil {
		mux.HandleFunc("/api/targets", scaleAPI.handleTargets)
		mux.HandleFunc("/api/scale", scaleAPI.handleScale)
	}
	if metricsAPI != nil {
		mux.HandleFunc("/api/metrics/pods", metricsAPI.handlePodMetrics)
		mux.HandleFunc("/api/metrics/deployments", metricsAPI.handleDeploymentMetrics)
	}
	mux.Handle(cfg.WSPath, http.HandlerFunc(hub.handleWS))
	mux.Handle("/", http.FileServer(http.Dir(cfg.StaticDir)))

	srv := &http.Server{
		Addr:              cfg.HTTPSAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	go func() {
		for msg := range pubsub.Channel() {
			hub.broadcast([]byte(msg.Payload))
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		hub.closeAll()
		return nil
	case err := <-errCh:
		return err
	}
}

func (h *Hub) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientConn := &client{conn: conn}
	h.add(clientConn)
	defer h.remove(clientConn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (h *Hub) broadcast(payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for clientConn := range h.clients {
		clientConn.mu.Lock()
		_ = clientConn.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_ = clientConn.conn.WriteMessage(websocket.TextMessage, payload)
		clientConn.mu.Unlock()
	}
}

func (h *Hub) add(clientConn *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	m := h.clients
	m[clientConn] = struct{}{}
}

func (h *Hub) remove(clientConn *client) {
	h.mu.Lock()
	delete(h.clients, clientConn)
	h.mu.Unlock()
	_ = clientConn.conn.Close()
}

func (h *Hub) closeAll() {
	h.mu.Lock()
	clients := h.clients
	h.clients = make(map[*client]struct{})
	h.mu.Unlock()

	for clientConn := range clients {
		_ = clientConn.conn.Close()
	}
}
