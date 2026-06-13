package health

import "github.com/username/project-name/internal/types"

type Module struct {
	handler *Handler
}

func NewModule(container *types.Container) *Module {
	service := NewService(container.DB, container.Cache)
	return &Module{handler: NewHandler(service)}
}
