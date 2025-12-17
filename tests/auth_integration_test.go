package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/tests/testutils"
	"github.com/stretchr/testify/require"
)

func Test_Auth_Integration2(t *testing.T) {
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

	// Create a list for User 1
	listID, _ := testutils.GivenTodoLists(t, tc.DB, domain.TodoList{UserID: user1.ID, Title: "User 1 Secrets"})

	t.Run("Authentication (Who are you?)", func(t *testing.T) {
		t.Run("Access protected route without token -> 401", func(t *testing.T) {
			resp, _ := testutils.TestRequest(t, server, http.MethodGet, "/api/lists", nil, nil)

			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Access protected route with invalid token -> 401", func(t *testing.T) {
			resp, _ := testutils.TestRequest(t, server, http.MethodGet, "/api/lists", testutils.AddBerrierTokenToHeader("invalid-rubbish-token", nil), nil)

			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Access protected route with valid token -> 200", func(t *testing.T) {
			resp, _ := testutils.TestRequest(t, server, http.MethodGet, "/api/lists", header1, nil)

			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("Authorization (Access Control)", func(t *testing.T) {

		t.Run("User 1 can see their own list", func(t *testing.T) {
			url := fmt.Sprintf("/api/lists/%d", listID)
			resp, _ := testutils.TestRequest(t, server, http.MethodGet, url, header1, nil)

			require.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("User 2 CANNOT see User 1's list -> 404 or 403", func(t *testing.T) {
			// User 2 has a valid token, but requests User 1's data
			url := fmt.Sprintf("/api/lists/%d", listID)
			resp, _ := testutils.TestRequest(t, server, http.MethodGet, url, header2, nil)

			// Should be Not Found (secure) or Forbidden
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

}
