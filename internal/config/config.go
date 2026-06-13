package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

const defaultJWTSecret = "change-me-before-production-change-me-before-production"

type Config struct {
	App           AppConfig           `validate:"required"`
	Server        ServerConfig        `validate:"required"`
	Database      DatabaseConfig      `validate:"required"`
	Auth          AuthConfig          `validate:"required"`
	Cache         CacheConfig         `validate:"required"`
	Mailer        MailerConfig        `validate:"required"`
	Storage       StorageConfig       `validate:"required"`
	CORS          CORSConfig          `validate:"required"`
	RateLimit     RateLimitConfig     `validate:"required"`
	Observability ObservabilityConfig `validate:"required"`
}

type AppConfig struct {
	Name string `validate:"required"`
	Env  string `validate:"required,oneof=development staging production"`
}

type ServerConfig struct {
	Port              string        `validate:"required,numeric"`
	ReadHeaderTimeout time.Duration `validate:"gt=0"`
	ReadTimeout       time.Duration `validate:"gt=0"`
	WriteTimeout      time.Duration `validate:"gt=0"`
	IdleTimeout       time.Duration `validate:"gt=0"`
	ShutdownTimeout   time.Duration `validate:"gt=0"`
	BodyLimitBytes    int64         `validate:"gt=0"`
}

type DatabaseConfig struct {
	URL             string        `validate:"required"`
	MaxOpenConns    int           `validate:"gte=1"`
	MaxIdleConns    int           `validate:"gte=1"`
	ConnMaxLifetime time.Duration `validate:"gt=0"`
	ConnMaxIdleTime time.Duration `validate:"gt=0"`
}

type AuthConfig struct {
	JWTSecret       string        `validate:"required,min=32"`
	TokenTTL        time.Duration `validate:"gt=0"`
	RefreshTokenTTL time.Duration `validate:"gt=0"`
	ResetTokenTTL   time.Duration `validate:"gt=0"`
	Issuer          string        `validate:"required"`
}

type MailerConfig struct {
	Provider     string `validate:"required,oneof=noop log smtp resend"`
	FromAddress  string `validate:"required,email"`
	SMTPHost     string
	SMTPPort     int `validate:"gte=0"`
	SMTPUsername string
	SMTPPassword string
	ResendAPIKey string
	ResendURL    string `validate:"omitempty,url"`
}

type StorageConfig struct {
	Provider string `validate:"required,oneof=noop local"`
	LocalDir string
}

type CacheConfig struct {
	RedisURL   string        `validate:"omitempty,url"`
	DefaultTTL time.Duration `validate:"gt=0"`
}

type CORSConfig struct {
	AllowedOrigins   []string `validate:"min=1,dive,required"`
	AllowedMethods   []string `validate:"min=1,dive,required"`
	AllowedHeaders   []string `validate:"min=1,dive,required"`
	AllowCredentials bool
	MaxAge           int `validate:"gte=0"`
}

type RateLimitConfig struct {
	Enabled  bool
	Requests int           `validate:"gte=1"`
	Window   time.Duration `validate:"gt=0"`
}

type ObservabilityConfig struct {
	SentryDSN   string `validate:"omitempty,url"`
	PosthogKey  string
	PosthogHost string        `validate:"omitempty,url"`
	Metrics     MetricsConfig `validate:"required"`
	Loki        LokiConfig    `validate:"required"`
}

type MetricsConfig struct {
	Enabled   bool
	Path      string `validate:"required"`
	Namespace string `validate:"required"`
}

type LokiConfig struct {
	Enabled           bool
	URL               string `validate:"omitempty,url"`
	TenantID          string
	BasicAuthUsername string
	BasicAuthPassword string
	BatchWait         time.Duration `validate:"gt=0"`
	Labels            map[string]string
}

func Load() (*Config, error) {
	readHeaderTimeout, err := parseDuration("SERVER_READ_HEADER_TIMEOUT", "5s")
	if err != nil {
		return nil, err
	}
	readTimeout, err := parseDuration("SERVER_READ_TIMEOUT", "15s")
	if err != nil {
		return nil, err
	}
	writeTimeout, err := parseDuration("SERVER_WRITE_TIMEOUT", "30s")
	if err != nil {
		return nil, err
	}
	idleTimeout, err := parseDuration("SERVER_IDLE_TIMEOUT", "60s")
	if err != nil {
		return nil, err
	}
	shutdownTimeout, err := parseDuration("SERVER_SHUTDOWN_TIMEOUT", "10s")
	if err != nil {
		return nil, err
	}
	bodyLimitBytes, err := parseInt64("SERVER_BODY_LIMIT_BYTES", "1048576")
	if err != nil {
		return nil, err
	}
	maxOpenConns, err := parseInt("DB_MAX_OPEN_CONNS", "25")
	if err != nil {
		return nil, err
	}
	maxIdleConns, err := parseInt("DB_MAX_IDLE_CONNS", "25")
	if err != nil {
		return nil, err
	}
	connMaxLifetime, err := parseDuration("DB_CONN_MAX_LIFETIME", "30m")
	if err != nil {
		return nil, err
	}
	connMaxIdleTime, err := parseDuration("DB_CONN_MAX_IDLE_TIME", "15m")
	if err != nil {
		return nil, err
	}
	tokenTTL, err := parseDuration("JWT_TOKEN_TTL", "24h")
	if err != nil {
		return nil, err
	}
	refreshTokenTTL, err := parseDuration("JWT_REFRESH_TOKEN_TTL", "168h")
	if err != nil {
		return nil, err
	}
	resetTokenTTL, err := parseDuration("JWT_RESET_TOKEN_TTL", "1h")
	if err != nil {
		return nil, err
	}
	cacheTTL, err := parseDuration("CACHE_DEFAULT_TTL", "15m")
	if err != nil {
		return nil, err
	}
	smtpPort, err := parseInt("MAILER_SMTP_PORT", "1025")
	if err != nil {
		return nil, err
	}
	allowCredentials, err := parseBool("CORS_ALLOW_CREDENTIALS", "true")
	if err != nil {
		return nil, err
	}
	maxAge, err := parseInt("CORS_MAX_AGE", "300")
	if err != nil {
		return nil, err
	}
	rateLimitEnabled, err := parseBool("RATE_LIMIT_ENABLED", "true")
	if err != nil {
		return nil, err
	}
	rateLimitRequests, err := parseInt("RATE_LIMIT_REQUESTS", "60")
	if err != nil {
		return nil, err
	}
	rateLimitWindow, err := parseDuration("RATE_LIMIT_WINDOW", "1m")
	if err != nil {
		return nil, err
	}
	metricsEnabled, err := parseBool("METRICS_ENABLED", "true")
	if err != nil {
		return nil, err
	}
	lokiEnabled, err := parseBool("LOKI_ENABLED", "false")
	if err != nil {
		return nil, err
	}
	lokiBatchWait, err := parseDuration("LOKI_BATCH_WAIT", "5s")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "go-backend-boilerplate"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Server: ServerConfig{
			Port:              getEnv("PORT", "8080"),
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
			ShutdownTimeout:   shutdownTimeout,
			BodyLimitBytes:    bodyLimitBytes,
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", ""),
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
			ConnMaxIdleTime: connMaxIdleTime,
		},
		Auth: AuthConfig{
			JWTSecret:       getEnv("JWT_SECRET", defaultJWTSecret),
			TokenTTL:        tokenTTL,
			RefreshTokenTTL: refreshTokenTTL,
			ResetTokenTTL:   resetTokenTTL,
			Issuer:          getEnv("JWT_ISSUER", "go-backend-boilerplate"),
		},
		Cache: CacheConfig{
			RedisURL:   getEnv("REDIS_URL", ""),
			DefaultTTL: cacheTTL,
		},
		Mailer: MailerConfig{
			Provider:     getEnv("MAILER_PROVIDER", "resend"),
			FromAddress:  getEnv("MAILER_FROM_ADDRESS", "noreply@example.com"),
			SMTPHost:     getEnv("MAILER_SMTP_HOST", "localhost"),
			SMTPPort:     smtpPort,
			SMTPUsername: getEnv("MAILER_SMTP_USERNAME", ""),
			SMTPPassword: getEnv("MAILER_SMTP_PASSWORD", ""),
			ResendAPIKey: getEnv("RESEND_API_KEY", ""),
			ResendURL:    getEnv("RESEND_URL", "https://api.resend.com"),
		},
		Storage: StorageConfig{
			Provider: getEnv("STORAGE_PROVIDER", "noop"),
			LocalDir: getEnv("STORAGE_LOCAL_DIR", "./tmp/storage"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
			AllowedMethods:   splitCSV(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")),
			AllowedHeaders:   splitCSV(getEnv("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-Request-ID")),
			AllowCredentials: allowCredentials,
			MaxAge:           maxAge,
		},
		RateLimit: RateLimitConfig{
			Enabled:  rateLimitEnabled,
			Requests: rateLimitRequests,
			Window:   rateLimitWindow,
		},
		Observability: ObservabilityConfig{
			SentryDSN:   getEnv("SENTRY_DSN", ""),
			PosthogKey:  getEnv("POSTHOG_KEY", ""),
			PosthogHost: getEnv("POSTHOG_HOST", "https://app.posthog.com"),
			Metrics: MetricsConfig{
				Enabled:   metricsEnabled,
				Path:      getEnv("METRICS_PATH", "/metrics"),
				Namespace: getEnv("METRICS_NAMESPACE", "go_backend"),
			},
			Loki: LokiConfig{
				Enabled:           lokiEnabled,
				URL:               getEnv("LOKI_URL", ""),
				TenantID:          getEnv("LOKI_TENANT_ID", ""),
				BasicAuthUsername: getEnv("LOKI_BASIC_AUTH_USERNAME", ""),
				BasicAuthPassword: getEnv("LOKI_BASIC_AUTH_PASSWORD", ""),
				BatchWait:         lokiBatchWait,
				Labels:            parseMap(getEnv("LOKI_LABELS", "")),
			},
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}

	if c.App.Env == "production" && c.Auth.JWTSecret == defaultJWTSecret {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}

	if !strings.HasPrefix(c.Observability.Metrics.Path, "/") {
		return fmt.Errorf("METRICS_PATH must start with '/'")
	}

	if c.Observability.Loki.Enabled && c.Observability.Loki.URL == "" {
		return fmt.Errorf("LOKI_URL must be set when Loki logging is enabled")
	}

	if c.Mailer.Provider == "smtp" && c.Mailer.SMTPHost == "" {
		return fmt.Errorf("MAILER_SMTP_HOST must be set when smtp mailer is enabled")
	}

	if c.Mailer.Provider == "resend" && c.Mailer.ResendAPIKey == "" {
		return fmt.Errorf("RESEND_API_KEY must be set when resend mailer is enabled")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func parseMap(value string) map[string]string {
	result := map[string]string{}
	if strings.TrimSpace(value) == "" {
		return result
	}

	pairs := strings.Split(value, ",")
	for _, pair := range pairs {
		key, rawValue, ok := strings.Cut(strings.TrimSpace(pair), "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		rawValue = strings.TrimSpace(rawValue)
		if key == "" || rawValue == "" {
			continue
		}

		result[key] = rawValue
	}

	return result
}

func parseDuration(key, fallback string) (time.Duration, error) {
	value := getEnv(key, fallback)
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s: %w", key, err)
	}
	return duration, nil
}

func parseInt(key, fallback string) (int, error) {
	value := getEnv(key, fallback)
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid int for %s: %w", key, err)
	}
	return parsed, nil
}

func parseInt64(key, fallback string) (int64, error) {
	value := getEnv(key, fallback)
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int64 for %s: %w", key, err)
	}
	return parsed, nil
}

func parseBool(key, fallback string) (bool, error) {
	value := getEnv(key, fallback)
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid bool for %s: %w", key, err)
	}
	return parsed, nil
}
