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

// ListUsers godoc
// @Summary List users
// @Description Returns the current user listing DTOs for the starter users module.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope{data=[]UserResponse}
// @Failure 401 {object} response.Envelope{error=response.ErrorBody}
// @Failure 403 {object} response.Envelope{error=response.ErrorBody}
// @Failure 500 {object} response.Envelope{error=response.ErrorBody}
// @Router /users/ [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		response.Error(w, r, apperrors.Internal("failed to list users"))
		return
	}

	response.Success(w, r, http.StatusOK, users)
}
