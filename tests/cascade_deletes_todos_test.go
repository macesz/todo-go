package tests

import (
	"fmt"
	"net/http"
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

	tc, server, services := testutils.ComposeServer(t)

	user := domain.User{
		Name:     "User One",
		Email:    "u1@example.com",
		Password: "pass",
	}
	header, err := testutils.GivenUser(t, services.TokenAuth, tc.DB, &user)

	if err != nil {
		t.Fatal(err)
	}

	// 1. Create a list for this user directly in DB
	todolistID, _ := testutils.GivenTodoLists(t, tc.DB, domain.TodoList{UserID: user.ID, Title: "User 1 Secrets"})

	require.NoError(t, err)

	// 2. Create one or more todos in that list
	_, err = testutils.GivenTodo(t, tc.DB, domain.Todo{
		UserID:     user.ID,
		TodoListID: todolistID,
		Title:      "Todo 1",
		Done:       false,
		CreatedAt:  time.Now(),
	})
	require.NoError(t, err)

	_, err = testutils.GivenTodo(t, tc.DB, domain.Todo{
		UserID:     user.ID,
		TodoListID: todolistID,
		Title:      "Todo 2",
		Done:       false,
		CreatedAt:  time.Now(),
	})
	require.NoError(t, err)

	// Sanity check: there are todos for this list
	var beforeCount int
	err = tc.DB.Get(&beforeCount, "SELECT COUNT(*) FROM todos WHERE todolist_id = $1", todolistID)
	require.NoError(t, err)
	require.Equal(t, 2, beforeCount)

	// 3. Delete the list via HTTP
	url := fmt.Sprintf("/api/lists/%d", todolistID)
	resp, _ := testutils.TestRequest(t, server, http.MethodDelete, url, header, nil)

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// 4. Assert that todos were deleted (cascade delete)
	var afterCount int
	err = tc.DB.Get(&afterCount, "SELECT COUNT(*) FROM todos WHERE todolist_id = $1", todolistID)
	require.NoError(t, err)
	require.Equal(t, 0, afterCount, "todos should be deleted when list is deleted")
}
