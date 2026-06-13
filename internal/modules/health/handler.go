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

// Live godoc
// @Summary Liveness probe
// @Description Returns a simple process-alive response.
// @Tags health
// @Produce json
// @Success 200 {object} response.Envelope{data=StatusResponse}
// @Router /health/live [get]
func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	response.Success(w, r, http.StatusOK, h.service.Live())
}

// Ready godoc
// @Summary Readiness probe
// @Description Returns dependency readiness for the API and its configured infrastructure.
// @Tags health
// @Produce json
// @Success 200 {object} response.Envelope{data=StatusResponse}
// @Success 503 {object} response.Envelope{data=StatusResponse}
// @Router /health/ready [get]
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	status, payload := h.service.Ready(r.Context())
	response.Success(w, r, status, payload)
}
