package middlewares

import (
	"net/http"

	apperrors "github.com/username/project-name/internal/errors"
	"github.com/username/project-name/internal/response"
)

type PermissionMode int

const (
	RequireAll PermissionMode = iota
	RequireAny
)

func RequirePermission(mode PermissionMode, required ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				response.Error(w, r, apperrors.Unauthorized("missing auth context"))
				return
			}

			if claims.Role == "super_admin" {
				next.ServeHTTP(w, r)
				return
			}

			userPerms := make(map[string]bool, len(claims.Permissions))
			for _, permission := range claims.Permissions {
				userPerms[permission] = true
			}

			switch mode {
			case RequireAll:
				for _, permission := range required {
					if !userPerms[permission] {
						response.Error(w, r, apperrors.Forbidden("missing required permission"))
						return
					}
				}
				next.ServeHTTP(w, r)
			case RequireAny:
				for _, permission := range required {
					if userPerms[permission] {
						next.ServeHTTP(w, r)
						return
					}
				}
				response.Error(w, r, apperrors.Forbidden("missing required permission"))
			default:
				response.Error(w, r, apperrors.Forbidden("invalid permission mode"))
			}
		})
	}
}
