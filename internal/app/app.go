package app

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/username/project-name/internal/config"
	analyticsplatform "github.com/username/project-name/internal/platform/analytics"
	"github.com/username/project-name/internal/platform/cache"
	"github.com/username/project-name/internal/platform/db"
	"github.com/username/project-name/internal/platform/logger"
	"github.com/username/project-name/internal/platform/mailer"
	"github.com/username/project-name/internal/platform/metrics"
	sentryplatform "github.com/username/project-name/internal/platform/sentry"
	"github.com/username/project-name/internal/platform/storage"
	"github.com/username/project-name/internal/types"
	"github.com/username/project-name/internal/utils"
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

	log, logCloser, err := logger.New(cfg.App, cfg.Observability)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	dbClient, err := db.New(ctx, cfg.Database)
	if err != nil {
		if logCloser != nil {
			_ = logCloser.Close()
		}
		_ = log.Sync()
		return nil, err
	}

	cacheClient, err := cache.NewRedis(ctx, cfg.Cache.RedisURL)
	if err != nil {
		_ = dbClient.Close()
		if logCloser != nil {
			_ = logCloser.Close()
		}
		_ = log.Sync()
		return nil, err
	}

	sentryClient, err := sentryplatform.Init(cfg.Observability.SentryDSN)
	if err != nil {
		_ = cacheClient.Close()
		_ = dbClient.Close()
		if logCloser != nil {
			_ = logCloser.Close()
		}
		_ = log.Sync()
		return nil, err
	}

	analyticsClient := analyticsplatform.NewNoop()
	if cfg.Observability.PosthogKey != "" {
		analyticsClient, err = analyticsplatform.NewPosthog(
			cfg.Observability.PosthogKey,
			cfg.Observability.PosthogHost,
		)
		if err != nil {
			_ = cacheClient.Close()
			_ = dbClient.Close()
			if logCloser != nil {
				_ = logCloser.Close()
			}
			_ = log.Sync()
			return nil, err
		}
	}

	var mailerClient mailer.Mailer
	switch cfg.Mailer.Provider {
	case "noop":
		mailerClient = mailer.NewNoop()
	case "resend":
		mailerClient = mailer.NewResend(
			cfg.Mailer.ResendAPIKey,
			cfg.Mailer.ResendURL,
			cfg.Mailer.FromAddress,
		)
	case "smtp":
		mailerClient = mailer.NewSMTP(
			cfg.Mailer.SMTPHost,
			cfg.Mailer.SMTPPort,
			cfg.Mailer.SMTPUsername,
			cfg.Mailer.SMTPPassword,
			cfg.Mailer.FromAddress,
		)
	default:
		mailerClient = mailer.NewLog(log, cfg.Mailer.FromAddress)
	}

	var storageClient storage.Storage
	switch cfg.Storage.Provider {
	case "local":
		storageClient = storage.NewLocal(cfg.Storage.LocalDir)
	default:
		storageClient = storage.NewNoop()
	}

	container := &types.Container{
		DB:        dbClient,
		Logger:    log,
		Analytics: analyticsClient,
		Cache:     cacheClient,
		Mailer:    mailerClient,
		Storage:   storageClient,
		JWT: utils.NewJWT(
			cfg.Auth.JWTSecret,
			cfg.Auth.Issuer,
			cfg.Auth.TokenTTL,
		),
		Passwords: utils.NewPasswordHasher(),
		Sentry:    sentryClient,
		Metrics:   metrics.New(cfg.Observability.Metrics.Enabled, cfg.Observability.Metrics.Path, cfg.Observability.Metrics.Namespace),
		LogCloser: logCloser,
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
