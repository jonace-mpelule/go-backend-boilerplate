package health_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/username/project-name/internal/modules/health"
	"github.com/username/project-name/internal/platform/cache"
)

type stubDB struct {
	err error
}

func (s stubDB) PingContext(context.Context) error {
	return s.err
}

type stubCache struct {
	err error
}

func (s stubCache) Set(context.Context, string, []byte, time.Duration) error { return nil }
func (s stubCache) Get(context.Context, string) ([]byte, error)              { return nil, cache.ErrCacheMiss }
func (s stubCache) Delete(context.Context, string) error                     { return nil }
func (s stubCache) Ping(context.Context) error                               { return s.err }
func (s stubCache) Close() error                                             { return nil }

func TestReadyReturnsUnavailableWhenDatabaseFails(t *testing.T) {
	service := health.NewService(stubDB{err: errors.New("db down")}, cache.NewNoop())

	status, payload := service.Ready(context.Background())
	if status != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", status)
	}

	if payload.Checks["database"] != "error" {
		t.Fatal("expected database health to be error")
	}
}

func TestReadyReturnsOKWhenDependenciesHealthy(t *testing.T) {
	service := health.NewService(stubDB{}, stubCache{})

	status, payload := service.Ready(context.Background())
	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}

	if payload.Status != "ok" {
		t.Fatalf("expected ok payload status, got %s", payload.Status)
	}
}
