//go:build unittest

package todolist

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/macesz/todo-go/delivery/web/todolist/mocks"
	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/tests/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestList tests the List handler with various scenarios
func TestList(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testUserID := int64(1)

	tests := []struct {
		name           string
		mockReturn     []*domain.TodoList
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - multiple lists",
			mockReturn: []*domain.TodoList{
				{
					ID:        1,
					UserID:    testUserID,
					Title:     "Shopping List",
					Color:     "#FF5733",
					Labels:    []string{"groceries", "urgent"},
					CreatedAt: fixedTime,
					Items:     []domain.Todo{},
				},
				{
					ID:        2,
					UserID:    testUserID,
					Title:     "Work Tasks",
					Color:     "#3357FF",
					Labels:    []string{"work"},
					CreatedAt: fixedTime,
					Items:     []domain.Todo{},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"ID":1,"UserID":1,"Title":"Shopping List","Color":"#FF5733","Labels":["groceries","urgent"],"CreatedAt":"2024-01-01T12:00:00Z","Items":[]},{"ID":2,"UserID":1,"Title":"Work Tasks","Color":"#3357FF","Labels":["work"],"CreatedAt":"2024-01-01T12:00:00Z","Items":[]}]`,
		},
		{
			name:           "Service error",
			mockReturn:     nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoListService(t)

			mockService.On("List", mock.Anything, testUserID).
				Return(tt.mockReturn, tt.mockError).
				Once()

			handlers := &TodoListHandlers{todoListService: mockService}

			req, err := http.NewRequest(http.MethodGet, "/lists", nil)
			require.NoError(t, err)

			// Add user context to simulate authenticated request
			req = testutils.WithUserContext(req, testUserID)

			rr := httptest.NewRecorder()
			handlers.List(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

// TestGetListByID tests the GetListByID handler with various scenarios
func TestGetListByID(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testUserID := int64(1)
	testListID := int64(1)

	tests := []struct {
		name           string
		urlParam       string
		shouldCallMock bool
		mockReturn     *domain.TodoList
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success - valid ID",
			urlParam:       "1",
			shouldCallMock: true,
			mockReturn: &domain.TodoList{
				ID:        testListID,
				UserID:    testUserID,
				Title:     "Shopping List",
				Color:     "#FF5733",
				Labels:    []string{"groceries"},
				CreatedAt: fixedTime,
				Items: []domain.Todo{
					{ID: 10, UserID: testUserID, TodoListID: testListID, Title: "Buy milk", Done: false, CreatedAt: fixedTime},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"user_id":1,"title":"Shopping List","color":"#FF5733","labels":["groceries"],"created_at":"2024-01-01T12:00:00Z","items":[{"id":10,"user_id":1,"todolist_id":1,"title":"Buy milk","done":false,"created_at":"2024-01-01T12:00:00Z"}]}`,
		},
		{
			name:           "List not found",
			urlParam:       "999",
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      domain.ErrListNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"todo list not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoListService(t)

			if tt.shouldCallMock {
				expectedID, _ := strconv.ParseInt(tt.urlParam, 10, 64)
				mockService.On("GetListByID", mock.Anything, testUserID, expectedID).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handler := &TodoListHandlers{todoListService: mockService}

			req, err := http.NewRequest(http.MethodGet, "/lists/"+tt.urlParam, nil)
			require.NoError(t, err)

			// Add user context
			req = testutils.WithUserContext(req, testUserID)

			// Add chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.GetListByID(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestCreate tests the Create handler with various scenarios
func TestCreate(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testUserID := int64(1)

	tests := []struct {
		name           string
		inputBody      string
		setupUserMock  func(*mocks.UserService)
		setupListMock  func(*mocks.TodoListService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - valid input",
			inputBody: `{"title":"Shopping List","color":"#FF5733","labels":["groceries","urgent"]}`,
			setupUserMock: func(m *mocks.UserService) {
				m.On("GetUser", mock.Anything, testUserID).
					Return(&domain.User{ID: testUserID, Name: "Test User", Email: "test@example.com"}, nil).
					Once()
			},
			setupListMock: func(m *mocks.TodoListService) {
				m.On("Create", mock.Anything, testUserID, "Shopping List", "#FF5733", []string{"groceries", "urgent"}).
					Return(&domain.TodoList{
						ID:        1,
						UserID:    testUserID,
						Title:     "Shopping List",
						Color:     "#FF5733",
						Labels:    []string{"groceries", "urgent"},
						CreatedAt: fixedTime,
					}, nil).
					Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"user_id":1,"title":"Shopping List","color":"#FF5733","labels":["groceries","urgent"],"created_at":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:      "Invalid JSON",
			inputBody: `{"title":"Broken JSON",}`, // ✅ Malformed JSON (extra comma)
			setupUserMock: func(m *mocks.UserService) {
				m.On("GetUser", mock.Anything, testUserID).
					Return(&domain.User{ID: testUserID, Name: "Test User", Email: "test@example.com"}, nil).
					Once()
			},
			setupListMock: func(m *mocks.TodoListService) {
				// ✅ Should not be called due to JSON parse error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid character '}' looking for beginning of object key string"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockUserService := mocks.NewUserService(t)
			mockListService := mocks.NewTodoListService(t)

			// Setup mocks
			tt.setupUserMock(mockUserService)
			tt.setupListMock(mockListService)

			// Create handlers with both services
			handlers := &TodoListHandlers{
				userService:     mockUserService,
				todoListService: mockListService,
			}

			// Create request
			req, err := http.NewRequest(http.MethodPost, "/lists", strings.NewReader(tt.inputBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Add user context
			req = testutils.WithUserContext(req, testUserID)

			// Create response recorder
			rr := httptest.NewRecorder()
			// Call handler
			handlers.Create(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			// Assert mock expectations
			mockUserService.AssertExpectations(t)
			mockListService.AssertExpectations(t)
		})
	}
}

// TestUpdate tests the Update handler with various scenarios
func TestUpdate(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC) // ✅ Add this
	testUserID := int64(1)

	tests := []struct {
		name           string
		urlParam       string
		inputBody      string
		shouldCallMock bool
		mockReturn     *domain.TodoList
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success - valid update",
			urlParam:       "1",
			inputBody:      `{"title":"Updated Shopping List","color":"#00FF00","labels":["groceries"]}`,
			shouldCallMock: true,
			mockReturn: &domain.TodoList{
				ID:        1,
				UserID:    testUserID,
				Title:     "Updated Shopping List",
				Color:     "#00FF00",
				Labels:    []string{"groceries"},
				CreatedAt: fixedTime,
				Deleted: false
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"user_id":1,"title":"Updated Shopping List","color":"#00FF00","labels":["groceries"],"created_at":"","deleted": false}`,
		},
		{
			name:           "List not found",
			urlParam:       "999",
			inputBody:      `{"title":"Updated List","color":"#00FF00","labels":[]}`,
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      domain.ErrListNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"todo list not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoListService(t)

			if tt.shouldCallMock {
				expectedID, _ := strconv.ParseInt(tt.urlParam, 10, 64)

				// Parse input to get expected values
				var input map[string]interface{}
				json.Unmarshal([]byte(tt.inputBody), &input)
				expectedTitle := input["title"].(string)
				expectedColor := input["color"].(string)
				expectedLabels := []string{}
				if labels, ok := input["labels"].([]interface{}); ok {
					for _, label := range labels {
						expectedLabels = append(expectedLabels, label.(string))
					}
				}

				mockService.On("Update", mock.Anything, testUserID, expectedID, expectedTitle, expectedColor, expectedLabels).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handlers := &TodoListHandlers{todoListService: mockService}

			req, err := http.NewRequest(http.MethodPut, "/lists/"+tt.urlParam, strings.NewReader(tt.inputBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Add user context
			req = testutils.WithUserContext(req, testUserID)

			// Add chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handlers.Update(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestDelete tests the Delete handler with various scenarios
func TestDelete(t *testing.T) {
	testUserID := int64(1)

	tests := []struct {
		name           string
		urlParam       string
		shouldCallMock bool
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success - valid ID",
			urlParam:       "1",
			shouldCallMock: true,
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:           "TodoList not found",
			urlParam:       "999",
			shouldCallMock: true,
			mockError:      domain.ErrListNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"todo list not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockListService := mocks.NewTodoListService(t)

			if tt.shouldCallMock {
				expectedID, _ := strconv.ParseInt(tt.urlParam, 10, 64)
				mockListService.On("Delete", mock.Anything, testUserID, expectedID).
					Return(tt.mockError).
					Once()
			}

			handlers := &TodoListHandlers{
				todoListService: mockListService,
			}

			req, err := http.NewRequest(http.MethodDelete, "/lists/"+tt.urlParam, nil)
			require.NoError(t, err)

			// Add user context
			req = testutils.WithUserContext(req, testUserID)

			// Add chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handlers.Delete(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			mockListService.AssertExpectations(t)
		})
	}
}
