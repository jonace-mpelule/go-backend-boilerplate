package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/username/project-name/internal/config"
	"github.com/username/project-name/internal/server"
)

type App struct {
	Config *config.Config
	Router *chi.Mux
	Server *server.Server
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	srv := server.New(cfg, router)

	return &App{
		Config: cfg,
		Router: router,
		Server: srv,
	}, nil
}

func (a *App) Run() error {
	return a.Server.Start()
}
