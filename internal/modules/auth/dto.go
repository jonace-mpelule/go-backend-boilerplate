package auth

import (
	"time"

	"github.com/username/project-name/internal/utils"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type UserRecord struct {
	ID                string
	Email             string
	PasswordHash      string
	Role              string
	Permissions       []string
	ResetTokenHash    *string
	ResetTokenExpires *time.Time
}

type RefreshSessionRecord struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type TokenPairResponse struct {
	AccessToken  string        `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string        `json:"refresh_token" example:"refresh-token-value"`
	ExpiresIn    time.Duration `json:"expires_in" swaggertype:"integer" example:"3600"`
}

type AuthResponse struct {
	User   ProfileResponse   `json:"user"`
	Tokens TokenPairResponse `json:"tokens"`
}

type ProfileResponse struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type JWTVerifier interface {
	Verify(tokenString string) (*utils.Claims, error)
}
