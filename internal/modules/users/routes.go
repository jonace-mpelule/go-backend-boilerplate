package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/username/project-name/internal/middlewares"
	"github.com/username/project-name/internal/permissions"
)

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Use(middlewares.AuthGuard(m.jwt))
		r.Use(middlewares.RequirePermission(middlewares.RequireAll, permissions.UserRead))
		r.Get("/", m.handler.ListUsers)
	})
}
