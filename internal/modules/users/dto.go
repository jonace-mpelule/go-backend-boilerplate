package users

import "github.com/username/project-name/internal/utils"

type UserRecord struct {
	ID    string
	Email string
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role,omitempty"`
}

type JWTVerifier interface {
	Verify(tokenString string) (*utils.Claims, error)
}
