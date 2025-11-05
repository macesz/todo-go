package todo

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/macesz/todo-go/delivery/web/mocks"
	"github.com/macesz/todo-go/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockTodoService is a mock implementation of the TodoService interface for testing
type MockTodoService struct {
	mu     sync.Mutex
	todos  map[int]domain.Todo
	nextID int
}

// TestListTodos tests the ListTodos handler with various scenarios
func TestListTodos(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

	// Define test cases for different scenarios
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
				{ID: 1, Title: "Test Todo 1", Done: false, CreatedAt: fixedTime},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"ID":1,"Title":"Test Todo 1","Done":false,"CreatedAt":"2024-01-01T12:00:00Z"}]` + "\n",
		},
		{
			name:           "Service error",
			mockReturn:     nil,
			mockError:      http.ErrServerClosed,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`, // Generic message for security
		},
	}

	// Run each test case in a subtest to isolate them and allow parallel execution

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			mockService.On("ListTodos", mock.Anything).Return(tt.mockReturn, tt.mockError)

			handlers := &TodoHandlers{Service: mockService}

			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/todos", nil)
			require.NoError(t, err)

			handlers.ListTodos(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockService.AssertExpectations(t)
			rr.Body.Reset() // Reset the response recorder for the next iteration
		})
	}
}

// TestCreateTodo tests the CreateTodo handler with various scenarios
func TestCreateTodo(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		inputTitle     string // For mock matching
		inputBody      string // Request body JSON
		shouldCallMock bool   // Whether to expect service call
		mockReturn     *domain.Todo
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid input",
			inputTitle:     "New Todo",
			shouldCallMock: true,
			inputBody:      `{"title": "New Todo"}`,
			mockReturn:     &domain.Todo{ID: 1, Title: "New Todo", Done: false, CreatedAt: fixedTime},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"title":"New Todo","done":false,"createdAt":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:           "Invalid JSON",
			inputBody:      `{"Title": "New Todo"`, // Malformed JSON
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"unexpected EOF"}`, // domain.ErrorResponse format
		},
		{
			name:           "Missing title (validation error)",
			inputBody:      `{"title":""}`,
			shouldCallMock: false, // Handler validates before service call
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"title is required"}`,
		},
		{
			name:           "Service error",
			inputTitle:     "New Todo",
			inputBody:      `{"title":"New Todo"}`,
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      errors.New("error creating todo"),
			expectedStatus: http.StatusInternalServerError,      // Matches handler
			expectedBody:   `{"error":"internal server error"}`, // Generic message for security
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			if tt.shouldCallMock {
				// Flexible title matching (handles trimming)
				mockService.On("CreateTodo", mock.Anything, mock.MatchedBy(func(title string) bool {
					return strings.TrimSpace(title) == tt.inputTitle
				})).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handlers := &TodoHandlers{Service: mockService}

			rr := httptest.NewRecorder()

			rbody := strings.NewReader(tt.inputBody)
			req, err := http.NewRequest(http.MethodPost, "/todos", rbody)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			handlers.CreateTodo(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			// Conditional assertion: Use Equal for empty (non-JSON) bodies; JSONEq otherwise
			if tt.expectedBody == "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestGetTodo tests the GetTodo handler with various scenarios
func TestGetTodo(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

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
			mockReturn:     &domain.Todo{ID: 1, Title: "Test Todo", Done: false, CreatedAt: fixedTime},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Test Todo","done":false,"createdAt":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:           "Non-integer ID",
			urlParam:       "abc",
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id must be an integer"}`,
		},
		{
			name:           "Missing ID",
			urlParam:       "",
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id is required"}`,
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
		{
			name:           "Internal server error", // NEW: Covers non-NotFound errors
			urlParam:       "1",
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      errors.New("database connection failed"), // Any non-domain.ErrNotFound error
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`, // Generic message
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			if tt.shouldCallMock {
				expectedID, err := strconv.ParseInt(tt.urlParam, 10, 64)
				if err != nil {
					t.Fatalf("invalid test setup: urlParam %q is not a valid int64: %v", tt.urlParam, err)
				}
				mockService.On("GetTodo", mock.Anything, expectedID).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handler := &TodoHandlers{Service: mockService}

			rr := httptest.NewRecorder()

			reqURL := "/todos/" + tt.urlParam

			req, err := http.NewRequest(http.MethodGet, reqURL, nil)
			require.NoError(t, err)

			// Manual chi context injection
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler.GetTodo(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody == "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}
			mockService.AssertExpectations(t)
			rr.Body.Reset() // Reset the response recorder for the next iteration
		})
	}
}

// TestUpdateTodo tests the UpdateTodo handler with various scenarios
func TestUpdateTodo(t *testing.T) {

	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		urlParam       string
		inputTitle     string // For mock matching
		inputDone      bool   // For mock matching
		inputBody      string // Request body JSON
		shouldCallMock bool
		mockReturn     *domain.Todo
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid input",
			urlParam:       "1",
			inputTitle:     "Updated Todo",
			inputDone:      true,
			inputBody:      `{"title":"Updated Todo","done":true}`,
			shouldCallMock: true,
			mockReturn:     &domain.Todo{ID: 1, Title: "Updated Todo", Done: true, CreatedAt: fixedTime},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Updated Todo","done":true,"createdAt":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:           "Invalid JSON",
			urlParam:       "1",
			inputBody:      `{"title": "Updated Todo", "done": true`, // Missing closing brace
			shouldCallMock: false,
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"unexpected EOF"}`, // Matches actual json.Decode error; change to custom if you update handler
		},
		{
			name:           "Non-integer ID",
			urlParam:       "abc",
			inputBody:      `{"title": "Updated Todo", "done": true}`,
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id must be an integer"}`,
		},
		{
			name:           "Missing ID",
			urlParam:       "",
			inputBody:      `{"title": "Updated Todo", "done": true}`,
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id is required"}`,
		},
		{
			name:           "Invalid data (empty title)",
			urlParam:       "1",
			inputBody:      `{"title": "", "done": true}`, // Should fail validation
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Key: 'UpdateTodoDTO.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`,
		},
		{
			name:           "Invalid data (title too long)", // New case to test max length
			urlParam:       "1",
			inputBody:      `{"title": "` + strings.Repeat("a", 256) + `", "done": true}`, // 256 chars > 255
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Key: 'UpdateTodoDTO.Title' Error:Field validation for 'Title' failed on the 'max' tag"}`,
		},
		{
			name:           "Service error (not found)",
			urlParam:       "1",
			inputTitle:     "Updated Todo",
			inputDone:      true,
			inputBody:      `{"title":"Updated Todo","done":true}`,
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      domain.ErrNotFound, // Use custom error to match errors.Is
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"todo not found"}`, // Matches domain.ErrNotFound.Error()
		},
		{
			name:           "Service error (internal)",
			urlParam:       "1",
			inputTitle:     "Updated Todo",
			inputDone:      true,
			inputBody:      `{"title":"Updated Todo","done":true}`,
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      errors.New("database failure"), // Not a custom error, so falls to 500
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			if tt.shouldCallMock {
				expectedID, err := strconv.ParseInt(tt.urlParam, 10, 64)
				if err != nil {
					t.Fatalf("invalid test setup: urlParam %q is not a valid int64: %v", tt.urlParam, err)
				}
				mockService.On("UpdateTodo", mock.Anything, expectedID, mock.MatchedBy(func(title string) bool {
					return strings.TrimSpace(title) == tt.inputTitle
				}), tt.inputDone).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handlers := &TodoHandlers{Service: mockService}

			rr := httptest.NewRecorder()

			reqURL := "/todos/" + tt.urlParam
			if tt.urlParam == "" {
				reqURL = "/todos/"
			}

			rbody := strings.NewReader(tt.inputBody)
			req, err := http.NewRequest(http.MethodPut, reqURL, rbody)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Manual chi context injection
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handlers.UpdateTodo(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			// Conditional assertion: Use Equal for empty/non-JSON bodies; JSONEq otherwise
			if tt.expectedBody == "" {
				assert.Equal(t, tt.expectedBody, strings.TrimSpace(rr.Body.String()))
			} else {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestDeleteTodo tests the DeleteTodo handler with various scenarios

func TestDeleteTodo(t *testing.T) {
	tests := []struct {
		name           string
		urlParam       string
		shouldCallMock bool // Whether to expect service call
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
			name:           "Non-integer ID",
			urlParam:       "abc",
			shouldCallMock: false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id must be an integer"}`,
		},
		{
			name:           "Missing ID",
			urlParam:       "",
			shouldCallMock: false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id is required"}`,
		},
		{
			name:           "Todo not found",
			urlParam:       "999",
			shouldCallMock: true,
			mockError:      errors.New("not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			// Only set up mock expectations for cases that should call the service
			if tt.shouldCallMock {
				expectedID, err := strconv.ParseInt(tt.urlParam, 10, 64)
				if err != nil {
					t.Fatalf("invalid test setup: urlParam %q is not a valid int64: %v", tt.urlParam, err)
				}
				mockService.On("DeleteTodo", mock.Anything, expectedID).
					Return(tt.mockError).
					Once()
			}

			handlers := &TodoHandlers{Service: mockService}

			rr := httptest.NewRecorder()
			reqURL := "/todos/" + tt.urlParam

			// URL construction
			req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
			require.NoError(t, err)

			// Use chi to set URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handlers.DeleteTodo(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody == "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
			rr.Body.Reset()
		})
	}
}
