package health

import (
	"context"
	"net/http"

	"github.com/username/project-name/internal/platform/cache"
)

type DBHealthChecker interface {
	PingContext(ctx context.Context) error
}

type Service struct {
	db    DBHealthChecker
	cache cache.Cache
}

func NewService(db DBHealthChecker, cacheClient cache.Cache) *Service {
	return &Service{
		db:    db,
		cache: cacheClient,
	}
}

func (s *Service) Live() StatusResponse {
	return StatusResponse{
		Status: "ok",
		Checks: map[string]string{
			"app": "ok",
		},
	}
}

func (s *Service) Ready(ctx context.Context) (int, StatusResponse) {
	checks := map[string]string{
		"database": "ok",
	}
	statusCode := http.StatusOK

	if err := s.db.PingContext(ctx); err != nil {
		checks["database"] = "error"
		statusCode = http.StatusServiceUnavailable
	}

	if err := s.cache.Ping(ctx); err == nil {
		checks["cache"] = "ok"
	} else if err == cache.ErrCacheDisabled {
		checks["cache"] = "disabled"
	} else {
		checks["cache"] = "error"
		statusCode = http.StatusServiceUnavailable
	}

	status := "ok"
	if statusCode != http.StatusOK {
		status = "degraded"
	}

	return statusCode, StatusResponse{
		Status: status,
		Checks: checks,
	}
}
