package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/username/project-name/internal/config"
	"github.com/username/project-name/internal/middlewares"
	"github.com/username/project-name/internal/modules"
	"github.com/username/project-name/internal/modules/auth"
	"github.com/username/project-name/internal/modules/health"
	"github.com/username/project-name/internal/modules/users"
	_ "github.com/username/project-name/internal/platform/docs"
	"github.com/username/project-name/internal/types"
	"go.uber.org/zap"
)

type Server struct {
	cfg       *config.Config
	router    *chi.Mux
	container *types.Container
	modules   []modules.Module
}

func NewServer(cfg *config.Config, router *chi.Mux, container *types.Container) *Server {
	s := &Server{
		cfg:       cfg,
		router:    router,
		container: container,
	}

	s.registerMiddleware()
	s.registerUtilityRoutes()
	s.registerModules()

	return s
}

func (s *Server) Start() error {
	addr := ":" + s.cfg.Server.Port

	server := &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadHeaderTimeout: s.cfg.Server.ReadHeaderTimeout,
		ReadTimeout:       s.cfg.Server.ReadTimeout,
		WriteTimeout:      s.cfg.Server.WriteTimeout,
		IdleTimeout:       s.cfg.Server.IdleTimeout,
	}

	serverErr := make(chan error, 1)

	go func() {
		s.container.Logger.Info("server starting", zap.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case <-stop:
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.ShutdownTimeout)
	defer cancel()

	shutdownErr := server.Shutdown(ctx)
	closeErr := s.container.Close(ctx)

	return errors.Join(shutdownErr, closeErr)
}

func (s *Server) registerMiddleware() {
	s.router.Use(chimiddleware.RequestID)
	if s.container.Metrics != nil && s.container.Metrics.Enabled() {
		s.router.Use(s.container.Metrics.Middleware)
	}
	s.router.Use(middlewares.Recoverer(s.container.Logger))
	s.router.Use(middlewares.RequestLogger(s.container.Logger))
	s.router.Use(middlewares.Security)
	s.router.Use(middlewares.CORS(s.cfg.CORS))
	s.router.Use(middlewares.BodyLimiter(s.cfg.Server.BodyLimitBytes))
	if s.cfg.RateLimit.Enabled {
		s.router.Use(middlewares.NewRateLimiter(s.cfg.RateLimit).Handler)
	}
}

func (s *Server) registerUtilityRoutes() {
	if s.container.Metrics != nil && s.container.Metrics.Enabled() {
		s.router.Handle(s.container.Metrics.Path(), s.container.Metrics.Handler())
	}
	s.router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusTemporaryRedirect)
	})
	s.router.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DeepLinking(true),
	))
}

func (s *Server) registerModules() {
	s.modules = []modules.Module{
		auth.NewModule(s.container, s.cfg.Auth),
		health.NewModule(s.container),
		users.NewModule(s.container, s.cfg.Cache.DefaultTTL),
	}

	for _, module := range s.modules {
		module.RegisterRoutes(s.router)
	}
}
