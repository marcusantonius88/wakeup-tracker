package httpapi

import (
	"net/http"

	"wakeup-tracker/backend/services/wake-session-service/internal/application"
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
	mux.HandleFunc("POST /wake-sessions/check-in", h.checkIn)
	mux.HandleFunc("GET /wake-sessions/history", h.history)
}

func (h *Handler) checkIn(w http.ResponseWriter, r *http.Request) {
	var request contracts.WakeCheckInRequest
	if err := infra.DecodeJSON(r, &request); err != nil {
		infra.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	response, err := h.service.CheckIn(r.Context(), request)
	if err != nil {
		infra.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	infra.WriteJSON(w, http.StatusCreated, response)
}

func (h *Handler) history(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "demo-user"
	}
	sessions, err := h.service.History(r.Context(), userID)
	if err != nil {
		infra.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	infra.WriteJSON(w, http.StatusOK, map[string]any{"sessions": sessions})
}
