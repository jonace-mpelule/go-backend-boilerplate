package health

import "github.com/go-chi/chi/v5"

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/health", func(r chi.Router) {
		r.Get("/", m.handler.Ready)
		r.Get("/live", m.handler.Live)
		r.Get("/ready", m.handler.Ready)
	})
}
