package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/username/project-name/internal/config"
	"github.com/username/project-name/internal/middlewares"
	"github.com/username/project-name/internal/modules"
	"github.com/username/project-name/internal/modules/health"
	"github.com/username/project-name/internal/modules/users"
	"github.com/username/project-name/internal/types"
)

type Server struct {
	cfg       *config.Config
	router    *chi.Mux
	container *types.Container
	modules   []modules.Module
}

func NewServer(
	cfg *config.Config,
	router *chi.Mux,
	container *types.Container,
) *Server {

	s := &Server{
		cfg:    cfg,
		router: router,
	}

	s.registerMiddleware()
	s.registerModules()

	return &Server{
		cfg:       cfg,
		router:    router,
		container: container,
	}
}

func (s *Server) Start() error {
	addr := ":" + s.cfg.Port
	log.Default().Print("Server Running On ", addr)

	return http.ListenAndServe(addr, s.router)
}

func (s *Server) registerMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middlewares.Security)
	s.router.Use(middleware.Logger)
	s.router.Use(middlewares.Sentry)
}

func (s *Server) registerModules() {
	s.modules = []modules.Module{
		health.NewModule(),
		users.NewModule(s.container),
	}

	for _, module := range s.modules {
		module.RegisterRoutes(s.router)
	}
}
