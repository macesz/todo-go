package middlewares

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/delivery/web/utils"
)

func UnloggedInRedirector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, _ := jwtauth.FromContext(r.Context())

		if token == nil || jwt.Validate(token) != nil {
			http.Redirect(w, r, "/login", 302)
		}

		next.ServeHTTP(w, r)
	})
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, utils.JsonError(err), http.StatusUnauthorized)
			return
		}

		if token == nil {
			http.Error(w, utils.JsonError(err), http.StatusUnauthorized)
			return
		}

		claim := token.PrivateClaims()
		user_id, ok := claim["user_id"].(float64)
		if !ok {
			err := errors.New("invalid user id in token")
			http.Error(w, utils.JsonError(err), http.StatusUnauthorized)
		}

		// CHECK USER ID IS VALID
		if user_id <= 0 {
			err := errors.New("invalid user id in token")
			http.Error(w, utils.JsonError(err), http.StatusUnauthorized)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, utils.JsonError(err), http.StatusUnauthorized)
			return
		}

		privateClaims := token.PrivateClaims()

		// Extract user information from token claims
		claims, err := auth.ClaimsFromToken(privateClaims)
		if err != nil {
			http.Error(w, utils.JsonError(err), http.StatusUnauthorized)
			return
		}

		userContext := auth.UserContext{
			ID:    claims.UserID,
			Name:  claims.Name,
			Email: claims.Email,
		}

		ctx := userContext.AddToContext(r.Context())

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
