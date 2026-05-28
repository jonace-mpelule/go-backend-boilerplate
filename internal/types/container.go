package types

import (
	"github.com/username/project-name/ent"
	"github.com/username/project-name/internal/platform/analytics"
	"go.uber.org/zap"
)

type Container struct {
	DB        *ent.Client
	Logger    *zap.Logger
	Analytics analytics.Analtics
}
