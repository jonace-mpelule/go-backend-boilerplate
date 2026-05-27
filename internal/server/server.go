package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/username/project-name/internal/config"
)

type Server struct {
	cfg    *config.Config
	router *chi.Mux
}

func New(cfg *config.Config, router *chi.Mux) *Server {
	return &Server{
		cfg:    cfg,
		router: router,
	}
}

func (s *Server) Start() error {
	s.registerRoutes()

	addr := ":" + s.cfg.Port

	log.Default().Print("Server Running On ", addr)

	return http.ListenAndServe(addr, s.router)
}

func (s *Server) registerRoutes() {

	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

}
