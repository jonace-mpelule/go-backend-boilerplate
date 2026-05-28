package health

import "github.com/go-chi/chi/v5"

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Get("/health", m.handler.Health)
}
