package middlewares

import (
	"net/http"

	"github.com/username/project-name/internal/utils"
)

type PermissionMode int

const (
	RequireAll PermissionMode = iota
	RequireAny
)

func RequirePermission(
	mode PermissionMode,
	required ...string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			claims, ok := r.Context().
				Value(UserContextKey).(*utils.Claims)

			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if claims.Role == "super_admin" {
				next.ServeHTTP(w, r)
				return
			}

			userPerms := make(map[string]bool)

			for _, p := range claims.Permissions {
				userPerms[p] = true
			}

			switch mode {

			case RequireAll:
				for _, p := range required {
					if !userPerms[p] {
						http.Error(w, "forbidden", http.StatusForbidden)
						return
					}
				}
				next.ServeHTTP(w, r)

			case RequireAny:

				for _, p := range required {
					if userPerms[p] {
						next.ServeHTTP(w, r)
						return
					}
				}

				http.Error(w, "forbidden", http.StatusForbidden)
			}
		})
	}
}
