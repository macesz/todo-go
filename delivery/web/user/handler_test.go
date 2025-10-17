package user

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/macesz/todo-go/delivery/web/user/mocks"
	"github.com/macesz/todo-go/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTodoService struct {
	mu     sync.Mutex
	todos  map[int]domain.User
	nextID int
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		inputName      string // For mock matching
		inputEmail     string // For mock matching
		inputPassword  string // For mock matching
		shouldCallMock bool   // Whether to expect service call
		inputBody      string
		mockReturn     *domain.User
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid input",
			inputName:      "Test User",
			inputEmail:     "test@example.com",
			inputPassword:  "password",
			shouldCallMock: true,
			inputBody:      `{"name":"Test User","email":"test@example.com","password":"password"}`,
			mockReturn:     &domain.User{ID: 1, Name: "Test User", Email: "test@example.com"},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"email":"test@example.com", "id":1, "name":"Test User"}`,
		}, {
			name:           "Invalid JSON",
			inputBody:      `{"name":"Test User"`, // Malformed (missing closing brace)
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"unexpected EOF"}`, // Match actual decoder error (run handler to confirm exact string)
		},
		{
			name:           "Internal server error",
			inputName:      "Test User",
			inputEmail:     "test@example.com",
			inputPassword:  "password",
			inputBody:      `{"name":"Test User","email":"test@example.com","password":"password"}`,
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      errors.New("database failure"), // Generic error → 500
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
		{
			name:           "Missing Name",
			inputBody:      `{"email":"test@example.com","password":"password"}`, // Valid JSON, missing name
			shouldCallMock: false,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Name is required"}`,
		}, {
			name:           "Missing Email",
			inputBody:      `{"name":"Test User","password":"password"}`, // Valid JSON, missing email
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Email is required"}`,
		}, {
			name:           "Missing Password",
			inputBody:      `{"name":"Test User","email":"test@example.com"}`, // Valid JSON, missing password
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Password is required"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewUserService(t)

			if tt.shouldCallMock {
				mockService.On("CreateUser", mock.Anything, tt.inputName, tt.inputEmail, tt.inputPassword).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handlers := &UserHandlers{
				Service: mockService,
			}

			rr := httptest.NewRecorder()

			rbody := strings.NewReader(tt.inputBody)
			req, err := http.NewRequest("POST", "/users", rbody)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json") // Required for JSON decoding in handler

			handlers.CreateUser(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String()) // Handles JSON whitespace/order

			mockService.AssertExpectations(t)

			rr.Body.Reset()
		})
	}

}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name           string
		urlParam       string
		shouldCallMock bool // Whether to expect service call
		mockReturn     *domain.User
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid id",
			urlParam:       "1",
			shouldCallMock: true,
			mockReturn:     &domain.User{ID: 1, Name: "Test User", Email: "test@example.com", Password: "hashedpassword123"}, // Password is included here for the mock but assumed omitted in response
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"Test User","email":"test@example.com"}`,
		}, {
			name:           "Non-integer ID",
			urlParam:       "abc",
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id must be an integer"}`,
		}, {
			name:           "User not found",
			urlParam:       "999",
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      domain.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"user not found"}`,
		}, {
			name:           "Missing ID",
			urlParam:       "",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id is required"}`,
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
			mockService := mocks.NewUserService(t)

			if tt.shouldCallMock {
				expectedId, err := strconv.ParseInt(tt.urlParam, 10, 64)
				if err != nil {
					t.Fatalf("invalid test setup: urlParam %q is not a valid int64", tt.urlParam)
				}

				mockService.On("GetUser", mock.Anything, expectedId).
					Return(tt.mockReturn, tt.mockError).Once()
			}

			handlers := &UserHandlers{
				Service: mockService,
			}

			rr := httptest.NewRecorder()

			reqURL := "/users/" + tt.urlParam

			req, err := http.NewRequest(http.MethodGet, reqURL, nil)
			require.NoError(t, err)

			// Manual chi context injection (consistent with your DeleteUser test)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handlers.GetUser(rr, req) // Assumes your handler method is named GetUser

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

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		urlParam       string
		shouldCallMock bool // Whether to expect service call
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid id",
			urlParam:       "1",
			shouldCallMock: true,
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		}, {
			name:           "Non-integer ID",
			urlParam:       "abc",
			shouldCallMock: false,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id must be an integer"}`,
		}, {
			name:           "User not found",
			urlParam:       "999",
			shouldCallMock: true,
			mockError:      domain.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"user not found"}`,
		}, {
			name:           "Missing ID",
			urlParam:       "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"id is required"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewUserService(t)

			if tt.shouldCallMock {
				expectedId, err := strconv.ParseInt(tt.urlParam, 10, 64)
				if err != nil {
					t.Fatalf("invalid test setup: urlParam %q is not a valid int64", tt.urlParam)
				}

				mockService.On("DeleteUser", mock.Anything, expectedId).
					Return(tt.mockError).Once()
			}

			handlers := &UserHandlers{
				Service: mockService,
			}

			rr := httptest.NewRecorder()

			reqURL := "/users/" + tt.urlParam

			req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
			require.NoError(t, err) // Better than t.Fatal

			// Use chi to set URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handlers.DeleteUser(rr, req)

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
