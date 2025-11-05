package user

import "github.com/go-chi/jwtauth/v5"

// UserHandlers groups HTTP handler functions.
// Like a Java controller class or JS route handler object.
type UserHandlers struct {
	Service   UserService
	TokenAuth *jwtauth.JWTAuth
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(service UserService) *UserHandlers {
	return &UserHandlers{
		Service: service,
	}
}
