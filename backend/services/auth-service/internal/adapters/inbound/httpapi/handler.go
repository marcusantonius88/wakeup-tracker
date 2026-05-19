package httpapi

import (
	"net/http"

	"wakeup-tracker/backend/services/auth-service/internal/application"
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
	mux.HandleFunc("POST /login", h.login)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var request contracts.LoginRequest
	if err := infra.DecodeJSON(r, &request); err != nil {
		infra.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	response, err := h.service.Login(r.Context(), request)
	if err != nil {
		infra.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	infra.WriteJSON(w, http.StatusOK, response)
}
