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

// Register godoc
// @Summary Register a user
// @Description Creates a new user account and returns access and refresh tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body RegisterRequest true "Registration payload"
// @Success 201 {object} response.Envelope{data=AuthResponse}
// @Failure 400 {object} response.Envelope{error=response.ErrorBody}
// @Failure 409 {object} response.Envelope{error=response.ErrorBody}
// @Failure 500 {object} response.Envelope{error=response.ErrorBody}
// @Router /auth/register [post]
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

// Login godoc
// @Summary Login
// @Description Authenticates a user and returns access and refresh tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body LoginRequest true "Login payload"
// @Success 200 {object} response.Envelope{data=AuthResponse}
// @Failure 400 {object} response.Envelope{error=response.ErrorBody}
// @Failure 401 {object} response.Envelope{error=response.ErrorBody}
// @Failure 500 {object} response.Envelope{error=response.ErrorBody}
// @Router /auth/login [post]
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

// Refresh godoc
// @Summary Refresh access token
// @Description Exchanges a refresh token for a new token pair.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body RefreshRequest true "Refresh payload"
// @Success 200 {object} response.Envelope{data=AuthResponse}
// @Failure 400 {object} response.Envelope{error=response.ErrorBody}
// @Failure 401 {object} response.Envelope{error=response.ErrorBody}
// @Failure 500 {object} response.Envelope{error=response.ErrorBody}
// @Router /auth/refresh [post]
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

// Me godoc
// @Summary Current user
// @Description Returns the authenticated user's profile.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope{data=ProfileResponse}
// @Failure 401 {object} response.Envelope{error=response.ErrorBody}
// @Failure 500 {object} response.Envelope{error=response.ErrorBody}
// @Router /auth/me [get]
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

// ForgotPassword godoc
// @Summary Request password reset
// @Description Starts the password reset flow and sends an email through the configured mailer.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body ForgotPasswordRequest true "Password reset request payload"
// @Success 200 {object} response.Envelope{data=MessageResponse}
// @Failure 400 {object} response.Envelope{error=response.ErrorBody}
// @Failure 500 {object} response.Envelope{error=response.ErrorBody}
// @Router /auth/forgot-password [post]
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

// ResetPassword godoc
// @Summary Confirm password reset
// @Description Completes a password reset using the emailed reset token.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body ResetPasswordRequest true "Password reset confirmation payload"
// @Success 200 {object} response.Envelope{data=MessageResponse}
// @Failure 400 {object} response.Envelope{error=response.ErrorBody}
// @Failure 401 {object} response.Envelope{error=response.ErrorBody}
// @Failure 500 {object} response.Envelope{error=response.ErrorBody}
// @Router /auth/reset-password [post]
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
