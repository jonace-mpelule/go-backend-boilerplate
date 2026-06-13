package middlewares

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/username/project-name/internal/errors"
	"github.com/username/project-name/internal/response"
	"github.com/username/project-name/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

type TokenVerifier interface {
	Verify(tokenString string) (*utils.Claims, error)
}

func AuthGuard(jwtHelper TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, r, apperrors.Unauthorized("missing authorization header"))
				return
			}

			tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || tokenString == "" {
				response.Error(w, r, apperrors.Unauthorized("invalid authorization header"))
				return
			}

			claims, err := jwtHelper.Verify(tokenString)
			if err != nil {
				response.Error(w, r, apperrors.Unauthorized("invalid token"))
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*utils.Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*utils.Claims)
	return claims, ok
}

func ContextWithClaims(ctx context.Context, claims *utils.Claims) context.Context {
	return context.WithValue(ctx, UserContextKey, claims)
}
