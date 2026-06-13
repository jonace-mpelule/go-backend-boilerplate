package users

import (
	"context"
	"encoding/json"
	"time"

	"github.com/username/project-name/internal/platform/analytics"
	"github.com/username/project-name/internal/platform/cache"
)

const usersListCacheKey = "users:list"

type ServiceContract interface {
	ListUsers(ctx context.Context) ([]UserResponse, error)
}

type RepositoryContract interface {
	List(ctx context.Context) ([]UserRecord, error)
}

type Service struct {
	repository RepositoryContract
	cache      cache.Cache
	analytics  analytics.Analytics
	cacheTTL   time.Duration
}

func NewService(
	repository RepositoryContract,
	cacheClient cache.Cache,
	analyticsClient analytics.Analytics,
	cacheTTL time.Duration,
) *Service {
	return &Service{
		repository: repository,
		cache:      cacheClient,
		analytics:  analyticsClient,
		cacheTTL:   cacheTTL,
	}
}

func (s *Service) ListUsers(ctx context.Context) ([]UserResponse, error) {
	if cached, err := s.cache.Get(ctx, usersListCacheKey); err == nil {
		var users []UserResponse
		if unmarshalErr := json.Unmarshal(cached, &users); unmarshalErr == nil {
			return users, nil
		}
	}

	records, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]UserResponse, 0, len(records))
	for _, record := range records {
		users = append(users, UserResponse{
			ID:    record.ID,
			Email: record.Email,
		})
	}

	if encoded, err := json.Marshal(users); err == nil {
		_ = s.cache.Set(ctx, usersListCacheKey, encoded, s.cacheTTL)
	}

	s.analytics.Track(ctx, "users.listed", map[string]any{
		"count": len(users),
	})

	return users, nil
}
