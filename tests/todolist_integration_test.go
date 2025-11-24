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
	"github.com/macesz/todo-go/dal/pgtodolist"
	"github.com/macesz/todo-go/dal/pguser"
	"github.com/macesz/todo-go/delivery/web/todolist"
	"github.com/macesz/todo-go/domain"
	todolistservice "github.com/macesz/todo-go/services/todolist"
	userservice "github.com/macesz/todo-go/services/user"
	"github.com/macesz/todo-go/tests/testutils"
	"github.com/stretchr/testify/require"
)

// setupTodoListTestServer creates a real server with all dependencies for todo list testing
func setupTodoListTestServer(t *testing.T) (*chi.Mux, *testutils.TestContainer, int64) {
	t.Helper()

	// Setup database
	tc := testutils.SetupTestDB(t)

	// Create stores
	todoListStore := pgtodolist.CreateStore(tc.DB)
	userStore := pguser.CreateStore(tc.DB)

	// Create services using constructors
	todoListSvc := todolistservice.NewTodoListService(todoListStore)
	userSvc := userservice.NewUserService(userStore)

	// Create test user
	testUser, err := userSvc.CreateUser(t.Context(), "Test User", "test@example.com", "password123")
	require.NoError(t, err)

	// Create handlers using constructor
	todoListHandlers := todolist.NewHandlers(todoListSvc, userSvc)

	// Setup router with proper routes
	r := chi.NewRouter()
	r.Get("/lists", todoListHandlers.List)
	r.Post("/lists", todoListHandlers.Create)
	r.Get("/lists/{id}", todoListHandlers.GetListByID)
	r.Put("/lists/{id}", todoListHandlers.Update)
	r.Delete("/lists/{id}", todoListHandlers.Delete)

	return r, tc, testUser.ID
}

func Test_TodoList_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Full CRUD Lifecycle", func(t *testing.T) {
		router, tc, userID := setupTodoListTestServer(t)
		defer testutils.CleanupDB(t, tc.DB)

		// 1. List todo lists (should be empty initially)
		t.Run("List empty todo lists", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/lists", nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var lists []*domain.TodoList
			err := json.NewDecoder(rr.Body).Decode(&lists)
			require.NoError(t, err)
			require.Empty(t, lists, "should have no lists initially")
		})

		t.Run("GET /lists errors", func(t *testing.T) {
			t.Run("no user context -> 403", func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/lists", nil)
				// no WithUserContext
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				require.Equal(t, http.StatusForbidden, rr.Code)
			})
		})
		// 2. Create a todo list
		var createdList domain.TodoListDTO

		t.Run("Create todo list", func(t *testing.T) {
			color := "#FF5733"
			payload := domain.CreateTodoListRequestDTO{
				Title:  "My Shopping List",
				Color:  &color,
				Labels: []string{"shopping", "groceries"},
			}
			body, _ := json.Marshal(payload)

			req := httptest.NewRequest(http.MethodPost, "/lists", bytes.NewReader(body))
			req = testutils.WithUserContext(req, userID)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusCreated, rr.Code)

			err := json.NewDecoder(rr.Body).Decode(&createdList)
			require.NoError(t, err)
			require.NotZero(t, createdList.ID, "list should have an ID")
			require.Equal(t, "My Shopping List", createdList.Title)
			require.Equal(t, "#FF5733", *createdList.Color)
			require.Equal(t, []string{"shopping", "groceries"}, createdList.Labels)
		})

		t.Run("Get todoList by ID", func(t *testing.T) {
			url := fmt.Sprintf("/lists/%d", createdList.ID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var fetchedTodoList domain.TodoListDTO

			err := json.NewDecoder(rr.Body).Decode(&fetchedTodoList)

			require.NoError(t, err)
			require.Equal(t, createdList.ID, fetchedTodoList.ID)
			require.Equal(t, createdList.Title, fetchedTodoList.Title)
		})

		t.Run("id valid but not found -> 404", func(t *testing.T) {
			url := fmt.Sprintf("/lists/%d", int64(999999))
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusNotFound, rr.Code)
		})

		t.Run("no user context -> 403", func(t *testing.T) {
			url := fmt.Sprintf("/lists/%d", int64(1))
			req := httptest.NewRequest(http.MethodGet, url, nil) // no WithUserContext
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusForbidden, rr.Code)
		})

		t.Run("list not found -> 404", func(t *testing.T) {
			color := "#FF5733"

			payload := domain.UpdateTodoListRequestDTO{
				Title:  "Does not exist",
				Color:  &color,
				Labels: []string{"x"},
			}
			body, _ := json.Marshal(payload)

			url := fmt.Sprintf("/lists/%d", int64(999999))
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			req = testutils.WithUserContext(req, userID)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusNotFound, rr.Code)
		})

		// 4. Update the todo
		t.Run("Update todo list", func(t *testing.T) {
			color := "#ADD8E6"
			payload := domain.UpdateTodoListRequestDTO{
				Title:  "Updated Test",
				Color:  &color,
				Labels: []string{"shopping"},
			}
			body, _ := json.Marshal(payload)

			url := fmt.Sprintf("/lists/%d", createdList.ID)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			req = testutils.WithUserContext(req, userID)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var updatedTodoList domain.TodoListDTO

			err := json.NewDecoder(rr.Body).Decode(&updatedTodoList)
			require.NoError(t, err)
			require.Equal(t, "Updated Test", updatedTodoList.Title)
			require.Equal(t, "#ADD8E6", *updatedTodoList.Color)
			require.Equal(t, []string{"shopping"}, updatedTodoList.Labels)
		})

		// 6. Delete the todo
		t.Run("Delete todo list", func(t *testing.T) {
			url := fmt.Sprintf("/lists/%d", createdList.ID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusNoContent, rr.Code)
		})
		t.Run("list not found -> 404", func(t *testing.T) {
			url := fmt.Sprintf("/lists/%d", int64(999999))
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req = testutils.WithUserContext(req, userID)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusNotFound, rr.Code)
		})

		//User Isoltaion
		t.Run("User Isolation", func(t *testing.T) {
			router, tc, user1ID := setupTodoListTestServer(t)
			defer testutils.CleanupDB(t, tc.DB)

			ctx := context.Background()

			// Create second user
			userStore := pguser.CreateStore(tc.DB)
			userSvc := userservice.NewUserService(userStore)
			user2, err := userSvc.CreateUser(ctx, "Test User 2", "test2@example.com", "password123")
			require.NoError(t, err)
			user2ID := user2.ID

			createdID, err := testutils.GivenTodoLists(t, tc.DB, domain.TodoList{
				UserID: user1ID,
				Title:  "My title",
			})
			require.Nil(t, err)

			// User 2 tries to access User 1's todo - should fail
			t.Run("User 2 cannot access User 1 todo", func(t *testing.T) {
				url := fmt.Sprintf("/lists/%d", createdID)
				req := httptest.NewRequest(http.MethodGet, url, nil)
				req = testutils.WithUserContext(req, user2ID)
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				require.Equal(t, http.StatusNotFound, rr.Code)
			})

			// User 2 lists todos - should be empty
			t.Run("User 2 sees only their todos", func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/lists", nil)
				req = testutils.WithUserContext(req, user2ID)
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				require.Equal(t, http.StatusOK, rr.Code)

				var todolists []domain.TodoListDTO
				json.NewDecoder(rr.Body).Decode(&todolists)
				require.Empty(t, todolists)
			})

		})
	})
}
