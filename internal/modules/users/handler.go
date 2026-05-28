package users

import (
	"net/http"

	"github.com/go-chi/render"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetUsers(
	w http.ResponseWriter,
	r *http.Request,
) {
	users, err := h.service.GetAllUsers(r.Context())

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"message": "ok",
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]any{
		"data": users,
	})
}
