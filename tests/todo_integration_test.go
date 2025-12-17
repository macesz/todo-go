package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/tests/testutils"
	"github.com/stretchr/testify/require"
)

// setupTestServer creates a real server with all dependencies

func Test_Todo_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc, server, services := testutils.ComposeServer(t)

	user := &domain.User{
		Name:     "User One",
		Email:    "u1@example.com",
		Password: "pass",
	}

	header, err := testutils.GivenUser(t, services.TokenAuth, tc.DB, user)
	if err != nil {
		t.Fatal(err)
	}

	user2 := &domain.User{
		Name:     "User Two",
		Email:    "u2@example.com",
		Password: "pass2",
	}

	header2, err := testutils.GivenUser(t, services.TokenAuth, tc.DB, user2)
	if err != nil {
		t.Fatal(err)
	}

	listID, err := testutils.GivenTodoLists(t, tc.DB, domain.TodoList{
		UserID: user.ID,
		Title:  "Integration List",
	})
	require.NoError(t, err)

	listID2, err := testutils.GivenTodoLists(t, tc.DB, domain.TodoList{
		UserID: user.ID,
		Title:  "Integration List2",
	})
	require.NoError(t, err)

	todoID, err := testutils.GivenTodo(t, tc.DB, domain.Todo{UserID: user.ID, TodoListID: listID2, Title: "Todo2", Done: false})
	require.NoError(t, err)

	t.Run("Full CRUD Lifecycle", func(t *testing.T) {

		// 1. List todos (should be empty)
		t.Run("List empty todos", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d/todos", listID)
			resp, respbody := testutils.TestRequest(t, server, http.MethodGet, url, header, nil)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var todos []domain.TodoDTO
			err := json.Unmarshal(respbody, &todos)
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

			url := fmt.Sprintf("/api/lists/%d/todos", listID)
			resp, respbody := testutils.TestRequest(t, server, http.MethodPost, url, header, bytes.NewReader(body))

			require.Equal(t, http.StatusCreated, resp.StatusCode)

			err := json.Unmarshal(respbody, &createdTodo)
			require.NoError(t, err)
			require.NotZero(t, createdTodo.ID)
			require.Equal(t, "Integration Test Todo", createdTodo.Title)
			require.Equal(t, int64(3), createdTodo.Priority)
			require.False(t, createdTodo.Done)
		})

		// 3. Get the created todo
		t.Run("Get todo by ID", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d/todos/%d", listID, createdTodo.ID)
			resp, respbody := testutils.TestRequest(t, server, http.MethodGet, url, header, nil)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var fetchedTodo domain.TodoDTO

			err := json.Unmarshal(respbody, &fetchedTodo)

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

			url := fmt.Sprintf("/api/lists/%d/todos/%d", listID, createdTodo.ID)
			resp, respbody := testutils.TestRequest(t, server, http.MethodPut, url, header, bytes.NewReader(body))

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var updatedTodo domain.TodoDTO
			err := json.Unmarshal(respbody, &updatedTodo)

			require.NoError(t, err)
			require.Equal(t, "Updated Integration Test", updatedTodo.Title)
			require.True(t, updatedTodo.Done)
			require.Equal(t, int64(5), updatedTodo.Priority)
		})

		// 5. List todos (should have one)
		t.Run("List todos after create", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d/todos", listID)
			resp, respbody := testutils.TestRequest(t, server, http.MethodGet, url, header, nil)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var todos []domain.TodoDTO

			err := json.Unmarshal(respbody, &todos)
			require.NoError(t, err)
			require.Len(t, todos, 1)
			require.Equal(t, "Updated Integration Test", todos[0].Title)
		})

		// 6. Delete the todo
		t.Run("Delete todo", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d/todos/%d", listID, todoID)
			resp, _ := testutils.TestRequest(t, server, http.MethodDelete, url, header, nil)

			require.Equal(t, http.StatusNoContent, resp.StatusCode)
		})

		// 7. Verify deletion
		t.Run("Get deleted todo returns 404", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d/todos/%d", listID, todoID)
			resp, _ := testutils.TestRequest(t, server, http.MethodDelete, url, header, nil)

			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("User Isolation", func(t *testing.T) {

		// User 2 tries to access User 1's todo - should fail
		t.Run("User 2 cannot access User 1 todo", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d/todos/%d", listID, todoID)
			resp, _ := testutils.TestRequest(t, server, http.MethodGet, url, header2, nil)

			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		// User 2 lists todos - should be empty
		t.Run("User 2 sees only their todos", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d/todos", listID2)
			resp, respbody := testutils.TestRequest(t, server, http.MethodGet, url, header2, nil)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var todos []domain.TodoDTO
			json.Unmarshal(respbody, &todos)
			require.Empty(t, todos)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {

		t.Run("Create with empty title", func(t *testing.T) {
			payload := domain.CreateTodoDTO{Title: "", Priority: 3}
			body, _ := json.Marshal(payload)

			url := fmt.Sprintf("/api/lists/%d/todos", listID)
			resp, _ := testutils.TestRequest(t, server, http.MethodPost, url, header, bytes.NewReader(body))

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}
