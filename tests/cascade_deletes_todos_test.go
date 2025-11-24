package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/tests/testutils"
	"github.com/stretchr/testify/require"
)

// // setupTodoListTestServer creates a real server with all dependencies for todo list testing
// func setupTodoListTestServer(t *testing.T) (*chi.Mux, *testutils.TestContainer, int64) {
// 	t.Helper()

// 	// Setup database
// 	tc := testutils.SetupTestDB(t)

// 	// Create stores
// 	todoListStore := pgtodolist.CreateStore(tc.DB)
// 	userStore := pguser.CreateStore(tc.DB)

// 	// Create services using constructors
// 	todoListSvc := todolistservice.NewTodoListService(todoListStore)
// 	userSvc := userservice.NewUserService(userStore)

// 	// Create test user
// 	testUser, err := userSvc.CreateUser(t.Context(), "Test User", "test@example.com", "password123")
// 	require.NoError(t, err)

// 	// Create handlers using constructor
// 	todoListHandlers := todolist.NewHandlers(todoListSvc, userSvc)

// 	// Setup router with proper routes
// 	r := chi.NewRouter()
// 	r.Get("/lists", todoListHandlers.List)
// 	r.Post("/lists", todoListHandlers.Create)
// 	r.Get("/lists/{id}", todoListHandlers.GetListByID)
// 	r.Put("/lists/{id}", todoListHandlers.Update)
// 	r.Delete("/lists/{id}", todoListHandlers.Delete)

// 	return r, tc, testUser.ID
// }
//
//

func TestTodoListHandlers_CascadeDeleteDeletesTodos(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	router, tc, userID := setupTodoListTestServer(t)
	defer testutils.CleanupDB(t, tc.DB)

	// 1. Create a list for this user directly in DB
	listID, err := testutils.GivenTodoLists(t, tc.DB, domain.TodoList{
		UserID:    userID,
		Title:     "List with todos",
		Color:     "#FFFFFF",
		Labels:    []string{"test"},
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)

	// 2. Create one or more todos in that list
	_, err = testutils.GivenTodo(t, tc.DB, domain.Todo{
		UserID:    userID,
		ListID:    listID,
		Title:     "Todo 1",
		Done:      false,
		Priority:  1,
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)

	_, err = testutils.GivenTodo(t, tc.DB, domain.Todo{
		UserID:    userID,
		ListID:    listID,
		Title:     "Todo 2",
		Done:      false,
		Priority:  2,
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)

	// Sanity check: there are todos for this list
	var beforeCount int
	err = tc.DB.Get(&beforeCount, "SELECT COUNT(*) FROM todos WHERE list_id = $1", listID)
	require.NoError(t, err)
	require.Equal(t, 2, beforeCount)

	// 3. Delete the list via HTTP
	url := fmt.Sprintf("/lists/%d", listID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req = testutils.WithUserContext(req, userID)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)

	// 4. Assert that todos were deleted (cascade delete)
	var afterCount int
	err = tc.DB.Get(&afterCount, "SELECT COUNT(*) FROM todos WHERE list_id = $1", listID)
	require.NoError(t, err)
	require.Equal(t, 0, afterCount, "todos should be deleted when list is deleted")
}
