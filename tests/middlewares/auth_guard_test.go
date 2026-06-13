package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/username/project-name/internal/middlewares"
	"github.com/username/project-name/internal/utils"
)

func TestAuthGuardInjectsClaims(t *testing.T) {
	jwtHelper := utils.NewJWT("12345678901234567890123456789012", "test-suite", time.Hour)
	token, err := jwtHelper.Generate("user-1", "admin", []string{"user:read"})
	if err != nil {
		t.Fatalf("expected token generation to succeed: %v", err)
	}

	handler := middlewares.AuthGuard(jwtHelper)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middlewares.ClaimsFromContext(r.Context())
		if !ok {
			t.Fatal("expected claims in context")
		}
		if claims.UserID != "user-1" {
			t.Fatalf("expected user-1, got %s", claims.UserID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestAuthGuardRejectsInvalidToken(t *testing.T) {
	jwtHelper := utils.NewJWT("12345678901234567890123456789012", "test-suite", time.Hour)
	handler := middlewares.AuthGuard(jwtHelper)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	request.Header.Set("Authorization", "Bearer invalid")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
