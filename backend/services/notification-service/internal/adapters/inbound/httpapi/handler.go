package httpapi

import (
	"net/http"

	"wakeup-tracker/backend/services/notification-service/internal/application"
	"wakeup-tracker/backend/shared/contracts"
	"wakeup-tracker/backend/shared/infra"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /notifications/send", h.send)
	mux.HandleFunc("GET /notifications", h.list)
}

func (h *Handler) send(w http.ResponseWriter, r *http.Request) {
	var request contracts.NotificationRequest
	if err := infra.DecodeJSON(r, &request); err != nil {
		infra.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	notification, err := h.service.Send(r.Context(), request)
	if err != nil {
		infra.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	infra.WriteJSON(w, http.StatusCreated, notification)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "demo-user"
	}
	notifications, err := h.service.List(r.Context(), userID)
	if err != nil {
		infra.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	infra.WriteJSON(w, http.StatusOK, map[string]any{"notifications": notifications})
}
