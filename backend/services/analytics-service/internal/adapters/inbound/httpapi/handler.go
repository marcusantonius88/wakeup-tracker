package httpapi

import (
	"net/http"

	"wakeup-tracker/backend/services/analytics-service/internal/application"
	"wakeup-tracker/backend/shared/events"
	"wakeup-tracker/backend/shared/infra"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /events", h.process)
	mux.HandleFunc("GET /analytics/consistency", h.get)
}

func (h *Handler) process(w http.ResponseWriter, r *http.Request) {
	var event events.Envelope
	if err := infra.DecodeJSON(r, &event); err != nil {
		infra.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	metrics, err := h.service.Process(r.Context(), event)
	if err != nil {
		infra.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	infra.WriteJSON(w, http.StatusAccepted, metrics)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "demo-user"
	}
	metrics, err := h.service.Get(r.Context(), userID)
	if err != nil {
		infra.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	infra.WriteJSON(w, http.StatusOK, metrics)
}
