package users

import (
	"net/http"

	apperrors "github.com/username/project-name/internal/errors"
	"github.com/username/project-name/internal/response"
)

type Handler struct {
	service ServiceContract
}

func NewHandler(service ServiceContract) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		response.Error(w, r, apperrors.Internal("failed to list users"))
		return
	}

	response.Success(w, r, http.StatusOK, users)
}
