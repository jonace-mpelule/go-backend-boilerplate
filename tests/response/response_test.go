package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/username/project-name/internal/errors"
	"github.com/username/project-name/internal/response"
)

func TestSuccess(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	response.Success(recorder, request, http.StatusOK, map[string]string{"status": "ok"})

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload response.Envelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if !payload.Success {
		t.Fatal("expected success response")
	}
}

func TestError(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	response.Error(recorder, request, apperrors.Forbidden("forbidden"))

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}

	var payload response.Envelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if payload.Success {
		t.Fatal("expected error response")
	}

	if payload.Error == nil || payload.Error.Code != "forbidden" {
		t.Fatal("expected forbidden error body")
	}
}
