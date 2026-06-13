package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/username/project-name/internal/middlewares"
)

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.With(middlewares.ValidationGuard[RegisterRequest]()).Post("/register", m.handler.Register)
		r.With(middlewares.ValidationGuard[LoginRequest]()).Post("/login", m.handler.Login)
		r.With(middlewares.ValidationGuard[RefreshRequest]()).Post("/refresh", m.handler.Refresh)
		r.With(middlewares.ValidationGuard[ForgotPasswordRequest]()).Post("/forgot-password", m.handler.ForgotPassword)
		r.With(middlewares.ValidationGuard[ResetPasswordRequest]()).Post("/reset-password", m.handler.ResetPassword)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthGuard(m.jwt))
			r.Get("/me", m.handler.Me)
		})
	})
}
