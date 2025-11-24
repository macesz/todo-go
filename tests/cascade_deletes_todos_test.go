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

func Test_CascadeDeleteDeletesTodos(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	router, tc, userID := setupTodoListTestServer(t)
	defer testutils.CleanupDB(t, tc.DB)

	// 1. Create a list for this user directly in DB
	todolistID, err := testutils.GivenTodoLists(t, tc.DB, domain.TodoList{
		UserID:    userID,
		Title:     "List with todos",
		Color:     "#FFFFFF",
		Labels:    []string{"test"},
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)

	// 2. Create one or more todos in that list
	_, err = testutils.GivenTodo(t, tc.DB, domain.Todo{
		UserID:     userID,
		TodoListID: todolistID,
		Title:      "Todo 1",
		Done:       false,
		Priority:   1,
		CreatedAt:  time.Now(),
	})
	require.NoError(t, err)

	_, err = testutils.GivenTodo(t, tc.DB, domain.Todo{
		UserID:     userID,
		TodoListID: todolistID,
		Title:      "Todo 2",
		Done:       false,
		Priority:   2,
		CreatedAt:  time.Now(),
	})
	require.NoError(t, err)

	// Sanity check: there are todos for this list
	var beforeCount int
	err = tc.DB.Get(&beforeCount, "SELECT COUNT(*) FROM todos WHERE todolist_id = $1", todolistID)
	require.NoError(t, err)
	require.Equal(t, 2, beforeCount)

	// 3. Delete the list via HTTP
	url := fmt.Sprintf("/lists/%d", todolistID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req = testutils.WithUserContext(req, userID)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)

	// 4. Assert that todos were deleted (cascade delete)
	var afterCount int
	err = tc.DB.Get(&afterCount, "SELECT COUNT(*) FROM todos WHERE todolist_id = $1", todolistID)
	require.NoError(t, err)
	require.Equal(t, 0, afterCount, "todos should be deleted when list is deleted")
}
