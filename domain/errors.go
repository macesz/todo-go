package domain

import "errors" // For creating custom errors

// Custom Errors
// These are defined as package-level variables for reuse across services and handlers.
// Check them with errors.Is(err, ErrNotFound) for type-safe handling.

var (
	// ErrNotFound is returned when a todo (or other resource) is not found.
	ErrNotFound = errors.New("todo not found")

	// ErrInvalidTitle is returned for invalid todo titles (e.g., empty or too long).
	ErrInvalidTitle = errors.New("title is required and must be between 1 and 255 characters")

	// ErrInvalidInput is a general error for validation failures.
	ErrInvalidInput = errors.New("invalid input")

	ErrUnauthorized = errors.New("unauthorized") // 401: Missing/invalid auth token
	ErrForbidden    = errors.New("forbidden")    // 403: Valid auth, but no permission

	// ErrDuplicate is returned if a duplicate resource exists (e.g., todo title or user email).
	ErrDuplicate = errors.New("resource already exists")

	// User-specific errors (add more as needed)
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)
