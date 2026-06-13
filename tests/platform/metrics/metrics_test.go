package metrics_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/username/project-name/internal/platform/metrics"
)

func TestMetricsHandlerReturnsOK(t *testing.T) {
	provider := metrics.New(true, "/metrics", "testsvc")

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	provider.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected metrics handler to return 200, got %d", recorder.Code)
	}
}

func TestMetricsMiddlewareRecordsRequests(t *testing.T) {
	provider := metrics.New(true, "/metrics", "testsvc")
	router := chi.NewRouter()
	router.Use(provider.Middleware)
	router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	router.Handle("/metrics", provider.Handler())

	userRequest := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	userRecorder := httptest.NewRecorder()
	router.ServeHTTP(userRecorder, userRequest)

	metricsRequest := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsRecorder := httptest.NewRecorder()
	router.ServeHTTP(metricsRecorder, metricsRequest)

	body := metricsRecorder.Body.String()
	if !strings.Contains(body, "testsvc_http_requests_total") {
		t.Fatal("expected requests total metric in output")
	}
	if !strings.Contains(body, "testsvc_http_responses_total") {
		t.Fatal("expected responses total metric in output")
	}
	if !strings.Contains(body, "route=\"/users/{id}\"") {
		t.Fatal("expected route pattern label in metrics output")
	}
}
