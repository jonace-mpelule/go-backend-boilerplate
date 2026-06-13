package health

import (
	"net/http"

	"github.com/username/project-name/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	response.Success(w, r, http.StatusOK, h.service.Live())
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	status, payload := h.service.Ready(r.Context())
	response.Success(w, r, status, payload)
}
