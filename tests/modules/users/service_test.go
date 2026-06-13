package users_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/username/project-name/internal/modules/users"
	"github.com/username/project-name/internal/platform/cache"
)

type stubRepository struct {
	listFn func(context.Context) ([]users.UserRecord, error)
}

func (s stubRepository) List(ctx context.Context) ([]users.UserRecord, error) {
	return s.listFn(ctx)
}

type stubCache struct {
	getFn func(context.Context, string) ([]byte, error)
	setFn func(context.Context, string, []byte, time.Duration) error
}

func (s stubCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if s.setFn != nil {
		return s.setFn(ctx, key, value, ttl)
	}
	return nil
}

func (s stubCache) Get(ctx context.Context, key string) ([]byte, error) {
	if s.getFn != nil {
		return s.getFn(ctx, key)
	}
	return nil, cache.ErrCacheMiss
}

func (s stubCache) Delete(context.Context, string) error { return nil }
func (s stubCache) Ping(context.Context) error           { return nil }
func (s stubCache) Close() error                         { return nil }

type stubAnalytics struct {
	tracked bool
}

func (s *stubAnalytics) Track(context.Context, string, map[string]any) {
	s.tracked = true
}

func (s *stubAnalytics) Close() error {
	return nil
}

func TestListUsersReturnsCachedUsers(t *testing.T) {
	cached, _ := json.Marshal([]users.UserResponse{{ID: "user-1", Email: "cached@example.com"}})

	service := users.NewService(
		stubRepository{listFn: func(context.Context) ([]users.UserRecord, error) {
			t.Fatal("repository should not be called on cache hit")
			return nil, nil
		}},
		stubCache{getFn: func(context.Context, string) ([]byte, error) {
			return cached, nil
		}},
		&stubAnalytics{},
		time.Minute,
	)

	result, err := service.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("expected list users to succeed: %v", err)
	}

	if len(result) != 1 || result[0].Email != "cached@example.com" {
		t.Fatal("expected cached users to be returned")
	}
}

func TestListUsersCachesRepositoryResults(t *testing.T) {
	analytics := &stubAnalytics{}
	cached := false

	service := users.NewService(
		stubRepository{listFn: func(context.Context) ([]users.UserRecord, error) {
			return []users.UserRecord{{ID: "user-1", Email: "user@example.com"}}, nil
		}},
		stubCache{
			getFn: func(context.Context, string) ([]byte, error) {
				return nil, cache.ErrCacheMiss
			},
			setFn: func(_ context.Context, key string, value []byte, ttl time.Duration) error {
				cached = key == "users:list" && ttl == time.Minute
				return nil
			},
		},
		analytics,
		time.Minute,
	)

	result, err := service.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("expected list users to succeed: %v", err)
	}

	if len(result) != 1 || result[0].ID != "user-1" {
		t.Fatal("expected repository results to be returned")
	}

	if !cached {
		t.Fatal("expected users list to be cached")
	}

	if !analytics.tracked {
		t.Fatal("expected analytics event to be tracked")
	}
}
