package auth

import (
	"errors"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/domain"
)
// CreateTokenAuth - Initialize JWT Auth with given secret, factory function
func CreateTokenAuth(secret string) *jwtauth.JWTAuth {
	// JWT Auth setup with HS256 and secret from config
	return jwtauth.New("HS256", []byte(secret), nil)
}

// JWT Claims struct, made private to the auth package -> encapsulation
type userClaims struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	EXP    int64  `json:"exp"`
}

// NewUserClaims - Convert domain.User to JWT claims

func NewUserClaims(u *domain.User, expiresIn time.Duration) userClaims {
	return userClaims{
		UserID: u.ID,
		Name:   u.Name,
		Email:  u.Email,
		EXP:    time.Now().Add(expiresIn).Unix(),
	}
}

// ToMap - Convert to map for jwtauth library
func (c userClaims) ToMap() map[string]any {
	return map[string]any{
		"user_id": c.UserID,
		"name":    c.Name,
		"email":   c.Email,
		"exp":     c.EXP,
	}
}

// ClaimsFromToken - Extract and validate claims from JWT token
// IMPORTANT: JWT stores numbers as float64, not int64!
// Extratct Claims from JWT token private claims
func ClaimsFromToken(claims map[string]any) (*userClaims, error) {
	userId, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user id in token")
	}
	name, ok := claims["name"].(string)
	if !ok {
		return nil, errors.New("invalid name in token")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email in token")
	}

	//Removed manual expiration extraction (JWT library handles this)
	return &userClaims{
		UserID: int64(userId),
		Name:   name,
		Email:  email,
	}, nil
}
