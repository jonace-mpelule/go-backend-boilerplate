package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/username/project-name/internal/config"
	analyticsplatform "github.com/username/project-name/internal/platform/analytics"
	"github.com/username/project-name/internal/platform/db"
	"github.com/username/project-name/internal/platform/logger"
	sentryplatform "github.com/username/project-name/internal/platform/sentry"
	"github.com/username/project-name/internal/types"
)

type App struct {
	Config    *config.Config
	Container *types.Container
	Server    *Server
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	dbClient, err := db.New(cfg.DBUrl)
	if err != nil {
		return nil, err
	}

	if err := sentryplatform.Init(
		cfg.SentryDSN,
	); err != nil {
		return nil, err

	}

	var analytics analyticsplatform.Analtics

	if cfg.PosthogKey != "" {
		posthogAnalytics, err := analyticsplatform.NewPosthog(
			cfg.PosthogKey,
			cfg.PosthogHost,
		)
		if err != nil {
			return nil, err
		}
		analytics = posthogAnalytics
	} else {
		analytics = analyticsplatform.NewNoop()
	}

	container := &types.Container{
		DB:        dbClient,
		Logger:    log,
		Analytics: analytics,
	}

	router := chi.NewRouter()
	srv := NewServer(cfg, router, container)

	return &App{
		Config:    cfg,
		Container: container,
		Server:    srv,
	}, nil
}

func (a *App) Run() error {
	return a.Server.Start()
}
