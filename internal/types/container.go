package types

import (
	"github.com/username/project-name/ent"
	"github.com/username/project-name/internal/platform/analytics"
	"github.com/username/project-name/internal/platform/cache"
	"go.uber.org/zap"
)

type Container struct {
	DB        *ent.Client
	Logger    *zap.Logger
	Analytics analytics.Analtics
	Cache     cache.Cache
}
