package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      string
		setupMock      func(m *mocks.UserService)
		expectedStatus int
		checkResponse  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name:      "Successful login",
			inputBody: `{"email":"test@example.com","password":"Password123"}`,
			setupMock: func(m *mocks.UserService) {
				m.On("Login", mock.Anything,
					"test@example.com",
					"Password123").Return(&domain.User{
					ID:    1,
					Name:  "Test User",
					Email: "test@example.com",
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.LoginResponseDTO
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err, "Should parse response JSON")

				assert.NotEmpty(t, response.Token, "Token should not be empty")

				assert.Equal(t, int64(1), response.User.ID)
				assert.Equal(t, "Test User", response.User.Name)
				assert.Equal(t, "test@example.com", response.User.Email)

				assert.True(t, strings.Count(response.Token, ".") == 2, "Token should have 3 parts separated by dots")
			},
		},
		{
			name:      "Invalid credentials - wrong password",
			inputBody: `{"email":"test@example.com","password":"WrongPassword123"}`,
			setupMock: func(m *mocks.UserService) {
				m.On("Login",
					mock.Anything,
					"test@example.com",
					"WrongPassword123",
				).Return(nil, domain.ErrInvalidCredentials).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid credentials", response.Error)
			},
		}, {
			name:           "Invalid JSON",
			inputBody:      `{"email":"test@example.com"`, // Malformed JSON
			setupMock:      func(m *mocks.UserService) {}, // No service call
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid request body", response.Error)
			},
		},
		{
			name:      "Internal server error",
			inputBody: `{"email":"test@example.com","password":"Password123"}`,
			setupMock: func(m *mocks.UserService) {
				m.On("Login",
					mock.Anything,
					"test@example.com",
					"Password123",
				).Return(nil, errors.New("database connection failed")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "internal server error", response.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewUserService(t)
			tt.setupMock(mockService)

			// Create JWT auth (with test secret)
			tokenAuth := jwtauth.New("HS256", []byte("test-secret-key-for-testing"), nil)
			handlers := &UserHandlers{
				Service:   mockService,
				TokenAuth: tokenAuth,
			}

			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handlers.Login(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code mismatch")

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}
			mockService.AssertExpectations(t)
		})
	}
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
			inputPassword:  "Password123",
			shouldCallMock: true,
			inputBody:      `{"name":"Test User","email":"test@example.com","password":"Password123"}`,
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
			inputPassword:  "Password123",
			inputBody:      `{"name":"Test User","email":"test@example.com","password":"Password123"}`,
			shouldCallMock: true,
			mockReturn:     nil,
			mockError:      errors.New("database failure"), // Generic error → 500
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
		{
			name:           "Missing Name",
			inputBody:      `{"email":"test@example.com","password":"Password123"}`, // Valid JSON, missing name
			shouldCallMock: false,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Name is required"}`,
		}, {
			name:           "Missing Email",
			inputBody:      `{"name":"Test User","password":"Password123"}`, // Valid JSON, missing email
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
