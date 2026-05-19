package infrastructure

import (
	"net/http"
	"time"

	"wakeup-tracker/backend/services/projection-service/internal/adapters/inbound/httpapi"
	"wakeup-tracker/backend/services/projection-service/internal/adapters/outbound/memory"
	"wakeup-tracker/backend/services/projection-service/internal/application"
	"wakeup-tracker/backend/shared/idempotency"
	"wakeup-tracker/backend/shared/infra"
	"wakeup-tracker/backend/shared/observability"
)

func NewServer() *http.Server {
	metrics := observability.NewMetrics()
	repository := memory.NewRepository()
	idempotencyStore := idempotency.NewMemoryStore()
	service := application.NewService(repository, idempotencyStore, metrics)
	handler := httpapi.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", infra.HealthHandler("projection-service"))
	mux.HandleFunc("GET /metrics", metrics.Handler())
	handler.Register(mux)

	return &http.Server{Addr: ":8086", Handler: withCORS(mux), ReadHeaderTimeout: 5 * time.Second}
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
