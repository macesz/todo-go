package todo

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

	"github.com/macesz/todo-go/tests/testutils"

	chi "github.com/go-chi/chi/v5"
	"github.com/macesz/todo-go/delivery/web/todo/mocks"
	"github.com/macesz/todo-go/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestListTodos tests the ListTodos handler with various scenarios
func TestListTodos(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testUserID := int64(1)
	testListID := int64(1)

	tests := []struct {
		name           string
		mockReturn     []*domain.Todo
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "No todos",
			mockReturn:     []*domain.Todo{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "[]",
		},
		{
			name: "One todo",
			mockReturn: []*domain.Todo{
				{ID: 1, UserID: testUserID, ListID: testListID, Title: "Test Todo 1", Done: false, Priority: 3, CreatedAt: fixedTime},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"ID":1,"UserID": 1, "ListID": 1, "Title":"Test Todo 1","Done":false,"Priority": 3,"CreatedAt":"2024-01-01T12:00:00Z"}]`,
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
			mockService := mocks.NewTodoService(t)

			// Updated to match new signature with ListID
			mockService.On("ListTodos", mock.Anything, testUserID, testListID).
				Return(tt.mockReturn, tt.mockError).
				Once()

			handlers := &TodoHandlers{todoService: mockService}

			req, err := http.NewRequest(http.MethodGet, "/{listID}/todos/", nil)
			require.NoError(t, err)

			// Add user context to simulate authenticated request
			req = testutils.WithUserContext(req, testUserID)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("listID", "1") // Add the listID parameter
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handlers.ListTodos(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestCreateTodo(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testUserID := int64(1)
	testListID := int64(1)

	tests := []struct {
		name           string
		inputBody      string
		setupUserMock  func(*mocks.UserService)
		setupTodoMock  func(*mocks.TodoService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Valid input",
			inputBody: `{"title": "New Todo", "priority": 2}`,
			setupUserMock: func(m *mocks.UserService) {
				m.On("GetUser", mock.Anything, testUserID).
					Return(&domain.User{ID: testUserID, Name: "Test User", Email: "test@example.com"}, nil).
					Once()
			},
			setupTodoMock: func(m *mocks.TodoService) {
				m.On("CreateTodo", mock.Anything, testUserID, testListID, "New Todo", int64(2)).
					Return(&domain.Todo{
						ID:        1,
						UserID:    testUserID,
						ListID:    testListID,
						Title:     "New Todo",
						Done:      false,
						Priority:  2,
						CreatedAt: fixedTime,
					}, nil).
					Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"user_id":1,"list_id":1,"title":"New Todo","done":false,"priority":2,"created_at":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:      "Missing title",
			inputBody: `{"title":"", "priority": 2}`,
			setupUserMock: func(m *mocks.UserService) {
				m.On("GetUser", mock.Anything, testUserID).
					Return(&domain.User{ID: testUserID, Name: "Test User", Email: "test@example.com"}, nil).
					Once()
			},
			setupTodoMock: func(m *mocks.TodoService) {
				// Should not be called due to validation error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"title is required"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockUserService := mocks.NewUserService(t)
			mockTodoService := mocks.NewTodoService(t)

			// Setup mocks
			tt.setupUserMock(mockUserService)
			tt.setupTodoMock(mockTodoService)

			// Create handlers with both services
			handlers := &TodoHandlers{
				userService: mockUserService,
				todoService: mockTodoService,
			}

			// Create request
			req, err := http.NewRequest(http.MethodPost, "/{listID}/todos", strings.NewReader(tt.inputBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Add user context
			req = testutils.WithUserContext(req, testUserID)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("listID", "1") // Add the listID parameter
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			rr := httptest.NewRecorder()
			// Call handler
			handlers.CreateTodo(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			// Assert mock expectations
			mockUserService.AssertExpectations(t)
			mockTodoService.AssertExpectations(t)
		})
	}
}

// TestGetTodo tests the GetTodo handler with various scenarios
func TestGetTodo(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testUserID := int64(1)

	tests := []struct {
		name           string
		urlParam       string
		shouldCallMock bool
		mockReturn     *domain.Todo
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid ID",
			urlParam:       "1",
			shouldCallMock: true,
			mockReturn:     &domain.Todo{ID: 1, UserID: testUserID, Title: "Test Todo", Done: false, Priority: 3, CreatedAt: fixedTime},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"userID":1,"title":"Test Todo","done":false,"priority":3,"createdAt":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:           "Todo not found",
			urlParam:       "999",
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"todo not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			if tt.shouldCallMock {
				expectedID, _ := strconv.ParseInt(tt.urlParam, 10, 64)
				// Updated to match new signature: GetTodo(ctx, userID, todoID)
				mockService.On("GetTodo", mock.Anything, testUserID, expectedID).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handler := &TodoHandlers{todoService: mockService}

			req, err := http.NewRequest(http.MethodGet, "/todos/"+tt.urlParam, nil)
			require.NoError(t, err)

			// Add user context
			req = testutils.WithUserContext(req, testUserID)

			// Add chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.GetTodo(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestUpdateTodo tests the UpdateTodo handler with various scenarios
func TestUpdateTodo(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testUserID := int64(1)

	tests := []struct {
		name           string
		urlParam       string
		inputBody      string
		shouldCallMock bool
		mockReturn     *domain.Todo
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid input",
			urlParam:       "1",
			inputBody:      `{"title":"Updated Todo","done":true,"priority":1}`,
			shouldCallMock: true,
			mockReturn:     &domain.Todo{ID: 1, UserID: testUserID, Title: "Updated Todo", Done: true, Priority: 1, CreatedAt: fixedTime},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"userID":1,"title":"Updated Todo","done":true,"priority":1,"createdAt":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:           "Todo not found",
			urlParam:       "1",
			inputBody:      `{"title":"Updated Todo","done":true,"priority":1}`,
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"todo not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			if tt.shouldCallMock {
				expectedID, _ := strconv.ParseInt(tt.urlParam, 10, 64)

				// Parse input to get expected values
				var input map[string]interface{}
				json.Unmarshal([]byte(tt.inputBody), &input)
				expectedTitle := input["title"].(string)
				expectedDone := input["done"].(bool)
				expectedPriority := int64(1)
				if p, ok := input["priority"].(float64); ok {
					expectedPriority = int64(p)
				}

				// Updated to match new signature: UpdateTodo(ctx, userID, todoID, title, done, priority)
				mockService.On("UpdateTodo", mock.Anything, testUserID, expectedID, expectedTitle, expectedDone, expectedPriority).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handlers := &TodoHandlers{todoService: mockService}

			req, err := http.NewRequest(http.MethodPut, "/todos/"+tt.urlParam, strings.NewReader(tt.inputBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Add user context
			req = testutils.WithUserContext(req, testUserID)

			// Add chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handlers.UpdateTodo(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestDeleteTodo tests the DeleteTodo handler with various scenarios
func TestDeleteTodo(t *testing.T) {
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
			name:           "Valid ID",
			urlParam:       "1",
			shouldCallMock: true,
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:           "Todo not found",
			urlParam:       "999",
			shouldCallMock: true,
			mockError:      domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"todo not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			if tt.shouldCallMock {
				expectedID, _ := strconv.ParseInt(tt.urlParam, 10, 64)
				// Updated to match new signature: DeleteTodo(ctx, userID, todoID)
				mockService.On("DeleteTodo", mock.Anything, testUserID, expectedID).
					Return(tt.mockError).
					Once()
			}

			handlers := &TodoHandlers{todoService: mockService}

			req, err := http.NewRequest(http.MethodDelete, "/todos/"+tt.urlParam, nil)
			require.NoError(t, err)

			// Add user context
			req = testutils.WithUserContext(req, testUserID)

			// Add chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handlers.DeleteTodo(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}
