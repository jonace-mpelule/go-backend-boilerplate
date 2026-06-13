package config_test

import (
	"testing"

	"github.com/username/project-name/internal/config"
)

func TestLoadValidConfig(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/template?sslmode=disable")
	t.Setenv("JWT_SECRET", "12345678901234567890123456789012")
	t.Setenv("RESEND_API_KEY", "re_test_123")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected config to load: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Fatalf("expected default port 8080, got %s", cfg.Server.Port)
	}
}

func TestLoadRejectsDefaultProductionJWTSecret(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/template?sslmode=disable")
	t.Setenv("JWT_SECRET", "change-me-before-production-change-me-before-production")
	t.Setenv("RESEND_API_KEY", "re_test_123")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected production config to reject default JWT secret")
	}
}

func TestLoadRejectsInvalidMetricsPath(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/template?sslmode=disable")
	t.Setenv("JWT_SECRET", "12345678901234567890123456789012")
	t.Setenv("RESEND_API_KEY", "re_test_123")
	t.Setenv("METRICS_PATH", "metrics")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected invalid metrics path to fail")
	}
}

func TestLoadRejectsEnabledLokiWithoutURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/template?sslmode=disable")
	t.Setenv("JWT_SECRET", "12345678901234567890123456789012")
	t.Setenv("RESEND_API_KEY", "re_test_123")
	t.Setenv("LOKI_ENABLED", "true")
	t.Setenv("LOKI_URL", "")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected Loki config without URL to fail")
	}
}

func TestLoadRejectsResendProviderWithoutAPIKey(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/template?sslmode=disable")
	t.Setenv("JWT_SECRET", "12345678901234567890123456789012")
	t.Setenv("MAILER_PROVIDER", "resend")
	t.Setenv("RESEND_API_KEY", "")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected resend provider without api key to fail")
	}
}
