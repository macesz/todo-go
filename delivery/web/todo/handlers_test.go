package todo

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/macesz/todo-go/delivery/web/mocks"
	"github.com/macesz/todo-go/domain"
	"github.com/stretchr/testify/mock"
)

// MockTodoService is a mock implementation of the TodoService interface for testing
type MockTodoService struct {
	mu     sync.Mutex
	todos  map[int]domain.Todo
	nextID int
}

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestListTodos tests the ListTodos handler with various scenarios
func TestListTodos(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

	// Define test cases for different scenarios
	tests := []struct {
		name           string
		mockReturn     []domain.Todo
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "No todos",
			mockReturn:     []domain.Todo{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "[]\n",
		},
		{
			name: "One todo",
			mockReturn: []domain.Todo{
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
			expectedBody:   `{"error":"http: Server closed"}` + "\n",
		},
	}

	// Run each test case in a subtest to isolate them and allow parallel execution

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			mockService.On("ListTodos", mock.Anything).Return(tt.mockReturn, tt.mockError)

			handlers := &TodoHandlers{Service: mockService}
			handler := http.HandlerFunc(handlers.ListTodos)

			rr := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/todos", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}

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
		inputBody      string
		mockReturn     domain.Todo
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid input",
			inputBody:      `{"Title": "New Todo"}`,
			mockReturn:     domain.Todo{ID: 1, Title: "New Todo", Done: false, CreatedAt: fixedTime},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"title":"New Todo","done":false,"createdAt":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:           "Invalid JSON",
			inputBody:      `{"Title": "New Todo"`, // Malformed JSON
			mockReturn:     domain.Todo{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"unexpected EOF"`,
		},
		{
			name:           "Service error",
			inputBody:      `{"Title": "New Todo"}`,
			mockReturn:     domain.Todo{},
			mockError:      errors.New("error creating todo"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error creating todo"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			if tt.mockError != nil || tt.mockReturn.ID != 0 {
				mockService.On("CreateTodo", mock.Anything, "New Todo").Return(tt.mockReturn, tt.mockError)
			}

			handlers := &TodoHandlers{Service: mockService}
			handler := http.HandlerFunc(handlers.CreateTodo)

			rr := httptest.NewRecorder()

			rbody := strings.NewReader(tt.inputBody)

			req, err := http.NewRequest("POST", "/todos", rbody)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}

			mockService.AssertExpectations(t)
			rr.Body.Reset() // Reset the response recorder for the next iteration
		})
	}
}

// TestGetTodo tests the GetTodo handler with various scenarios
func TestGetTodo(t *testing.T) {
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		url            string
		mockReturn     domain.Todo
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid ID",
			url:            "/todos/1",
			mockReturn:     domain.Todo{ID: 1, Title: "Test Todo", Done: false, CreatedAt: fixedTime},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Test Todo","done":false,"createdAt":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:           "Non-integer ID",
			url:            "/todos/abc",
			mockReturn:     domain.Todo{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id must be an integer"}`,
		},
		{
			name:           "Missing ID",
			url:            "/todos/",
			mockReturn:     domain.Todo{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id is required"}`,
		},
		{
			name:           "Todo not found",
			url:            "/todos/999",
			mockReturn:     domain.Todo{},
			mockError:      errors.New("not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			if tt.mockError != nil || tt.mockReturn.ID != 0 {
				mockService.On("GetTodo", mock.Anything, mock.AnythingOfType("int")).Return(tt.mockReturn, tt.mockError)
			}

			handler := &TodoHandlers{Service: mockService}
			// handler := http.HandlerFunc(handlers.GetTodo)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/todos/"+tt.url, nil)

			// Use chi to set URL params, (rctx is for routing context), so we can simulate URL parameters
			rctx := chi.NewRouteContext()
			// Extract the ID from the URL path and set it as a URL parameter
			rctx.URLParams.Add("id", strings.TrimPrefix(tt.url, "/todos/"))
			// Add the routing context to the request's context
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler.GetTodo(rr, req)
			// handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
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
		inputBody      string
		mockID         int    // Expected ID for mock
		mockTitle      string // Expected title for mock
		mockDone       bool   // Expected done for mock
		shouldCallMock bool   // Whether to expect service call
		mockReturn     domain.Todo
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Valid input",
			urlParam:  "1",
			inputBody: `{"title": "Updated Todo", "done": true}`,
			mockID:    1,
			mockTitle: "Updated Todo",
			mockDone:  true,
			mockReturn: domain.Todo{
				ID: 1, Title: "Updated Todo", Done: true, CreatedAt: fixedTime,
			},
			mockError:      nil,
			shouldCallMock: true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Updated Todo","done":true,"createdAt":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:           "Invalid JSON",
			urlParam:       "1",
			inputBody:      `{"title": "Updated Todo", "done": true`, // Missing closing brace
			shouldCallMock: false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"unexpected EOF"`,
		},
		{
			name:           "Non-integer ID",
			urlParam:       "abc",
			inputBody:      `{"title": "Updated Todo", "done": true}`,
			shouldCallMock: false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id must be an integer"}`,
		},
		{
			name:           "Missing ID",
			urlParam:       "",
			inputBody:      `{"title": "Updated Todo", "done": true}`,
			shouldCallMock: false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id is required"}`,
		},
		{
			name:           "Invalid data (empty title)",
			urlParam:       "1",
			inputBody:      `{"title": "", "done": true}`, // Should fail validation
			shouldCallMock: false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"title is required and must be between 1 and 255 characters; done is required"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewTodoService(t)

			// Only set up mock expectations for cases that should call the service
			if tt.shouldCallMock {
				mockService.On("UpdateTodo", mock.Anything, tt.mockID, tt.mockTitle, tt.mockDone).
					Return(tt.mockReturn, tt.mockError)
			}

			handlers := &TodoHandlers{Service: mockService}

			rr := httptest.NewRecorder()
			rbody := strings.NewReader(tt.inputBody)

			// URL construction
			req, err := http.NewRequest("PUT", "/todos/"+tt.urlParam, rbody)
			if err != nil {
				t.Fatal(err)
			}

			// Set Content-Type header (important for JSON parsing)
			req.Header.Set("Content-Type", "application/json")

			// Use chi to set URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handlers.UpdateTodo(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			gotBody := strings.TrimSpace(rr.Body.String())
			if gotBody != tt.expectedBody {
				t.Errorf("handler returned unexpected body:\ngot:  %q\nwant: %q",
					gotBody, tt.expectedBody)
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
		mockID         int  // Expected ID for mock
		shouldCallMock bool // Whether to expect service call
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid ID",
			urlParam:       "1",
			mockID:         1,
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
			mockID:         999,
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
				mockService.On("DeleteTodo", mock.Anything, tt.mockID).
					Return(tt.mockError)
			}

			handlers := &TodoHandlers{Service: mockService}

			rr := httptest.NewRecorder()

			// URL construction
			req, err := http.NewRequest("DELETE", "/todos/"+tt.urlParam, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Use chi to set URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handlers.DeleteTodo(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			gotBody := strings.TrimSpace(rr.Body.String())
			if gotBody != tt.expectedBody {
				t.Errorf("handler returned unexpected body:\ngot:  %q\nwant: %q",
					gotBody, tt.expectedBody)
			}

			mockService.AssertExpectations(t)
		})
	}
}
