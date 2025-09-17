# Service Layer Testing in Go — Internal Guide

This guide is a practical introduction to testing the service layer in Go. It uses your TodoService tests as a reference point and expands on key concepts, patterns, and techniques to write fast, reliable, and maintainable tests.

## Goals

- Understand what to test in the service layer and why
- Learn table-driven testing with subtests and mocks
- Use testify’s mocking and assertions effectively
- Handle contexts, time, and parallelism safely
- Avoid flaky tests and common pitfalls
- Run, debug, and measure tests with confidence

---

## What is the Service Layer?

The service layer holds your application’s business logic. It sits between transport (HTTP, gRPC, CLI) and the data layer (stores, repositories).

- It orchestrates workflows and enforces rules
- It passes `context.Context` downstream
- It translates and wraps data-layer errors into domain errors
- It remains testable by depending on interfaces, not concrete implementations

Testing the service layer focuses on your logic, not the database or HTTP. We usually use mocks for the store to isolate behaviors.

---

## Test Types Overview

- **Unit tests (focus here)**: Isolate the service logic by mocking dependencies. Fast and deterministic.
- **Integration tests**: Use a real DB or test container; validate repository and queries.
- **End-to-end tests**: Full stack (transport + service + store); slower, used for high-confidence workflows.

For unit tests, aim for:
- Happy path
- Expected error paths
- Edge cases and boundaries
- Context cancellation/timeout behavior
- Input validation

---

## Patterns in Your Current Tests

You’re already using several good patterns:

- **Table-driven tests** with `tests := []struct{ ... }{ ... }`
- **Subtests** with `t.Run(...)`
- **Parallelism** with `t.Parallel()`
- **Mocks** with testify’s `On(...).Return(...).Once()`
- Clear separation of `fields`, `args`, `want`, and `initMocks`

These are solid foundations. Below are recommendations to level up.

---

## Recommendations and Improvements

### 1) Use require/assert to simplify comparisons
Manually comparing slices and structs is verbose and brittle:

```go
if len(res) != len(tt.want) { ... }
for i := range res {
  if res[i] != tt.want[i] { ... }
}
```

Prefer testify:
```go
require.NoError(t, err)
require.Equal(t, tt.want, res) // Great diffs on failure
```

If order is not important use `ElementsMatch`. If you need deeply customized comparison, consider `google/go-cmp`.

### 2) Always assert mock expectations
Ensure every expected call happened:

```go
store := mocks.NewTodoStore(tt)
tt.Cleanup(func() { store.AssertExpectations(tt) })
```

Doing this in `tt.Cleanup` guarantees it runs even if the test fails early.

### 3) Use fixed time values for determinism
Avoid `time.Now()` in tests. Instead, use a fixed timestamp:
```go
fixed := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
```
Use this in your expected results and mock returns. This ensures tests are deterministic and don’t fail due to timing issues.

### 4) Handle time safely
Comparing `time.Time` using `==` is risky due to monotonic clock bits. You currently capture `now` once and reuse it, which works, but a better approach is:

- Inject a `Clock` interface into your service and use a fake clock in tests, or
- Use `require.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Millisecond)`

Example clock approach:
```go
type Clock interface { Now() time.Time }
type realClock struct{}
func (realClock) Now() time.Time { return time.Now() }

type TodoService struct {
  Store domain.TodoStore
  Clock Clock
}
```

### 5) Validate inputs and test validation paths
Add and test cases like empty title, excessively long title, invalid ID, etc. This belongs in the service layer.

### 6) Context handling tests
Verify that the service respects cancellation and timeouts:

- Use `context.WithTimeout`
- Make the mock observe the ctx (e.g., with `mock.MatchedBy` or `Run` to inspect args)
- Ensure service returns promptly on cancellation

### 7) Avoid overspecifying mocks
Use argument matchers if exact values aren’t essential:
```go
store.On("List", mock.Anything).Return(...).Once()
```
Or to inspect and validate an arg:
```go
store.On("Create", mock.Anything, mock.AnythingOfType("string")).
  Run(func(args mock.Arguments) {
    title := args.String(1)
    require.NotEmpty(tt, title)
  }).
  Return(domain.Todo{...}, nil).
  Once()
```

### 8) Determinism: no randomness, no time dependencies, no shared mutable globals
Parallel tests are great for speed, but avoid shared globals or non-thread-safe singletons. If you must use globals, guard with `sync.Mutex` or avoid `t.Parallel()` for those cases.

---

## An Improved Test Skeleton

Below is a refactor of one test showing some best practices. Adapt similarly for the others.

```go
package todo

import (
  "context"
  "errors"
  "testing"
  "time"

  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/mock"

  "github.com/macesz/todo-go/domain"
  "github.com/macesz/todo-go/services/todo/mocks"
)

func TestListTodos(t *testing.T) {
  t.Parallel()

  type fields struct {
    Store *mocks.TodoStore
  }

  type args struct {
    ctx context.Context
  }

  fixed := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

  tests := []struct {
    name      string
    fields    fields
    args      args
    want      []domain.Todo
    wantErr   bool
    initMocks func(tt *testing.T, ta *args, s *TodoService)
  }{
    {
      name:   "success",
      fields: fields{},
      args:   args{ctx: context.Background()},
      want: []domain.Todo{
        {ID: 1, Title: "Test Todo 1", Done: false, CreatedAt: fixed},
        {ID: 2, Title: "Test Todo 2", Done: true, CreatedAt: fixed},
      },
      initMocks: func(tt *testing.T, ta *args, s *TodoService) {
        store := mocks.NewTodoStore(tt)
        //Here we ensure expectations are checked
        tt.Cleanup(func() { store.AssertExpectations(tt) })

        store.
          On("List", ta.ctx).
          Return([]domain.Todo{
            {ID: 1, Title: "Test Todo 1", Done: false, CreatedAt: fixed},
            {ID: 2, Title: "Test Todo 2", Done: true, CreatedAt: fixed},
          }, nil).
          Once()

        s.Store = store
      },
    },
    {
      name:    "store error",
      fields:  fields{},
      args:    args{ctx: context.Background()},
      wantErr: true,
      initMocks: func(tt *testing.T, ta *args, s *TodoService) {
        store := mocks.NewTodoStore(tt)
        tt.Cleanup(func() { store.AssertExpectations(tt) })

        store.
          On("List", ta.ctx).
          Return(nil, errors.New("not found")).
          Once()

        s.Store = store
      },
    },
  }

  for _, tc := range tests {
    tc := tc
    t.Run(tc.name, func(t *testing.T) {
      t.Parallel()

      s := &TodoService{
        Store: tc.fields.Store,
      }

      tc.initMocks(t, &tc.args, s)

      got, err := s.ListTodos(tc.args.ctx)
      if tc.wantErr {
        require.Error(t, err)
        return
      }
      require.NoError(t, err)
      require.Equal(t, tc.want, got)
    })
  }
}
```

Key changes:
- `require` for concise assertions and better diffs
- `tt.Cleanup(store.AssertExpectations(tt))` to enforce mock expectations
- Fixed time value for determinism
- Removed manual slice comparison

Apply the same style to `CreateTodo`, `GetTodo`, `UpdateTodo`, and `DeleteTodo`.

---

## What to Test in CRUD Services

- List:
  - Empty list
  - Non-empty list
  - Store error propagation
- Create:
  - Valid title
  - Empty/invalid title (validation error)
  - Store error (e.g., conflict)
- Get:
  - Found
  - Not found (translate to domain error)
- Update:
  - Valid update (title and/or done)
  - Not found
  - Validation error (empty title)
- Delete:
  - Success
  - Not found
  - Idempotency (optional policy decision)

Also consider:
- Context cancellation/timeouts
- Idempotent operations where appropriate
- Observability (logging) — test by injecting a logger or observing side effects if necessary

---

## Using Testify Mocks Effectively

Common patterns:
- `On("Method", args...).Return(result, err).Once()`
- Matchers: `mock.Anything`, `mock.AnythingOfType("string")`, `mock.MatchedBy(func(T) bool { ... })`
- Inspect args with `.Run(func(args mock.Arguments) { ... })`
- Verify: `store.AssertExpectations(t)` or `mock.AssertExpectationsForObjects(t, store1, store2, ...)`

Avoid:
- Over-specifying exact args when not relevant
- Leaving expectations un-asserted (can hide missing calls)

---

## Context Testing Example

```go
func TestServiceRespectsContext(t *testing.T) {
  t.Parallel()

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
  t.Cleanup(cancel)

  store := mocks.NewTodoStore(t)
  t.Cleanup(func() { store.AssertExpectations(t) })

  // Ensure service passes the exact ctx through to the store
  store.
    On("List", ctx).
    Return(nil, context.DeadlineExceeded).
    Once()

  s := &TodoService{Store: store}
  _, err := s.ListTodos(ctx)
  require.ErrorIs(t, err, context.DeadlineExceeded)
}
```

If matching the exact `ctx` ref is awkward, relax with `mock.Anything` and assert in `.Run` that the ctx is done.

---

## Running and Measuring Tests

- Run all tests:
  `go test ./...`
- Verbose:
  `go test -v ./...`
- Disable cache (useful while debugging):
  `go test -count=1 ./...`
- Shuffle test order (catches hidden coupling):
  `go test -shuffle=on ./...`
- Race detector:
  `go test -race ./...`
- Coverage:
  `go test -coverprofile=cover.out ./... && go tool cover -func=cover.out`
  `go tool cover -html=cover.out` to open in browser
- Focus on one test:
  `go test -run TestUpdateTodo ./... -v`

---

## Common Pitfalls and How to Avoid Them

- **Flaky tests due to time**: Inject a clock or use fixed times; avoid sleeping.
- **Brittle struct equality with time**: Use `WithinDuration` or `Time.Equal`.
- **Unasserted mocks**: Always call `AssertExpectations`.
- **Parallel tests sharing state**: Avoid shared mutable globals or use locks.
- **Overspecified mocks**: Match only what matters; otherwise tests become high maintenance.
- **Testing the database in unit tests**: Keep unit tests isolated; use integration tests for DB.
- **Comparing large slices manually**: Use `require.Equal` or `cmp.Diff`.

---

## Checklist Before Submitting

- Dependencies are interfaces; service is DI-friendly
- Table-driven tests cover happy, error, and edge cases
- `t.Parallel()` used where safe
- `require/assert` used for clear, fail-fast checks
- Mocks have `.Once()` and `AssertExpectations`
- Deterministic time and no sleeps
- Run with `-race`, `-shuffle=on`, and coverage tracked

---

## Appendix: Minimal Test Template

```go
func TestService_Method(t *testing.T) {
  t.Parallel()

  type fields struct {
    Store *mocks.SomeStore
  }
  type args struct {
    ctx context.Context
    // other inputs...
  }
  tests := []struct {
    name      string
    fields    fields
    args      args
    want      SomeType
    wantErr   bool
    initMocks func(tt *testing.T, ta *args, s *Service)
  }{
    {
      name: "success",
      args: args{ctx: context.Background()},
      initMocks: func(tt *testing.T, ta *args, s *Service) {
        store := mocks.NewSomeStore(tt)
        tt.Cleanup(func() { store.AssertExpectations(tt) })
        store.On("Method", ta.ctx /*...*/).Return(/*...*/).Once()
        s.Store = store
      },
      want: /*...*/,
    },
  }

  for _, tc := range tests {
    tc := tc
    t.Run(tc.name, func(t *testing.T) {
      t.Parallel()
      s := &Service{Store: tc.fields.Store}
      if tc.initMocks != nil {
        tc.initMocks(t, &tc.args, s)
      }
      got, err := s.Method(tc.args.ctx /*...*/)
      if tc.wantErr {
        require.Error(t, err)
        return
      }
      require.NoError(t, err)
      require.Equal(t, tc.want, got)
    })
  }
}
```

With these patterns, your service layer tests will be fast, clear, and robust.
