package logger_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/username/project-name/internal/config"
	"github.com/username/project-name/internal/platform/logger"
)

func TestLoggerInitializesWithoutLoki(t *testing.T) {
	log, closer, err := logger.New(
		config.AppConfig{Name: "svc", Env: "development"},
		config.ObservabilityConfig{},
	)
	if err != nil {
		t.Fatalf("expected logger init to succeed: %v", err)
	}
	if closer != nil {
		t.Fatal("expected no Loki closer when Loki is disabled")
	}
	_ = log.Sync()
}

func TestLoggerPushesToLokiWhenEnabled(t *testing.T) {
	var (
		mu       sync.Mutex
		payloads [][]byte
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/loki/api/v1/push" {
			t.Fatalf("unexpected push path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()

		mu.Lock()
		payloads = append(payloads, body)
		mu.Unlock()

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	log, closer, err := logger.New(
		config.AppConfig{Name: "svc", Env: "development"},
		config.ObservabilityConfig{
			Loki: config.LokiConfig{
				Enabled:   true,
				URL:       server.URL,
				BatchWait: 10 * time.Millisecond,
			},
		},
	)
	if err != nil {
		t.Fatalf("expected logger init to succeed: %v", err)
	}

	log.Info("hello loki")
	if err := closer.Close(); err != nil {
		t.Fatalf("expected Loki closer to flush successfully: %v", err)
	}
	_ = log.Sync()

	mu.Lock()
	defer mu.Unlock()

	if len(payloads) == 0 {
		t.Fatal("expected at least one Loki payload")
	}

	var decoded map[string]any
	if err := json.Unmarshal(payloads[0], &decoded); err != nil {
		t.Fatalf("expected valid Loki push payload: %v", err)
	}
}
