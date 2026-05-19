package httpapi

import (
	"net/http"

	"wakeup-tracker/backend/services/device-validation-service/internal/application"
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
	mux.HandleFunc("POST /devices/validate", h.validate)
}

func (h *Handler) validate(w http.ResponseWriter, r *http.Request) {
	var request contracts.DeviceValidationRequest
	if err := infra.DecodeJSON(r, &request); err != nil {
		infra.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	response, proofID, err := h.service.Validate(r.Context(), request)
	if err != nil {
		infra.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	status := http.StatusOK
	if !response.Allowed {
		status = http.StatusForbidden
	}
	infra.WriteJSON(w, status, map[string]any{"validation": response, "device_proof_id": proofID})
}
