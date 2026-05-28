package middlewares

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

func Sentry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				http.Error(
					w,
					"Internal Server Error",
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
