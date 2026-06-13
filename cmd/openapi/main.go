package main

import (
	"encoding/json"
	"os"
)

func main() {
	spec := map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "Go Backend Boilerplate API",
			"version":     "1.0.0",
			"description": "Generated OpenAPI document for the boilerplate auth and starter modules.",
		},
		"servers": []map[string]string{
			{"url": "http://localhost:8080"},
		},
		"paths": map[string]any{
			"/health": map[string]any{
				"get": successOperation("Readiness health check", "HealthStatus"),
			},
			"/health/live": map[string]any{
				"get": successOperation("Liveness health check", "HealthStatus"),
			},
			"/health/ready": map[string]any{
				"get": successOperation("Readiness health check", "HealthStatus"),
			},
			"/metrics": map[string]any{
				"get": map[string]any{
					"summary": "Prometheus metrics endpoint",
					"responses": map[string]any{
						"200": map[string]any{
							"description": "Prometheus metrics output",
						},
					},
				},
			},
			"/auth/register": map[string]any{
				"post": requestOperation("Register a new user", "RegisterRequest", "AuthResponse"),
			},
			"/auth/login": map[string]any{
				"post": requestOperation("Login with email and password", "LoginRequest", "AuthResponse"),
			},
			"/auth/refresh": map[string]any{
				"post": requestOperation("Refresh access token", "RefreshRequest", "AuthResponse"),
			},
			"/auth/forgot-password": map[string]any{
				"post": requestOperation("Request password reset", "ForgotPasswordRequest", "MessageResponse"),
			},
			"/auth/reset-password": map[string]any{
				"post": requestOperation("Reset password with token", "ResetPasswordRequest", "MessageResponse"),
			},
			"/auth/me": map[string]any{
				"get": securedOperation("Current authenticated user", "ProfileResponse"),
			},
			"/users/": map[string]any{
				"get": securedOperation("List users", "UsersListResponse"),
			},
		},
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
			"schemas": map[string]any{
				"Envelope": objectSchema(map[string]any{
					"success": map[string]string{"type": "boolean"},
				}, []string{"success"}),
				"ErrorBody": objectSchema(map[string]any{
					"code":    map[string]string{"type": "string"},
					"message": map[string]string{"type": "string"},
				}, []string{"code", "message"}),
				"RegisterRequest": objectSchema(map[string]any{
					"email":    map[string]string{"type": "string", "format": "email"},
					"password": map[string]string{"type": "string", "minLength": "8"},
				}, []string{"email", "password"}),
				"LoginRequest": objectSchema(map[string]any{
					"email":    map[string]string{"type": "string", "format": "email"},
					"password": map[string]string{"type": "string", "minLength": "8"},
				}, []string{"email", "password"}),
				"RefreshRequest": objectSchema(map[string]any{
					"refresh_token": map[string]string{"type": "string"},
				}, []string{"refresh_token"}),
				"ForgotPasswordRequest": objectSchema(map[string]any{
					"email": map[string]string{"type": "string", "format": "email"},
				}, []string{"email"}),
				"ResetPasswordRequest": objectSchema(map[string]any{
					"token":        map[string]string{"type": "string"},
					"new_password": map[string]string{"type": "string", "minLength": "8"},
				}, []string{"token", "new_password"}),
				"ProfileResponse": objectSchema(map[string]any{
					"id":          map[string]string{"type": "string"},
					"email":       map[string]string{"type": "string", "format": "email"},
					"role":        map[string]string{"type": "string"},
					"permissions": map[string]any{"type": "array", "items": map[string]string{"type": "string"}},
				}, []string{"id", "email", "role", "permissions"}),
				"TokenPairResponse": objectSchema(map[string]any{
					"access_token":  map[string]string{"type": "string"},
					"refresh_token": map[string]string{"type": "string"},
					"expires_in":    map[string]string{"type": "integer"},
				}, []string{"access_token", "refresh_token", "expires_in"}),
				"AuthResponseData": objectSchema(map[string]any{
					"user":   refSchema("ProfileResponse"),
					"tokens": refSchema("TokenPairResponse"),
				}, []string{"user", "tokens"}),
				"MessageResponse": objectSchema(map[string]any{
					"message": map[string]string{"type": "string"},
				}, []string{"message"}),
				"HealthStatus": objectSchema(map[string]any{
					"status": map[string]string{"type": "string"},
					"checks": map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
				}, []string{"status", "checks"}),
				"UserResponse": objectSchema(map[string]any{
					"id":    map[string]string{"type": "string"},
					"email": map[string]string{"type": "string", "format": "email"},
					"role":  map[string]string{"type": "string"},
				}, []string{"id", "email"}),
				"UsersListResponse": map[string]any{
					"type":  "array",
					"items": refSchema("UserResponse"),
				},
			},
		},
	}

	file, err := os.Create("internal/platform/docs/openapi.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(spec); err != nil {
		panic(err)
	}
}

func successOperation(summary, schema string) map[string]any {
	return map[string]any{
		"summary": summary,
		"responses": map[string]any{
			"200": jsonResponse(schema),
		},
	}
}

func requestOperation(summary, requestSchema, responseSchema string) map[string]any {
	return map[string]any{
		"summary": summary,
		"requestBody": map[string]any{
			"required": true,
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": refSchema(requestSchema),
				},
			},
		},
		"responses": map[string]any{
			"200": jsonResponse(responseSchema),
			"201": jsonResponse(responseSchema),
			"400": errorResponse(),
			"401": errorResponse(),
			"409": errorResponse(),
		},
	}
}

func securedOperation(summary, responseSchema string) map[string]any {
	return map[string]any{
		"summary": summary,
		"security": []map[string][]string{
			{"bearerAuth": {}},
		},
		"responses": map[string]any{
			"200": jsonResponse(responseSchema),
			"401": errorResponse(),
		},
	}
}

func jsonResponse(schema string) map[string]any {
	return map[string]any{
		"description": "Successful response",
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": refSchema(schema),
			},
		},
	}
}

func errorResponse() map[string]any {
	return map[string]any{
		"description": "Error response",
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": refSchema("ErrorBody"),
			},
		},
	}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": properties,
		"required":   required,
	}
}

func refSchema(name string) map[string]any {
	return map[string]any{
		"$ref": "#/components/schemas/" + name,
	}
}
