package users

import (
	"context"

	"github.com/username/project-name/ent"
	"github.com/username/project-name/internal/platform/analytics"
)

type Service struct {
	repository *Repository
	analytics  analytics.Analtics
}

func NewService(
	repository *Repository,
	analytics analytics.Analtics,
) *Service {
	return &Service{
		repository: repository,
		analytics:  analytics,
	}
}

func (s *Service) GetAllUsers(
	ctx context.Context,
) ([]*ent.User, error) {
	return s.repository.GetAll(ctx)
}
