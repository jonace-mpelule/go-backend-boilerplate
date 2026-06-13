package users

import (
	"time"

	"github.com/username/project-name/internal/types"
)

type Module struct {
	handler *Handler
	jwt     JWTVerifier
}

func NewModule(container *types.Container, cacheTTL time.Duration) *Module {
	repo := NewRepository(container.DB.Ent)

	service := NewService(repo, container.Cache, container.Analytics, cacheTTL)

	handler := NewHandler(service)

	return &Module{
		handler: handler,
		jwt:     container.JWT,
	}
}
