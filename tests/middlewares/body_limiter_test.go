package middlewares_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/username/project-name/internal/middlewares"
)

func TestBodyLimiterRestrictsReads(t *testing.T) {
	handler := middlewares.BodyLimiter(4)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err == nil {
			t.Fatal("expected body read to fail when body exceeds limit")
		}
		w.WriteHeader(http.StatusRequestEntityTooLarge)
	}))

	request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("12345"))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d", recorder.Code)
	}
}
