
# Testing Service Layer in Go: Complete Beginner's Guide

## Table of Contents
1. [Introduction](#introduction)
2. [Service Layer Architecture](#service-layer-architecture)
3. [Dependency Injection Pattern](#dependency-injection-pattern)
4. [Mocking Dependencies](#mocking-dependencies)
5. [Table-Driven Tests with initMocks](#table-driven-tests-with-initmocks)
6. [Parallel Testing](#parallel-testing)
7. [Testing CRUD Operations](#testing-crud-operations)
8. [Error Handling & Edge Cases](#error-handling--edge-cases)
9. [Best Practices](#best-practices)
10. [Common Pitfalls](#common-pitfalls)
11. [Interview Questions](#interview-questions)
12. [Exercises](#exercises)

## Introduction

Service layer testing focuses on testing business logic while isolating external dependencies like databases, APIs, or file systems. This guide teaches you how to write comprehensive, maintainable service tests in Go.

### What is the Service Layer?

The service layer sits between your handlers (presentation) and data access (repository/store). It contains:
- **Business logic**: Rules and workflows
- **Data transformation**: Converting between domain models
- **Orchestration**: Coordinating multiple operations
- **Validation**: Business rule validation

### Why Test the Service Layer?

- **Business Logic Verification**: Ensure core functionality works correctly
- **Independence**: Test without external dependencies
- **Speed**: Fast tests without I/O operations
- **Reliability**: Consistent results without external state
- **Documentation**: Tests serve as business requirement documentation

## Service Layer Architecture

### Typical Structure

```go
// Domain model
type Todo struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Done      bool      `json:"done"`
    CreatedAt time.Time `json:"createdAt"`
}

// Store interface (dependency)
type TodoStore interface {
    List(ctx context.Context) ([]Todo, error)
    Create(ctx context.Context, title string) (Todo, error)
    Get(ctx context.Context, id int) (Todo, error)
    Update(ctx context.Context, id int, title string, done bool) (Todo, error)
    Delete(ctx context.Context, id int) error
}

// Service struct with dependency
type TodoService struct {
    Store TodoStore  // Dependency injection
}

// Service methods
func (s *TodoService) ListTodos(ctx context.Context) ([]Todo, error) {
    return s.Store.List(ctx)
}

func (s *TodoService) CreateTodo(ctx context.Context, title string) (Todo, error) {
    // Business logic can go here (validation, transformation, etc.)
    if strings.TrimSpace(title) == "" {
        return Todo{}, errors.New("title cannot be empty")
    }
    
    return s.Store.Create(ctx, title)
}
```

### Key Principles

1. **Dependency Injection**: Accept interfaces, not concrete types
2. **Single Responsibility**: Each service handles one domain area
3. **Context Propagation**: Pass context for cancellation and timeouts
4. **Error Handling**: Return meaningful errors for business failures

## Dependency Injection Pattern

### The Problem Without DI

```go
// BAD: Hard-coded dependency
type TodoService struct {
    db *sql.DB  // Directly couples to database
}

func (s *TodoService) ListTodos() ([]Todo, error) {
    // Direct database queries - hard to test!
    rows, err := s.db.Query("SELECT * FROM todos")
    // ...
}
```

**Problems:**
- ❌ Hard to test (requires real database)
- ❌ Tightly coupled to specific implementation
- ❌ Cannot easily swap implementations
- ❌ Tests are slow and fragile

### The Solution With DI

```go
// GOOD: Interface-based dependency
type TodoStore interface {
    List(ctx context.Context) ([]Todo, error)
    // Other methods...
}

type TodoService struct {
    Store TodoStore  // Interface, not concrete type
}

func NewTodoService(store TodoStore) *TodoService {
    return &TodoService{Store: store}
}
```

**Benefits:**
- ✅ Easy to test with mocks
- ✅ Loosely coupled
- ✅ Swappable implementations
- ✅ Fast, reliable tests

## Mocking Dependencies

### Using testify/mock

Install testify:
```bash
go get github.com/stretchr/testify/mock
```

### Mock Generation

Use `mockery` to generate mocks from interfaces:

```bash
# Install mockery
go install github.com/vektra/mockery/v2@latest

# Generate mock from interface
mockery --name=TodoStore --dir=. --output=mocks
```

This generates:
```go
// mocks/TodoStore.go
type TodoStore struct {
    mock.Mock
}

func (m *TodoStore) List(ctx context.Context) ([]domain.Todo, error) {
    ret := m.Called(ctx)
    return ret.Get(0).([]domain.Todo), ret.Error(1)
}
// ... other methods
```

### Mock Usage Patterns

```go
func TestListTodos(t *testing.T) {
    // Create mock
    mockStore := mocks.NewTodoStore(t)
    
    // Set expectations
    mockStore.On("List", mock.Anything).Return(
        []Todo{{ID: 1, Title: "Test"}}, nil)
    
    // Create service with mock
    service := &TodoService{Store: mockStore}
    
    // Test service
    todos, err := service.ListTodos(context.Background())
    
    // Verify results
    assert.NoError(t, err)
    assert.Len(t, todos, 1)
    
    // Verify mock was called
    mockStore.AssertExpectations(t)
}
```

## Table-Driven Tests with initMocks

### The Pattern

Table-driven tests with `initMocks` provide clean, maintainable test structure:

```go
func TestCreateTodo(t *testing.T) {
    type args struct {
        ctx   context.Context
        title string
    }
    
    tests := []struct {
        name      string
        args      args
        initMocks func(t *testing.T, args *args, s *TodoService)
        want      Todo
        wantErr   bool
    }{
        {
            name: "success",
            args: args{
                ctx:   context.Background(),
                title: "New Todo",
            },
            initMocks: func(t *testing.T, ta *args, s *TodoService) {
                store := mocks.NewTodoStore(t)
                store.On("Create", ta.ctx, ta.title).Return(
                    Todo{ID: 1, Title: ta.title}, nil)
                s.Store = store
            },
            want: Todo{ID: 1, Title: "New Todo"},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := &TodoService{}
            tt.initMocks(t, &tt.args, s)
            
            got, err := s.CreateTodo(tt.args.ctx, tt.args.title)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

### Benefits of initMocks Pattern

1. **Clean Test Cases**: Logic separated from setup
2. **Flexible Mocking**: Different mock behavior per test
3. **Reusable**: Pattern works for any service test
4. **Readable**: Clear what each test is doing
5. **Maintainable**: Easy to add new test cases

### initMocks Function Signature

```go
// Standard signature for initMocks
initMocks func(t *testing.T, args *args, s *ServiceType)
```

**Parameters:**
- `t *testing.T`: For creating fresh mocks
- `args *args`: Access to test arguments for mock setup
- `s *ServiceType`: Service instance to inject mocks into

## Parallel Testing

### Why Use Parallel Tests?

```go
func TestListTodos(t *testing.T) {
    t.Parallel()  // Enable parallel execution
    
    // Test implementation...
}
```

**Benefits:**
- ⚡ **Faster test execution**: Tests run concurrently
- 🔄 **Better CPU utilization**: Utilize multiple cores
- ⏰ **Reduced waiting time**: Especially important for larger test suites

### Parallel Testing Rules

1. **Independent tests only**: No shared state between tests
2. **Fresh mocks**: Create new mocks for each test
3. **Isolated resources**: No shared files, databases, etc.
4. **Thread-safe code**: Service code must be thread-safe

### Parallel Pattern

```go
func TestService(t *testing.T) {
    t.Parallel()  // Parent test parallel
    
    tests := []struct{
        // test cases...
    }{}
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // Each subtest parallel too
            
            // Test implementation with fresh mocks
        })
    }
}
```

## Testing CRUD Operations

### Create Operation Testing

```go
func TestCreateTodo(t *testing.T) {
    tests := []struct {
        name      string
        title     string
        mockSetup func(*mocks.TodoStore, string)
        want      Todo
        wantErr   bool
    }{
        {
            name:  "valid title",
            title: "New Todo",
            mockSetup: func(m *mocks.TodoStore, title string) {
                m.On("Create", mock.Anything, title).Return(
                    Todo{ID: 1, Title: title}, nil)
            },
            want:    Todo{ID: 1, Title: "New Todo"},
            wantErr: false,
        },
        {
            name:  "empty title",
            title: "",
            mockSetup: func(m *mocks.TodoStore, title string) {
                // No mock setup - validation should fail before store call
            },
            want:    Todo{},
            wantErr: true,
        },
        {
            name:  "store error",
            title: "New Todo",
            mockSetup: func(m *mocks.TodoStore, title string) {
                m.On("Create", mock.Anything, title).Return(
                    Todo{}, errors.New("database error"))
            },
            want:    Todo{},
            wantErr: true,
        },
    }
    
    // Test execution...
}
```

### Read Operation Testing

```go
func TestGetTodo(t *testing.T) {
    tests := []struct {
        name      string
        id        int
        mockTodo  Todo
        mockError error
        wantErr   bool
    }{
        {
            name:     "existing todo",
            id:       1,
            mockTodo: Todo{ID: 1, Title: "Test Todo"},
            wantErr:  false,
        },
        {
            name:      "non-existent todo",
            id:        999,
            mockTodo:  Todo{},
            mockError: errors.New("not found"),
            wantErr:   true,
        },
    }
    
    // Implementation...
}
```

### Update Operation Testing

```go
func TestUpdateTodo(t *testing.T) {
    tests := []struct {
        name     string
        id       int
        title    string
        done     bool
        mockRet  Todo
        mockErr  error
        wantErr  bool
    }{
        {
            name:    "successful update",
            id:      1,
            title:   "Updated",
            done:    true,
            mockRet: Todo{ID: 1, Title: "Updated", Done: true},
        },
        {
            name:    "todo not found",
            id:      999,
            title:   "Updated",
            done:    true,
            mockErr: errors.New("not found"),
            wantErr: true,
        },
    }
    
    // Implementation...
}
```

### Delete Operation Testing

```go
func TestDeleteTodo(t *testing.T) {
    tests := []struct {
        name    string
        id      int
        mockErr error
        wantErr bool
    }{
        {
            name:    "successful delete",
            id:      1,
            wantErr: false,
        },
        {
            name:    "todo not found", 
            id:      999,
            mockErr: errors.New("not found"),
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockStore := mocks.NewTodoStore(t)
            mockStore.On("Delete", mock.Anything, tt.id).Return(tt.mockErr)
            
            service := &TodoService{Store: mockStore}
            err := service.DeleteTodo(context.Background(), tt.id)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Error Handling & Edge Cases

### Types of Errors to Test

1. **Validation Errors**: Invalid input data
2. **Not Found Errors**: Resources that don't exist
3. **Store Errors**: Database/persistence failures
4. **Business Rule Violations**: Domain-specific constraints

### Context Testing

```go
func TestWithCancelledContext(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel immediately
    
    mockStore := mocks.NewTodoStore(t)
    // Mock may not be called if service checks context first
    
    service := &TodoService{Store: mockStore}
    _, err := service.ListTodos(ctx)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "context canceled")
}
```

### Timeout Testing

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer cancel()
    
    mockStore := mocks.NewTodoStore(t)
    mockStore.On("List", mock.Anything).Return(nil, context.DeadlineExceeded)
    
    service := &TodoService{Store: mockStore}
    _, err := service.ListTodos(ctx)
    
    assert.Error(t, err)
    assert.Equal(t, context.DeadlineExceeded, err)
}
```

## Best Practices

### 1. Test Structure - AAA Pattern

```go
func TestCreateTodo(t *testing.T) {
    // ARRANGE
    mockStore := mocks.NewTodoStore(t)
    mockStore.On("Create", mock.Anything, "Test").Return(
        Todo{ID: 1, Title: "Test"}, nil)
    service := &TodoService{Store: mockStore}
    
    // ACT  
    result, err := service.CreateTodo(context.Background(), "Test")
    
    // ASSERT
    assert.NoError(t, err)
    assert.Equal(t, "Test", result.Title)
    mockStore.AssertExpectations(t)
}
```

### 2. Fresh Mocks Per Test

```go
// GOOD: Fresh mock per test
func TestListTodos(t *testing.T) {
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockStore := mocks.NewTodoStore(t)  // Fresh mock
            tt.initMocks(t, &tt.args, service)
        })
    }
}
```

### 3. Meaningful Test Names

```go
// GOOD test names
"TestCreateTodo/Success_ReturnsCreatedTodo"
"TestCreateTodo/EmptyTitle_ReturnsValidationError"
"TestCreateTodo/StoreError_ReturnsError"

// POOR test names  
"TestCreateTodo/Test1"
"TestCreateTodo/Test2"
```

### 4. Test Edge Cases

Always test:
- **Happy path**: Normal successful operations
- **Validation failures**: Invalid input
- **Not found scenarios**: Resources that don't exist
- **Store failures**: Persistence layer errors
- **Context cancellation**: Timeout and cancellation scenarios

### 5. Assert Mock Expectations

```go
func TestService(t *testing.T) {
    mockStore := mocks.NewTodoStore(t)
    // ... set up expectations
    
    // ... execute service method
    
    // ALWAYS verify mocks were called as expected
    mockStore.AssertExpectations(t)
}
```

### 6. Use Specific Mock Matchers

```go
// GOOD: Specific expectations
mockStore.On("Create", mock.Anything, "specific title").Return(...)

// ACCEPTABLE: When value doesn't matter for test
mockStore.On("List", mock.Anything).Return(...)

// AVOID: Overuse of mock.Anything
mockStore.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
```

## Common Pitfalls

### 1. Shared Mock State

❌ **Wrong**:
```go
var sharedMock *mocks.TodoStore  // Shared across tests

func TestA(t *testing.T) {
    sharedMock.On("List", ...).Return(...)
}

func TestB(t *testing.T) {
    // TestB contaminated by TestA's expectations!
}
```

✅ **Correct**:
```go
func TestA(t *testing.T) {
    mockStore := mocks.NewTodoStore(t)  // Fresh per test
}
```

### 2. Not Testing Error Cases

❌ **Wrong**:
```go
// Only testing success case
func TestCreateTodo(t *testing.T) {
    // Only happy path test
}
```

✅ **Correct**:
```go
func TestCreateTodo(t *testing.T) {
    tests := []struct{}{
        {"success", /* ... */},
        {"validation_error", /* ... */},
        {"store_error", /* ... */},
    }
}
```

### 3. Ignoring Context

❌ **Wrong**:
```go
func (s *TodoService) CreateTodo(ctx context.Context, title string) {
    // Ignoring context - no cancellation support
    return s.Store.Create(ctx, title)
}
```

✅ **Correct**:
```go
func (s *TodoService) CreateTodo(ctx context.Context, title string) {
    // Check context before expensive operations
    if err := ctx.Err(); err != nil {
        return Todo{}, err
    }
    
    return s.Store.Create(ctx, title)
}
```

### 4. Over-Mocking

❌ **Wrong**:
```go
// Mocking everything, even simple operations
func TestBusinessLogic(t *testing.T) {
    mockValidator := mocks.NewValidator(t)
    mockLogger := mocks.NewLogger(t) 
    mockTimer := mocks.NewTimer(t)
    // ... too many mocks for simple logic
}
```

✅ **Correct**:
```go
// Only mock external dependencies
func TestBusinessLogic(t *testing.T) {
    mockStore := mocks.NewTodoStore(t)  // External dependency
    // Use real objects for simple logic
}
```

### 5. Not Asserting Mock Expectations

❌ **Wrong**:
```go
func TestService(t *testing.T) {
    mock.On("Method", ...).Return(...)
    service.DoSomething()
    // Missing: mock.AssertExpectations(t)
}
```

✅ **Correct**:
```go
func TestService(t *testing.T) {
    mock.On("Method", ...).Return(...)
    service.DoSomething()
    mock.AssertExpectations(t)  // Verify mock was called
}
```

## Interview Questions

### Basic Questions

**Q1: What is the service layer and why do we test it separately?**

**A:** The service layer contains business logic and sits between handlers and data access. We test it separately to:
- Verify business rules in isolation
- Test without external dependencies
- Ensure fast, reliable tests
- Document business requirements through tests

**Q2: What is dependency injection and how does it help testing?**

**A:** Dependency injection provides dependencies through interfaces rather than hard-coding them. Benefits for testing:
- Easy to mock dependencies
- Loose coupling between components  
- Swappable implementations
- Fast tests without I/O

**Q3: How do you create a mock in Go using testify?**

**A:**
```go
// Generate with mockery
mockStore := mocks.NewTodoStore(t)

// Set expectations
mockStore.On("Method", args).Return(result, error)

// Use in service
service := &Service{Store: mockStore}

// Verify calls
mockStore.AssertExpectations(t)
```

### Intermediate Questions

**Q4: What is the initMocks pattern and why is it useful?**

**A:** The initMocks pattern separates mock setup from test case definition:
```go
initMocks: func(t *testing.T, args *args, s *Service) {
    mockStore := mocks.NewTodoStore(t)
    mockStore.On("Method", args.param).Return(result, nil)
    s.Store = mockStore
}
```

Benefits: cleaner test cases, flexible mock setup, reusable pattern.

**Q5: When should you use parallel testing and what are the requirements?**

**A:** Use `t.Parallel()` when:
- Tests are independent (no shared state)
- Want faster test execution
- Tests don't modify global variables
- Fresh mocks created per test

Requirements: thread-safe code, isolated resources, no race conditions.

**Q6: How do you test context cancellation in services?**

**A:**
```go
ctx, cancel := context.WithCancel(context.Background())
cancel()

_, err := service.Method(ctx)
assert.Error(t, err)
assert.Contains(t, err.Error(), "context canceled")
```

### Advanced Questions

**Q7: How would you test a service that orchestrates multiple store operations?**

**A:** Test transaction-like behavior:
```go
func TestComplexOperation(t *testing.T) {
    mockStore := mocks.NewTodoStore(t)
    
    // Set up multiple expectations in order
    mockStore.On("Get", mock.Anything, 1).Return(todo, nil).Once()
    mockStore.On("Update", mock.Anything, 1, "new title", true).Return(updatedTodo, nil).Once()
    mockStore.On("Create", mock.Anything, "audit log").Return(auditTodo, nil).Once()
    
    service := &TodoService{Store: mockStore}
    err := service.ComplexOperation(ctx, 1)
    
    assert.NoError(t, err)
    mockStore.AssertExpectations(t) // Verifies all calls were made
}
```

**Q8: How do you test business logic that doesn't directly use external dependencies?**

**A:** Create unit tests without mocks for pure business logic:
```go
func TestValidateTitle(t *testing.T) {
    tests := []struct {
        title   string
        wantErr bool
    }{
        {"valid", false},
        {"", true},
        {strings.Repeat("a", 256), true},
    }
    
    for _, tt := range tests {
        err := validateTitle(tt.title)
        assert.Equal(t, tt.wantErr, err != nil)
    }
}
```

**Q9: How would you test a service with multiple dependencies?**

**A:** Mock each dependency separately:
```go
type TodoService struct {
    Store TodoStore
    Cache CacheService
    Logger Logger
}

func TestWithMultipleDeps(t *testing.T) {
    mockStore := mocks.NewTodoStore(t)
    mockCache := mocks.NewCacheService(t) 
    mockLogger := mocks.NewLogger(t)
    
    // Set expectations on each mock
    mockStore.On("Get", mock.Anything, 1).Return(todo, nil)
    mockCache.On("Set", mock.Anything, mock.Anything).Return(nil)
    mockLogger.On("Info", mock.Anything)
    
    service := &TodoService{
        Store: mockStore,
        Cache: mockCache, 
        Logger: mockLogger,
    }
}
```

## Exercises

### Exercise 1: Basic Service Test

Write tests for this simple service:

```go
type UserService struct {
    Store UserStore
}

type UserStore interface {
    GetByEmail(ctx context.Context, email string) (User, error)
}

func (s *UserService) FindUserByEmail(ctx context.Context, email string) (User, error) {
    if email == "" {
        return User{}, errors.New("email is required")
    }
    
    return s.Store.GetByEmail(ctx, email)
}
```

<details>
<summary>Solution</summary>

```go
func TestFindUserByEmail(t *testing.T) {
    tests := []struct {
        name      string
        email     string
        mockUser  User
        mockError error
        wantErr   bool
    }{
        {
            name:     "valid email",
            email:    "user@example.com", 
            mockUser: User{ID: 1, Email: "user@example.com"},
            wantErr:  false,
        },
        {
            name:    "empty email",
            email:   "",
            wantErr: true,
        },
        {
            name:      "user not found",
            email:     "notfound@example.com",
            mockError: errors.New("user not found"),
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockStore := mocks.NewUserStore(t)
            
            if tt.email != "" {
                mockStore.On("GetByEmail", mock.Anything, tt.email).
                    Return(tt.mockUser, tt.mockError)
            }
            
            service := &UserService{Store: mockStore}
            
            user, err := service.FindUserByEmail(context.Background(), tt.email)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.mockUser, user)
            }
            
            mockStore.AssertExpectations(t)
        })
    }
}
```
</details>

### Exercise 2: Service with Business Logic

Write tests for a service that has business rules:

```go
func (s *TodoService) CompleteTodo(ctx context.Context, id int) error {
    todo, err := s.Store.Get(ctx, id)
    if err != nil {
        return err
    }
    
    if todo.Done {
        return errors.New("todo is already completed")
    }
    
    _, err = s.Store.Update(ctx, id, todo.Title, true)
    return err
}
```

### Exercise 3: Service with Multiple Operations

Write tests for a service method that performs multiple store operations:

```go
func (s *TodoService) TransferTodo(ctx context.Context, fromList, toList int, todoID int) error {
    // Get todo from source list
    todo, err := s.Store.GetFromList(ctx, fromList, todoID)
    if err != nil {
        return err
    }
    
    // Remove from source list
    if err := s.Store.RemoveFromList(ctx, fromList, todoID); err != nil {
        return err
    }
    
    // Add to destination list
    if err := s.Store.AddToList(ctx, toList, todo); err != nil {
        // Rollback: add back to source list
        s.Store.AddToList(ctx, fromList, todo)
        return err
    }
    
    return nil
}
```

## Summary

Service layer testing in Go requires understanding:

1. **Dependency injection** for testable architecture
2. **Mocking patterns** with testify/mock
3. **Table-driven tests** with initMocks for clean structure
4. **Parallel testing** for faster execution
5. **CRUD operation testing** patterns
6. **Error handling** and edge case coverage
7. **Best practices** for maintainable tests

Master these concepts and you'll write robust, maintainable service tests that provide confidence in your business logic while remaining fast and reliable!

## Further Reading

- [testify/mock Documentation](https://pkg.go.dev/github.com/stretchr/testify/mock)
- [mockery Tool](https://github.com/vektra/mockery)
- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Context Package](https://golang.org/pkg/context/)
- [Dependency Injection in Go](https://blog.golang.org/wire)
