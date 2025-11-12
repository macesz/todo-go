package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/macesz/todo-go/dal/pgtodo"
	"github.com/macesz/todo-go/dal/pguser"
	"github.com/macesz/todo-go/delivery/web/todo"
	"github.com/macesz/todo-go/domain"
	todoservice "github.com/macesz/todo-go/services/todo"
	userservice "github.com/macesz/todo-go/services/user"
	"github.com/macesz/todo-go/tests/testutils"
	"github.com/stretchr/testify/require"
)

// setupTestServer creates a real server with all dependencies
func setupTestServer(t *testing.T) (*chi.Mux, *testutils.TestContainer, int64) {
	t.Helper()

	// Setup database
	tc := testutils.SetupTestDB(t)

	// Create stores
	todoStore := pgtodo.CreateStore(tc.DB)
	userStore := pguser.CreateStore(tc.DB)

	// Create services using constructors
	todoSvc := todoservice.NewTodoService(todoStore)
	userSvc := userservice.NewUserService(userStore)

	// Create test user
	testUser, err := userSvc.CreateUser(t.Context(), "Test User", "test@example.com", "password123")

	require.NoError(t, err)

	// Create handlers using constructor (add this if you don't have it)
	todoHandlers := todo.NewHandlers(todoSvc, userSvc)

	// Setup router
	r := chi.NewRouter()
	r.Get("/todos", todoHandlers.ListTodos)
	r.Post("/todos", todoHandlers.CreateTodo)
	r.Get("/todos/{id}", todoHandlers.GetTodo)
	r.Put("/todos/{id}", todoHandlers.UpdateTodo)
	r.Delete("/todos/{id}", todoHandlers.DeleteTodo)

	return r, tc, testUser.ID
}
func TestTodoHandlers_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Full CRUD Lifecycle", func(t *testing.T) {
		router, tc, userID := setupTestServer(t)
		defer testutils.CleanupDB(t, tc.DB)

		// 1. List todos (should be empty)
		t.Run("List empty todos", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var todos []domain.TodoDTO
			err := json.NewDecoder(rr.Body).Decode(&todos)
			require.NoError(t, err)
			require.Empty(t, todos)
		})

		// 2. Create a todo
		var createdTodo domain.TodoDTO
		t.Run("Create todo", func(t *testing.T) {
			payload := domain.CreateTodoDTO{
				Title:    "Integration Test Todo",
				Priority: 3,
			}
			body, _ := json.Marshal(payload)

			req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
			req = testutils.WithUserContext(req, userID)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusCreated, rr.Code)

			err := json.NewDecoder(rr.Body).Decode(&createdTodo)
			require.NoError(t, err)
			require.NotZero(t, createdTodo.ID)
			require.Equal(t, "Integration Test Todo", createdTodo.Title)
			require.Equal(t, int64(3), createdTodo.Priority)
			require.False(t, createdTodo.Done)
		})

		// 3. Get the created todo
		t.Run("Get todo by ID", func(t *testing.T) {
			url := fmt.Sprintf("/todos/%d", createdTodo.ID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var fetchedTodo domain.TodoDTO
			err := json.NewDecoder(rr.Body).Decode(&fetchedTodo)
			require.NoError(t, err)
			require.Equal(t, createdTodo.ID, fetchedTodo.ID)
			require.Equal(t, createdTodo.Title, fetchedTodo.Title)
		})

		// 4. Update the todo
		t.Run("Update todo", func(t *testing.T) {
			payload := domain.UpdateTodoDTO{
				Title:    "Updated Integration Test",
				Done:     true,
				Priority: 5,
			}
			body, _ := json.Marshal(payload)

			url := fmt.Sprintf("/todos/%d", createdTodo.ID)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			req = testutils.WithUserContext(req, userID)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var updatedTodo domain.TodoDTO
			err := json.NewDecoder(rr.Body).Decode(&updatedTodo)
			require.NoError(t, err)
			require.Equal(t, "Updated Integration Test", updatedTodo.Title)
			require.True(t, updatedTodo.Done)
			require.Equal(t, int64(5), updatedTodo.Priority)
		})

		// 5. List todos (should have one)
		t.Run("List todos after create", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var todos []domain.TodoDTO
			err := json.NewDecoder(rr.Body).Decode(&todos)
			require.NoError(t, err)
			require.Len(t, todos, 1)
			require.Equal(t, "Updated Integration Test", todos[0].Title)
		})

		// 6. Delete the todo
		t.Run("Delete todo", func(t *testing.T) {
			url := fmt.Sprintf("/todos/%d", createdTodo.ID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusNoContent, rr.Code)
		})

		// 7. Verify deletion
		t.Run("Get deleted todo returns 404", func(t *testing.T) {
			url := fmt.Sprintf("/todos/%d", createdTodo.ID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusNotFound, rr.Code)
		})
	})

	t.Run("User Isolation", func(t *testing.T) {
		router, tc, user1ID := setupTestServer(t)
		defer testutils.CleanupDB(t, tc.DB)

		ctx := context.Background()

		// Create second user
		userStore := pguser.CreateStore(tc.DB)
		userSvc := userservice.NewUserService(userStore)
		user2, err := userSvc.CreateUser(ctx, "Test User 2", "test2@example.com", "password123")
		require.NoError(t, err)
		user2ID := user2.ID

		// User 1 creates a todo
		payload := domain.CreateTodoDTO{Title: "User 1 Todo", Priority: 3}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
		req = testutils.WithUserContext(req, user1ID)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusCreated, rr.Code)

		var user1Todo domain.TodoDTO
		json.NewDecoder(rr.Body).Decode(&user1Todo)

		// User 2 tries to access User 1's todo - should fail
		t.Run("User 2 cannot access User 1 todo", func(t *testing.T) {
			url := fmt.Sprintf("/todos/%d", user1Todo.ID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = testutils.WithUserContext(req, user2ID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusNotFound, rr.Code)
		})

		// User 2 lists todos - should be empty
		t.Run("User 2 sees only their todos", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			req = testutils.WithUserContext(req, user2ID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var todos []domain.TodoDTO
			json.NewDecoder(rr.Body).Decode(&todos)
			require.Empty(t, todos)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		router, tc, userID := setupTestServer(t)
		defer testutils.CleanupDB(t, tc.DB)

		t.Run("Create with empty title", func(t *testing.T) {
			payload := domain.CreateTodoDTO{Title: "", Priority: 3}
			body, _ := json.Marshal(payload)

			req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
			req = testutils.WithUserContext(req, userID)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusBadRequest, rr.Code)
		})
	})
}
