package auth

import (
	"github.com/username/project-name/internal/config"
	"github.com/username/project-name/internal/types"
)

type Module struct {
	handler *Handler
	jwt     JWTVerifier
}

func NewModule(container *types.Container, cfg config.AuthConfig) *Module {
	repo := NewRepository(container.DB.Ent)
	service := NewService(
		repo,
		container.JWT,
		container.Passwords,
		container.Mailer,
		container.Analytics,
		cfg.TokenTTL,
		cfg.RefreshTokenTTL,
		cfg.ResetTokenTTL,
	)

	return &Module{
		handler: NewHandler(service),
		jwt:     container.JWT,
	}
}
