package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthGuard(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			authHeader := r.Header.Get(
				"Authorization",
			)

			if authHeader == "" {
				http.Error(
					w,
					"unauthorized",
					http.StatusUnauthorized,
				)
				return
			}

			tokenString := strings.TrimPrefix(
				authHeader,
				"Bearer ",
			)

			token, err := jwt.Parse(
				tokenString,
				func(token *jwt.Token) (
					any,
					error,
				) {
					return []byte(secret), nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(
					w,
					"invalid token",
					http.StatusUnauthorized,
				)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				UserContextKey,
				token.Claims,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
