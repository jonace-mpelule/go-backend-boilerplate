package errors

import "net/http"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Details any    `json:"details,omitempty"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func Wrap(err error, code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

func Validation(message string, details any) *AppError {
	appErr := New("validation_error", message, http.StatusBadRequest)
	appErr.Details = details
	return appErr
}

func Unauthorized(message string) *AppError {
	return New("unauthorized", message, http.StatusUnauthorized)
}

func Forbidden(message string) *AppError {
	return New("forbidden", message, http.StatusForbidden)
}

func NotFound(message string) *AppError {
	return New("not_found", message, http.StatusNotFound)
}

func Conflict(message string) *AppError {
	return New("conflict", message, http.StatusConflict)
}

func Internal(message string) *AppError {
	return New("internal_error", message, http.StatusInternalServerError)
}
