package infrastructure

import (
	"context"
	"net/http"
	"os"
	"time"

	"wakeup-tracker/backend/services/wake-session-service/internal/adapters/inbound/httpapi"
	"wakeup-tracker/backend/services/wake-session-service/internal/adapters/outbound/memory"
	"wakeup-tracker/backend/services/wake-session-service/internal/application"
	"wakeup-tracker/backend/shared/idempotency"
	"wakeup-tracker/backend/shared/infra"
	"wakeup-tracker/backend/shared/observability"
)

func NewServer() *http.Server {
	logger := observability.NewLogger("wake-session-service")
	metrics := observability.NewMetrics()
	repository := memory.NewWakeSessionRepository()
	eventBus := infra.NewInMemoryEventBus(logger)
	publisher := infra.NewRetryingPublisher(eventBus, 3, 100*time.Millisecond, metrics, logger)
	var idempotencyStore interface {
		Reserve(context.Context, string, time.Duration) (bool, error)
	} = idempotency.NewMemoryStore()
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		idempotencyStore = idempotency.NewRedisStore(redisAddr)
	}
	service := application.NewService(repository, publisher, idempotencyStore, metrics, logger)
	handler := httpapi.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", infra.HealthHandler("wake-session-service"))
	mux.HandleFunc("GET /metrics", metrics.Handler())
	handler.Register(mux)

	return &http.Server{
		Addr:              ":8082",
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Shutdown(ctx context.Context, server *http.Server) error {
	return server.Shutdown(ctx)
}
