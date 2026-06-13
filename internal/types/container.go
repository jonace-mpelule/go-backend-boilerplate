package types

import (
	"context"
	"io"

	"github.com/username/project-name/internal/platform/analytics"
	"github.com/username/project-name/internal/platform/cache"
	"github.com/username/project-name/internal/platform/db"
	"github.com/username/project-name/internal/platform/mailer"
	"github.com/username/project-name/internal/platform/metrics"
	"github.com/username/project-name/internal/platform/storage"
	sentryplatform "github.com/username/project-name/internal/platform/sentry"
	"github.com/username/project-name/internal/utils"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Container struct {
	DB        *db.Client
	Logger    *zap.Logger
	Analytics analytics.Analytics
	Cache     cache.Cache
	Mailer    mailer.Mailer
	Storage   storage.Storage
	JWT       *utils.JWT
	Passwords *utils.PasswordHasher
	Sentry    *sentryplatform.Client
	Metrics   *metrics.Provider
	LogCloser io.Closer
}

func (c *Container) Close(ctx context.Context) error {
	if c == nil {
		return nil
	}

	var closeErr error

	if c.Analytics != nil {
		closeErr = multierr.Append(closeErr, c.Analytics.Close())
	}

	if c.Cache != nil {
		closeErr = multierr.Append(closeErr, c.Cache.Close())
	}

	if c.LogCloser != nil {
		closeErr = multierr.Append(closeErr, c.LogCloser.Close())
	}

	if c.DB != nil {
		closeErr = multierr.Append(closeErr, c.DB.Close())
	}

	if c.Logger != nil {
		closeErr = multierr.Append(closeErr, c.Logger.Sync())
	}

	if c.Sentry != nil {
		c.Sentry.Flush(ctx)
	}

	return closeErr
}
