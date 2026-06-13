package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	apperrors "github.com/username/project-name/internal/errors"
	"github.com/username/project-name/internal/request"
	"github.com/username/project-name/internal/response"
)

var validationEngine = validator.New()

type ValidationFieldError struct {
	Field       string `json:"field"`
	Message     string `json:"message"`
	Explanation string `json:"explanation,omitempty"`
	Rule        string `json:"rule,omitempty"`
}

type ValidationErrorDetails struct {
	Summary string                 `json:"summary"`
	Fields  []ValidationFieldError `json:"fields,omitempty"`
}

func ValidationGuard[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload := new(T)
			if appErr := decodeAndValidateBody(r, payload); appErr != nil {
				response.Error(w, r, appErr)
				return
			}

			ctx := request.WithValidatedPayload(r.Context(), payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func decodeAndValidateBody[T any](r *http.Request, payload *T) *apperrors.AppError {
	if r.Body == nil {
		return validationError("request body is required", []ValidationFieldError{
			{
				Field:       "body",
				Message:     "request body is required",
				Explanation: "Provide a JSON request body for this endpoint.",
				Rule:        "required",
			},
		})
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(payload); err != nil {
		return mapDecodeError(err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return validationError("request body must contain a single JSON object", []ValidationFieldError{
			{
				Field:       "body",
				Message:     "request body must contain a single JSON object",
				Explanation: "Remove trailing JSON values or duplicate payload fragments.",
				Rule:        "single_object",
			},
		})
	}

	if err := validationEngine.Struct(payload); err != nil {
		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			return validationError("request validation failed", mapValidationErrors[T](fieldErrors))
		}

		return validationError("request validation failed", []ValidationFieldError{
			{
				Field:       "body",
				Message:     "request validation failed",
				Explanation: "The request body did not satisfy the expected contract.",
				Rule:        "validation",
			},
		})
	}

	return nil
}

func mapDecodeError(err error) *apperrors.AppError {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return validationError("invalid JSON syntax", []ValidationFieldError{
			{
				Field:       "body",
				Message:     "invalid JSON syntax",
				Explanation: fmt.Sprintf("The JSON payload is malformed near byte offset %d.", syntaxErr.Offset),
				Rule:        "json_syntax",
			},
		})
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := typeErr.Field
		if field == "" {
			field = "body"
		}
		return validationError("invalid field type", []ValidationFieldError{
			{
				Field:       normalizeFieldName(field),
				Message:     fmt.Sprintf("must be a valid %s", typeErr.Type.String()),
				Explanation: fmt.Sprintf("The value provided for %s does not match the expected JSON type.", normalizeFieldName(field)),
				Rule:        "type",
			},
		})
	}

	message := err.Error()
	if strings.HasPrefix(message, "json: unknown field ") {
		field := strings.Trim(message[len("json: unknown field "):], "\"")
		return validationError("unknown field in request body", []ValidationFieldError{
			{
				Field:       field,
				Message:     "is not allowed",
				Explanation: "Remove this field from the request body because it is not part of the accepted contract.",
				Rule:        "unknown_field",
			},
		})
	}

	if errors.Is(err, io.EOF) {
		return validationError("request body is required", []ValidationFieldError{
			{
				Field:       "body",
				Message:     "request body is required",
				Explanation: "Provide a JSON request body for this endpoint.",
				Rule:        "required",
			},
		})
	}

	return validationError("invalid request body", []ValidationFieldError{
		{
			Field:       "body",
			Message:     "could not be parsed",
			Explanation: "Ensure the request body is valid JSON and matches the expected schema.",
			Rule:        "invalid_body",
		},
	})
}

func mapValidationErrors[T any](fieldErrors validator.ValidationErrors) []ValidationFieldError {
	results := make([]ValidationFieldError, 0, len(fieldErrors))
	typeInfo := reflectedType[T]()

	for _, fieldErr := range fieldErrors {
		fieldName := jsonFieldName(typeInfo, fieldErr.StructField())
		results = append(results, ValidationFieldError{
			Field:       fieldName,
			Message:     validationMessage(fieldErr),
			Explanation: validationExplanation(fieldErr, fieldName),
			Rule:        fieldErr.Tag(),
		})
	}

	return results
}

func validationMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", fieldErr.Param())
	default:
		return "is invalid"
	}
}

func validationExplanation(fieldErr validator.FieldError, field string) string {
	switch fieldErr.Tag() {
	case "required":
		return fmt.Sprintf("Provide a value for %s.", field)
	case "email":
		return fmt.Sprintf("%s must contain a valid email address.", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long.", field, fieldErr.Param())
	default:
		return fmt.Sprintf("%s did not satisfy the validation rule %s.", field, fieldErr.Tag())
	}
}

func validationError(summary string, fields []ValidationFieldError) *apperrors.AppError {
	return apperrors.Validation(summary, ValidationErrorDetails{
		Summary: summary,
		Fields:  fields,
	})
}

func reflectedType[T any]() reflect.Type {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	if typ.Kind() == reflect.Pointer {
		return typ.Elem()
	}
	return typ
}

func jsonFieldName(typ reflect.Type, field string) string {
	if typ.Kind() != reflect.Struct {
		return normalizeFieldName(field)
	}

	if structField, ok := typ.FieldByName(field); ok {
		tag := structField.Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if name != "" && name != "-" {
			return name
		}
	}

	return normalizeFieldName(field)
}

func normalizeFieldName(field string) string {
	return strings.ToLower(field)
}
