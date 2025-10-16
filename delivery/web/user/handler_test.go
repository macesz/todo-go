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
			inputBody:      `{"name":"Test User", "email":"test@example.com", "password": "password"}`,
			mockReturn: &domain.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: `{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
			}`,
		}, {
			name:           "Invalid JSON",
			inputBody:      `{"Name:"Test User", Email:"test@example.com", password: "password","}`,
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid JSON"}`,
		},
		{
			name:           "Service error",
			inputName:      "Test User",
			inputEmail:     "test@example.com",
			inputPassword:  "password123",
			inputBody:      `{"Name:"Test User", Email:"test@example.com", password: "password","}`,
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      errors.New("error creating user"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"error creating user"}`,
		},
		{
			name:           "Missing Name",
			inputBody:      `{"name":"",Email:"test@example.com", password: "password","}`,
			shouldCallMock: false,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"missing required name"}`,
		}, {
			name:           "Missing Email",
			inputBody:      `{"name":"Test User",Email:"", password: "password","}`,
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"missing required email"}`,
		}, {
			name:           "Missing Password",
			inputBody:      `{"name":"Test User",Email:"test@example.com", password: "","}`,
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"missing required password"}`,
		}, {
			name:           "Wrong Email",
			inputBody:      `{"name":"Test User",Email:"test.user.com", password: "password","}`,
			shouldCallMock: false,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"wrong email"}`,
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

func GetUser(t *testing.T) {
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
			expectedBody:   `{"error":"invalid user ID"}`,
		}, {
			name:           "User not found",
			urlParam:       "abc",
			shouldCallMock: false,
			mockReturn:     nil,
			mockError:      errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"not found"}`,
		}, {
			name:           "Missing ID",
			urlParam:       "",
			mockReturn:     nil,
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
			if tt.urlParam == "" {
				reqURL = "/users/"
			}

			req, err := http.NewRequest(http.MethodGet, reqURL, nil)
			require.NoError(t, err)

			// Manual chi context injection (consistent with your DeleteUser test)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handlers.GetUser(rr, req) // Assumes your handler method is named GetUser

			require.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

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
			expectedBody:   `"invalid user ID"`,
		}, {
			name:           "User not found",
			urlParam:       "999",
			shouldCallMock: true,
			mockError:      errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"not found"}`,
		}, {
			name:           "Missing ID",
			urlParam:       "",
			shouldCallMock: false,
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
