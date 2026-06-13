package response

import (
	"net/http"

	"github.com/go-chi/render"
	apperrors "github.com/username/project-name/internal/errors"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Envelope struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, r *http.Request, status int, data any) {
	render.Status(r, status)
	render.JSON(w, r, Envelope{
		Success: true,
		Data:    data,
	})
}

func Created(w http.ResponseWriter, r *http.Request, data any) {
	Success(w, r, http.StatusCreated, data)
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	render.NoContent(w, r)
}

func Error(w http.ResponseWriter, r *http.Request, err *apperrors.AppError) {
	if err == nil {
		err = apperrors.Internal("internal server error")
	}

	render.Status(r, err.Status)
	render.JSON(w, r, Envelope{
		Success: false,
		Error: &ErrorBody{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		},
	})
}
