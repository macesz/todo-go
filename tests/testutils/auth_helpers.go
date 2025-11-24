package testutils

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/domain"
)

// GenerateTestToken creates a real signed JWT string for a test user
func GenerateTestToken(tokenAuth *jwtauth.JWTAuth, user *domain.User) (string, error) {
	// Create claims using your existing logic
	claims := auth.NewUserClaims(user, time.Hour)

	// Encode using the library
	_, tokenString, err := tokenAuth.Encode(claims.ToMap())
	return tokenString, err
}

// SetupTestAuth creates the JWT for testing
func SetupTestAuth() *jwtauth.JWTAuth {
	return auth.CreateTokenAuth("my-super-secret-test-key-12345")
}
