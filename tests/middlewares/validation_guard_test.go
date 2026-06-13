package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/username/project-name/internal/middlewares"
	"github.com/username/project-name/internal/request"
)

type validationPayload struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
}

type validationResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details struct {
			Summary string `json:"summary"`
			Fields  []struct {
				Field       string `json:"field"`
				Message     string `json:"message"`
				Explanation string `json:"explanation"`
				Rule        string `json:"rule"`
			} `json:"fields"`
		} `json:"details"`
	} `json:"error"`
}

func TestValidationGuardInjectsValidatedPayload(t *testing.T) {
	handler := middlewares.ValidationGuard[validationPayload]()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := request.ValidatedPayload[validationPayload](r.Context())
		if !ok {
			t.Fatal("expected validated payload in context")
		}
		if payload.Email != "user@example.com" || payload.Name != "Jon" {
			t.Fatalf("unexpected payload: %+v", payload)
		}
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"email":"user@example.com","name":"Jon"}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}
}

func TestValidationGuardRejectsUnknownFields(t *testing.T) {
	handler := middlewares.ValidationGuard[validationPayload]()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"email":"user@example.com","name":"Jon","role":"admin"}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	var body validationResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if body.Error.Code != "validation_error" {
		t.Fatalf("expected validation_error, got %s", body.Error.Code)
	}
	if len(body.Error.Details.Fields) != 1 || body.Error.Details.Fields[0].Field != "role" {
		t.Fatalf("expected unknown role field error, got %+v", body.Error.Details.Fields)
	}
}

func TestValidationGuardRejectsTrailingJSON(t *testing.T) {
	handler := middlewares.ValidationGuard[validationPayload]()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"email":"user@example.com","name":"Jon"}{"email":"extra@example.com"}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestValidationGuardReturnsClearFieldErrors(t *testing.T) {
	handler := middlewares.ValidationGuard[validationPayload]()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"email":"bad-email"}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	var body validationResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if body.Error.Details.Summary != "request validation failed" {
		t.Fatalf("unexpected summary: %s", body.Error.Details.Summary)
	}
	if len(body.Error.Details.Fields) < 2 {
		t.Fatalf("expected multiple field errors, got %+v", body.Error.Details.Fields)
	}
}
