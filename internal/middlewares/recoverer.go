package middlewares

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/username/project-name/internal/errors"
	"github.com/username/project-name/internal/response"
	"go.uber.org/zap"
)

func Recoverer(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					sentry.CurrentHub().Recover(recovered)
					logger.Error(
						"request panic recovered",
						zap.Any("panic", recovered),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
					)
					response.Error(w, r, errors.Internal("internal server error"))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
