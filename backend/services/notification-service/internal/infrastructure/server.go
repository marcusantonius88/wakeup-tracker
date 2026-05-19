package infrastructure

import (
	"net/http"
	"time"

	"wakeup-tracker/backend/services/notification-service/internal/adapters/inbound/httpapi"
	"wakeup-tracker/backend/services/notification-service/internal/adapters/outbound/memory"
	"wakeup-tracker/backend/services/notification-service/internal/application"
	"wakeup-tracker/backend/shared/infra"
	"wakeup-tracker/backend/shared/observability"
)

func NewServer() *http.Server {
	logger := observability.NewLogger("notification-service")
	metrics := observability.NewMetrics()
	repository := memory.NewNotificationRepository()
	eventBus := infra.NewInMemoryEventBus(logger)
	publisher := infra.NewRetryingPublisher(eventBus, 3, 100*time.Millisecond, metrics, logger)
	service := application.NewService(repository, publisher, metrics)
	handler := httpapi.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", infra.HealthHandler("notification-service"))
	mux.HandleFunc("GET /metrics", metrics.Handler())
	handler.Register(mux)

	return &http.Server{Addr: ":8084", Handler: withCORS(mux), ReadHeaderTimeout: 5 * time.Second}
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
