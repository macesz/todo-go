
# Testing HTTP Handlers in Go: Complete Beginner's Guide

## Table of Contents
1. [Introduction](#introduction)
2. [Basic Testing Setup](#basic-testing-setup)
3. [Testing Patterns](#testing-patterns)
4. [Mocking Dependencies](#mocking-dependencies)
5. [URL Parameters with Chi Router](#url-parameters-with-chi-router)
6. [Common Pitfalls](#common-pitfalls)
7. [Best Practices](#best-practices)
8. [Interview Questions](#interview-questions)
9. [Exercises](#exercises)

## Introduction

HTTP handler testing is crucial for building reliable web APIs in Go. This guide covers everything from basic handler testing to advanced scenarios with mocks and URL parameters.

### Why Test HTTP Handlers?

- **Verify API behavior**: Ensure your endpoints return correct status codes and responses
- **Catch regressions**: Prevent breaking changes when refactoring
- **Document behavior**: Tests serve as living documentation
- **Build confidence**: Deploy with assurance that your API works correctly

## Basic Testing Setup

### Essential Imports

```go
import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)
```

### Simple Handler Test

```go
// Handler to test
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"alive": true}`))
}

// Test function
func TestHealthCheckHandler(t *testing.T) {
    // 1. Create HTTP request
    req, err := http.NewRequest("GET", "/health", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    // 2. Create response recorder
    rr := httptest.NewRecorder()
    
    // 3. Create handler and call it
    handler := http.HandlerFunc(HealthCheckHandler)
    handler.ServeHTTP(rr, req)
    
    // 4. Assert status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    
    // 5. Assert response body
    expected := `{"alive": true}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
```

### Key Components Explained

- **`httptest.NewRequest()`**: Creates HTTP request for testing
- **`httptest.NewRecorder()`**: Records HTTP response for inspection
- **`handler.ServeHTTP()`**: Invokes the handler with request/response
- **`rr.Code`**: Gets HTTP status code from response
- **`rr.Body.String()`**: Gets response body as string

## Testing Patterns

### 1. Table-Driven Tests

**Benefits**: Test multiple scenarios efficiently, reduce code duplication

```go
func TestGetTodo(t *testing.T) {
    tests := []struct {
        name           string
        urlParam       string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "Valid ID",
            urlParam:       "1",
            expectedStatus: http.StatusOK,
            expectedBody:   `{"id":1,"title":"Test Todo"}`,
        },
        {
            name:           "Invalid ID",
            urlParam:       "abc",
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"id must be an integer"}`,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation here
        })
    }
}
```

### 2. Testing Request Body

For handlers that accept JSON input:

```go
func TestCreateTodo(t *testing.T) {
    inputBody := `{"title": "New Todo"}`
    
    // Create request with body
    req, err := http.NewRequest("POST", "/todos", strings.NewReader(inputBody))
    if err != nil {
        t.Fatal(err)
    }
    
    // Set Content-Type header
    req.Header.Set("Content-Type", "application/json")
    
    // Rest of test...
}
```

### 3. Testing Different HTTP Methods

```go
// GET request
req := httptest.NewRequest("GET", "/todos", nil)

// POST request with body
req := httptest.NewRequest("POST", "/todos", strings.NewReader(body))

// PUT request with body
req := httptest.NewRequest("PUT", "/todos/1", strings.NewReader(body))

// DELETE request
req := httptest.NewRequest("DELETE", "/todos/1", nil)
```

## Mocking Dependencies

### Why Mock?

- **Isolation**: Test handlers independently of database/external services
- **Speed**: Tests run faster without I/O operations
- **Control**: Simulate error conditions and edge cases
- **Reliability**: Tests don't depend on external state

### Service Interface Pattern

```go
// Define service interface
type TodoService interface {
    GetTodo(ctx context.Context, id int) (domain.Todo, error)
    CreateTodo(ctx context.Context, title string) (domain.Todo, error)
}

// Handler struct with dependency
type TodoHandlers struct {
    Service TodoService
}

// Handler method
func (h *TodoHandlers) GetTodo(w http.ResponseWriter, r *http.Request) {
    // Handler implementation using h.Service
}
```

### Using Testify Mock

```go
func TestGetTodo(t *testing.T) {
    // Create mock service
    mockService := mocks.NewTodoService(t)
    
    // Set expectations
    mockService.On("GetTodo", mock.Anything, 1).Return(
        domain.Todo{ID: 1, Title: "Test"}, nil)
    
    // Create handler with mock
    handlers := &TodoHandlers{Service: mockService}
    
    // Test handler
    // ...
    
    // Verify mock was called
    mockService.AssertExpectations(t)
}
```

### Mock Best Practices

1. **Create fresh mocks per test**: Avoid test interference
2. **Set precise expectations**: Use specific parameter values when possible
3. **Use `mock.Anything` sparingly**: Only when parameter value doesn't matter
4. **Always assert expectations**: Verify mocks were called as expected

## URL Parameters with Chi Router

### The Challenge

Chi router extracts URL parameters (like `/todos/{id}`) and stores them in request context. In tests, we need to manually set this context.

### Wrong Way ❌

```go
// This won't work - URLParam returns empty string
req := httptest.NewRequest("GET", "/todos/1", nil)
handlers.GetTodo(rr, req) // id will be empty!
```

### Right Way ✅

```go
func TestGetTodo(t *testing.T) {
    req := httptest.NewRequest("GET", "/todos/1", nil)
    
    // Create chi route context
    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("id", "1")  // Set the parameter
    
    // Add context to request
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
    
    // Now URLParam will work in handler
    handlers.GetTodo(rr, req)
}
```

### URL Parameter Testing Pattern

```go
tests := []struct {
    urlParam string
    // other fields
}{
    {"1", /* ... */},
    {"abc", /* ... */},
    {"", /* ... */},
}

for _, tt := range tests {
    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("id", tt.urlParam)
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}
```

## Common Pitfalls

### 1. Double Handler Calls

❌ **Wrong**:
```go
handlers.GetTodo(rr, req)        // First call
handler.ServeHTTP(rr, req)       // Second call - duplicates response!
```

✅ **Correct**:
```go
handlers.GetTodo(rr, req)        // Single call
```

### 2. Wrong URL Construction

❌ **Wrong**:
```go
url := "/todos/1"
req := httptest.NewRequest("GET", "/todos/"+url, nil)  // "/todos//todos/1"
```

✅ **Correct**:
```go
urlParam := "1"
req := httptest.NewRequest("GET", "/todos/"+urlParam, nil)  // "/todos/1"
```

### 3. Forgetting Content-Type Header

❌ **Wrong**:
```go
req := httptest.NewRequest("POST", "/todos", strings.NewReader(body))
// Handler can't parse JSON without Content-Type!
```

✅ **Correct**:
```go
req := httptest.NewRequest("POST", "/todos", strings.NewReader(body))
req.Header.Set("Content-Type", "application/json")
```

### 4. Inconsistent Error Response Format

❌ **Wrong**:
```go
// Handler returns: "error message"
// Test expects: {"error": "error message"}
```

✅ **Correct**:
```go
// Ensure handler and test expect same format
writeJSON(w, status, map[string]string{"error": err.Error()})
```

### 5. Mock Contamination Between Tests

❌ **Wrong**:
```go
var mockService *mocks.TodoService  // Shared across tests

func TestA(t *testing.T) {
    mockService.On("GetTodo", ...).Return(...)
}

func TestB(t *testing.T) {
    // TestB affected by TestA's expectations!
}
```

✅ **Correct**:
```go
func TestA(t *testing.T) {
    mockService := mocks.NewTodoService(t)  // Fresh mock per test
    // ...
}
```

## Best Practices

### 1. Test Structure

Use **AAA pattern** (Arrange, Act, Assert):

```go
func TestGetTodo(t *testing.T) {
    // ARRANGE
    mockService := mocks.NewTodoService(t)
    mockService.On("GetTodo", mock.Anything, 1).Return(todo, nil)
    handlers := &TodoHandlers{Service: mockService}
    
    // ACT
    req := httptest.NewRequest("GET", "/todos/1", nil)
    rr := httptest.NewRecorder()
    handlers.GetTodo(rr, req)
    
    // ASSERT
    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, expectedJSON, rr.Body.String())
    mockService.AssertExpectations(t)
}
```

### 2. Use Descriptive Test Names

```go
// Good test names
"TestGetTodo/Valid_ID_Returns_Todo"
"TestGetTodo/Invalid_ID_Returns_BadRequest"
"TestGetTodo/Nonexistent_ID_Returns_NotFound"

// Poor test names
"TestGetTodo/Test1"
"TestGetTodo/Test2"
```

### 3. Test Edge Cases

Always test:
- **Happy path**: Normal, successful operations
- **Validation errors**: Invalid input data
- **Not found scenarios**: Resources that don't exist
- **Malformed requests**: Invalid JSON, wrong content types
- **Service errors**: When dependencies fail

### 4. Use Helper Functions

```go
func setupTodoTest(t *testing.T, mockSetup func(*mocks.TodoService)) (*TodoHandlers, *httptest.ResponseRecorder) {
    mockService := mocks.NewTodoService(t)
    if mockSetup != nil {
        mockSetup(mockService)
    }
    
    handlers := &TodoHandlers{Service: mockService}
    rr := httptest.NewRecorder()
    
    return handlers, rr
}
```

### 5. Consistent Response Formats

Ensure all handlers return consistent JSON structure:

```go
// Success response
{"id": 1, "title": "Todo", "done": false}

// Error response
{"error": "error message"}
```

## Interview Questions

### Basic Questions

**Q1: What is `httptest.ResponseRecorder` and why do we use it?**

**A:** `ResponseRecorder` implements `http.ResponseWriter` and records HTTP responses for testing. It captures status codes, headers, and body content so we can assert on them in tests.

**Q2: How do you test a handler that expects JSON input?**

**A:** Create request with JSON body using `strings.NewReader()`, set `Content-Type: application/json` header, and pass it to the handler.

**Q3: What's the difference between unit testing and integration testing for HTTP handlers?**

**A:** 
- **Unit testing**: Test handler in isolation using mocks for dependencies
- **Integration testing**: Test handler with real dependencies (database, external APIs)

### Intermediate Questions

**Q4: How do you test URL parameters with chi router?**

**A:** Manually create chi route context and add URL parameters:
```go
rctx := chi.NewRouteContext()
rctx.URLParams.Add("id", "1")
req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
```

**Q5: Why should you create fresh mocks for each test case?**

**A:** To avoid test contamination. Shared mocks carry expectations from previous tests, causing unpredictable test failures.

**Q6: How do you test error scenarios when dependencies fail?**

**A:** Use mocks to return errors:
```go
mockService.On("GetTodo", mock.Anything, 1).Return(domain.Todo{}, errors.New("database error"))
```

### Advanced Questions

**Q7: How would you test middleware?**

**A:** Create a test handler, wrap it with middleware, and verify the middleware behavior:
```go
testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
})
wrappedHandler := myMiddleware(testHandler)
```

**Q8: How do you test handlers that depend on request context values?**

**A:** Set context values in the test request:
```go
ctx := context.WithValue(req.Context(), "userID", 123)
req = req.WithContext(ctx)
```

**Q9: What are table-driven tests and when should you use them?**

**A:** Table-driven tests use a slice of test cases to test multiple scenarios with the same test logic. Use them when testing similar scenarios with different inputs/outputs.

## Exercises

### Exercise 1: Basic Handler Test
Write a test for this simple handler:

```go
func PingHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("pong"))
}
```

<details>
<summary>Solution</summary>

```go
func TestPingHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/ping", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(PingHandler)
    handler.ServeHTTP(rr, req)
    
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    
    expected := "pong"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
```
</details>

### Exercise 2: JSON Handler Test
Write a test for this handler that accepts JSON:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if req.Name == "" || req.Email == "" {
        http.Error(w, "Name and email required", http.StatusBadRequest)
        return
    }
    
    response := map[string]interface{}{
        "id":    1,
        "name":  req.Name,
        "email": req.Email,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

<details>
<summary>Solution</summary>

```go
func TestCreateUserHandler(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "Valid input",
            requestBody:    `{"name":"John","email":"john@example.com"}`,
            expectedStatus: http.StatusOK,
            expectedBody:   `{"id":1,"name":"John","email":"john@example.com"}`,
        },
        {
            name:           "Missing name",
            requestBody:    `{"email":"john@example.com"}`,
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Name and email required\n",
        },
        {
            name:           "Invalid JSON",
            requestBody:    `{"name":"John"`,
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Invalid JSON\n",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, err := http.NewRequest("POST", "/users", strings.NewReader(tt.requestBody))
            if err != nil {
                t.Fatal(err)
            }
            req.Header.Set("Content-Type", "application/json")

            rr := httptest.NewRecorder()
            handler := http.HandlerFunc(CreateUserHandler)
            handler.ServeHTTP(rr, req)

            if status := rr.Code; status != tt.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.expectedStatus)
            }

            if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.expectedBody) {
                t.Errorf("handler returned unexpected body: got %v want %v",
                    rr.Body.String(), tt.expectedBody)
            }
        })
    }
}
```
</details>

### Exercise 3: Mocked Handler Test
Create a test for a handler that uses a service dependency. Include both success and error cases.

## Summary

HTTP handler testing in Go requires understanding:

1. **Basic testing setup** with `httptest` package
2. **Table-driven tests** for comprehensive coverage
3. **Mocking dependencies** for isolated unit tests
4. **URL parameter handling** with routers like chi
5. **Common pitfalls** and how to avoid them
6. **Best practices** for maintainable tests

Master these concepts and you'll be well-prepared for both real-world development and technical interviews!

## Further Reading

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [HTTP Test Package](https://golang.org/pkg/net/http/httptest/)
- [Chi Router Documentation](https://github.com/go-chi/chi)
