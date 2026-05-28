package users

import (
	"github.com/username/project-name/internal/types"
)

type Module struct {
	handler *Handler
}

func NewModule(container *types.Container) *Module {
	repo := NewRepository(container.DB)

	service := NewService(repo, container.Analytics)

	handler := NewHandler(service)

	return &Module{
		handler: handler,
	}
}
