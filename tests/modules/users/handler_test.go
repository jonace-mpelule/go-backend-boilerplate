package users_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/username/project-name/internal/modules/users"
)

type stubService struct {
	listFn func(context.Context) ([]users.UserResponse, error)
}

func (s stubService) ListUsers(ctx context.Context) ([]users.UserResponse, error) {
	return s.listFn(ctx)
}

func TestHandlerListUsersSuccess(t *testing.T) {
	handler := users.NewHandler(stubService{
		listFn: func(context.Context) ([]users.UserResponse, error) {
			return []users.UserResponse{{ID: "user-1", Email: "user@example.com"}}, nil
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()

	handler.ListUsers(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestHandlerListUsersFailure(t *testing.T) {
	handler := users.NewHandler(stubService{
		listFn: func(context.Context) ([]users.UserResponse, error) {
			return nil, assertErr{}
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()

	handler.ListUsers(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

type assertErr struct{}

func (assertErr) Error() string { return "boom" }
