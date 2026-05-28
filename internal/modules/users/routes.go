package users

import "github.com/go-chi/chi/v5"

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", m.handler.GetUsers)
	})
}
