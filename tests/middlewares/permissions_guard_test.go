package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/username/project-name/internal/middlewares"
	"github.com/username/project-name/internal/utils"
)

func TestRequirePermissionRequireAny(t *testing.T) {
	handler := middlewares.RequirePermission(middlewares.RequireAny, "user:read", "user:create")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	request = request.WithContext(middlewares.ContextWithClaims(request.Context(), &utils.Claims{
		UserID:      "user-1",
		Role:        "admin",
		Permissions: []string{"user:create"},
	}))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestRequirePermissionAllowsSuperAdmin(t *testing.T) {
	handler := middlewares.RequirePermission(middlewares.RequireAll, "user:read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	request = request.WithContext(middlewares.ContextWithClaims(request.Context(), &utils.Claims{
		UserID: "user-1",
		Role:   "super_admin",
	}))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}
