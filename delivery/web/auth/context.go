package auth

import "context"

type contextKey string

const userContextKey contextKey = "userContext"

type UserContext struct {
	ID    int64
	Name  string
	Email string
}

// NewUserContext - Create from JWT claims
func NewUserContext(claims *userClaims) *UserContext {
	return &UserContext{
		ID:    claims.UserID,
		Email: claims.Email,
		Name:  claims.Name,
	}
}

// AddToContext - Store user context in request context
func (u *UserContext) AddToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// FromContext - Retrieve user context from request context
func UserFromContext(ctx context.Context) (*UserContext, bool) {
	userCtx, ok := ctx.Value(userContextKey).(*UserContext)
	return userCtx, ok
}
