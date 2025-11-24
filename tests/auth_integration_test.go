package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/dal/pgtodolist"
	"github.com/macesz/todo-go/dal/pguser"
	"github.com/macesz/todo-go/delivery/web/middlewares"
	"github.com/macesz/todo-go/delivery/web/todolist"
	"github.com/macesz/todo-go/domain"
	todolistservice "github.com/macesz/todo-go/services/todolist"
	userservice "github.com/macesz/todo-go/services/user"
	"github.com/macesz/todo-go/tests/testutils"
	"github.com/stretchr/testify/require"
)

// setupAuthTestServer builds the router with the REAL auth middleware chain
func setupAuthTestServer(t *testing.T) (*chi.Mux, *testutils.TestContainer, *jwtauth.JWTAuth) {
	t.Helper()

	tc := testutils.SetupTestDB(t)

	// 1. Setup Auth Strategy
	tokenAuth := testutils.SetupTestAuth()

	// 2. Setup Services
	userStore := pguser.CreateStore(tc.DB)
	listStore := pgtodolist.CreateStore(tc.DB)

	userSvc := userservice.NewUserService(userStore)
	listSvc := todolistservice.NewTodoListService(listStore)

	listHandlers := todolist.NewHandlers(listSvc, userSvc)

	// 3. Setup Router with Middleware Chain (Copying from your server.go)
	r := chi.NewRouter()

	// === PROTECTED ROUTES CONFIGURATION ===
	r.Group(func(r chi.Router) {
		// IMPORTANT: This is the chain we are testing
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middlewares.Authenticator)
		r.Use(middlewares.UserContext)

		r.Get("/lists", listHandlers.List)
		r.Get("/lists/{id}", listHandlers.GetListByID)
	})

	return r, tc, tokenAuth
}

func Test_Auth_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	router, tc, tokenAuth := setupAuthTestServer(t)
	defer testutils.CleanupDB(t, tc.DB)

	ctx := t.Context()
	userSvc := userservice.NewUserService(pguser.CreateStore(tc.DB))

	// Create two users
	user1, _ := userSvc.CreateUser(ctx, "User One", "u1@example.com", "pass")
	user2, _ := userSvc.CreateUser(ctx, "User Two", "u2@example.com", "pass")

	// Generate JWT Tokens for them
	token1, err := testutils.GenerateTestToken(tokenAuth, user1)
	require.NoError(t, err)

	token2, err := testutils.GenerateTestToken(tokenAuth, user2)
	require.NoError(t, err)

	// Create a list for User 1
	listID, _ := testutils.GivenTodoLists(t, tc.DB, domain.TodoList{UserID: user1.ID, Title: "User 1 Secrets"})

	t.Run("Authentication (Who are you?)", func(t *testing.T) {

		t.Run("Access protected route without token -> 401", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/lists", nil)
			// No Authorization header set
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusUnauthorized, rr.Code)
		})

		t.Run("Access protected route with invalid token -> 401", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/lists", nil)
			req.Header.Set("Authorization", "Bearer invalid-rubbish-token")
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusUnauthorized, rr.Code)
		})

		t.Run("Access protected route with valid token -> 200", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/lists", nil)
			// Set valid Authorization header
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token1))
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)
		})
	})

	t.Run("Authorization (Access Control)", func(t *testing.T) {

		t.Run("User 1 can see their own list", func(t *testing.T) {
			url := fmt.Sprintf("/lists/%d", listID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token1))

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)
		})

		t.Run("User 2 CANNOT see User 1's list -> 404 or 403", func(t *testing.T) {
			// User 2 has a valid token, but requests User 1's data
			url := fmt.Sprintf("/lists/%d", listID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token2))

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Should be Not Found (secure) or Forbidden
			require.Equal(t, http.StatusNotFound, rr.Code)
		})
	})
}
