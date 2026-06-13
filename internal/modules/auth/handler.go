package auth

import (
	"net/http"

	"github.com/username/project-name/internal/middlewares"
	"github.com/username/project-name/internal/request"
	"github.com/username/project-name/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	input, ok := request.ValidatedPayload[RegisterRequest](r.Context())
	if !ok {
		response.Error(w, r, nil)
		return
	}

	result, appErr := h.service.Register(r.Context(), *input)
	if appErr != nil {
		response.Error(w, r, appErr)
		return
	}

	response.Created(w, r, result)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	input, ok := request.ValidatedPayload[LoginRequest](r.Context())
	if !ok {
		response.Error(w, r, nil)
		return
	}

	result, appErr := h.service.Login(r.Context(), *input)
	if appErr != nil {
		response.Error(w, r, appErr)
		return
	}

	response.Success(w, r, http.StatusOK, result)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	input, ok := request.ValidatedPayload[RefreshRequest](r.Context())
	if !ok {
		response.Error(w, r, nil)
		return
	}

	result, appErr := h.service.Refresh(r.Context(), *input)
	if appErr != nil {
		response.Error(w, r, appErr)
		return
	}

	response.Success(w, r, http.StatusOK, result)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.ClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, r, nil)
		return
	}

	profile, appErr := h.service.Me(r.Context(), claims.UserID)
	if appErr != nil {
		response.Error(w, r, appErr)
		return
	}

	response.Success(w, r, http.StatusOK, profile)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	input, ok := request.ValidatedPayload[ForgotPasswordRequest](r.Context())
	if !ok {
		response.Error(w, r, nil)
		return
	}

	if appErr := h.service.ForgotPassword(r.Context(), *input); appErr != nil {
		response.Error(w, r, appErr)
		return
	}

	response.Success(w, r, http.StatusOK, MessageResponse{Message: "If the account exists, a reset email has been sent"})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	input, ok := request.ValidatedPayload[ResetPasswordRequest](r.Context())
	if !ok {
		response.Error(w, r, nil)
		return
	}

	if appErr := h.service.ResetPassword(r.Context(), *input); appErr != nil {
		response.Error(w, r, appErr)
		return
	}

	response.Success(w, r, http.StatusOK, MessageResponse{Message: "Password reset successful"})
}
