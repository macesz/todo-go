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

// setupTodoListTestServer creates a real server with all dependencies for todo list testing

func Test_TodoList_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc, server, services := testutils.ComposeServer(t)

	user1 := domain.User{
		Name:     "User One",
		Email:    "u1@example.com",
		Password: "pass",
	}
	header1, err := testutils.GivenUser(t, services.TokenAuth, tc.DB, &user1)
	if err != nil {
		t.Fatal(err)
	}

	user2 := domain.User{
		Name:     "User Two",
		Email:    "u2@example.com",
		Password: "pass2",
	}

	header2, err := testutils.GivenUser(t, services.TokenAuth, tc.DB, &user2)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("Full CRUD Lifecycle", func(t *testing.T) {
		defer testutils.CleanupDB(t, tc.DB)

		// 1. List todo lists (should be empty initially)
		t.Run("List empty todo lists", func(t *testing.T) {

			resp, respbody := testutils.TestRequest(t, server, http.MethodGet, "/api/lists", header1, nil)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var lists []*domain.TodoList
			err := json.Unmarshal(respbody, &lists)

			require.NoError(t, err)
			require.Empty(t, lists, "should have no lists initially")
		})

		t.Run("GET /lists errors", func(t *testing.T) {
			t.Run("no user context -> 403", func(t *testing.T) {

				resp, _ := testutils.TestRequest(t, server, http.MethodGet, "/api/lists", nil, nil)

				require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
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

			resp, respbody := testutils.TestRequest(t, server, http.MethodPost, "/api/lists", header1, bytes.NewReader(body))

			require.Equal(t, http.StatusCreated, resp.StatusCode)

			err := json.Unmarshal(respbody, &createdList)

			require.NoError(t, err)
			require.NotZero(t, createdList.ID, "list should have an ID")
			require.Equal(t, "My Shopping List", createdList.Title)
			require.Equal(t, "#FF5733", *createdList.Color)
			require.Equal(t, []string{"shopping", "groceries"}, createdList.Labels)
		})

		t.Run("Get todoList by ID", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d", createdList.ID)

			resp, respbody := testutils.TestRequest(t, server, http.MethodGet, url, header1, nil)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var fetchedTodoList domain.TodoListDTO

			err := json.Unmarshal(respbody, &fetchedTodoList)

			require.NoError(t, err)
			require.Equal(t, createdList.ID, fetchedTodoList.ID)
			require.Equal(t, createdList.Title, fetchedTodoList.Title)
		})

		t.Run("id valid but not found -> 404", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d", int64(999999))

			resp, _ := testutils.TestRequest(t, server, http.MethodGet, url, header1, nil)

			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("no user context -> 401", func(t *testing.T) {

			url := fmt.Sprintf("/api/lists/%d", int64(1))

			resp, _ := testutils.TestRequest(t, server, http.MethodGet, url, nil, nil)

			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		// 4. Update the todo
		t.Run("Update todo list", func(t *testing.T) {
			color := "#ADD8E6"
			payload := domain.UpdateTodoListRequestDTO{
				Title:   "Updated Test",
				Color:   &color,
				Labels:  []string{"shopping"},
				Deleted: false,
			}
			body, _ := json.Marshal(payload)

			url := fmt.Sprintf("/api/lists/%d", createdList.ID)
			resp, respbody := testutils.TestRequest(t, server, http.MethodPut, url, header1, bytes.NewReader(body))

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var updatedTodoList domain.TodoListDTO

			err := json.Unmarshal(respbody, &updatedTodoList)

			require.NoError(t, err)
			require.Equal(t, "Updated Test", updatedTodoList.Title)
			require.Equal(t, "#ADD8E6", *updatedTodoList.Color)
			require.Equal(t, []string{"shopping"}, updatedTodoList.Labels)
		})

		// 6. Delete the todo
		t.Run("Delete todo list", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d", createdList.ID)

			resp, _ := testutils.TestRequest(t, server, http.MethodDelete, url, header1, nil)

			require.Equal(t, http.StatusNoContent, resp.StatusCode)
		})

		//User Isoltaion
		t.Run("User Isolation", func(t *testing.T) {

			createdID, err := testutils.GivenTodoLists(t, tc.DB, domain.TodoList{
				UserID: user1.ID,
				Title:  "My title",
			})
			require.Nil(t, err)

			// User 2 tries to access User 1's todo - should fail
			t.Run("User 2 cannot access User 1 todo", func(t *testing.T) {
				url := fmt.Sprintf("/api/lists/%d", createdID)

				resp, _ := testutils.TestRequest(t, server, http.MethodGet, url, header2, nil)

				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			// User 2 lists todos - should be empty
			t.Run("User 2 sees only their todos", func(t *testing.T) {

				resp, respbody := testutils.TestRequest(t, server, http.MethodGet, "/api/lists", header2, nil)

				require.Equal(t, http.StatusOK, resp.StatusCode)

				var todolists []domain.TodoListDTO
				json.Unmarshal(respbody, &todolists)

				require.Empty(t, todolists)
			})

		})
	})
}
