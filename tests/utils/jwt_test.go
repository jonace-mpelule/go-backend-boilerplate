package utils_test

import (
	"testing"
	"time"

	"github.com/username/project-name/internal/utils"
)

func TestJWTGenerateAndVerify(t *testing.T) {
	helper := utils.NewJWT("12345678901234567890123456789012", "test-suite", time.Hour)

	token, err := helper.Generate("user-1", "admin", []string{"user:read"})
	if err != nil {
		t.Fatalf("expected token generation to succeed: %v", err)
	}

	claims, err := helper.Verify(token)
	if err != nil {
		t.Fatalf("expected token verification to succeed: %v", err)
	}

	if claims.UserID != "user-1" || claims.Role != "admin" {
		t.Fatal("expected claims to round-trip")
	}
}
