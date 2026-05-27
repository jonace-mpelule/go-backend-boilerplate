package config

import (
	"os"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Env   string `validate:"required,oneof=development staging production"`
	Port  string `validate:"required,numeric"`
	DBUrl string `validate:"required,url"`

	SentryDSN  string
	PosthogKey string
}

func Load() (*Config, error) {

	cfg := &Config{
		Env:        getEnv("APP_ENV", "development"),
		Port:       getEnv("PORT", "8080"),
		DBUrl:      getEnv("DATABASE_URL", ""),
		SentryDSN:  getEnv("SENTRY_DSN", ""),
		PosthogKey: getEnv("POSTHOG_KEY", ""),
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil

}

func (c *Config) validate() error {

	validate := validator.New()
	return validate.Struct(c)

}

func getEnv(key, fallback string) string {

	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback

}
