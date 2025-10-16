
# Go API Error Handling and Best Practices: Lessons from Todo App

This document summarizes the key concepts learned while improving the `UpdateTodo` handler in a Go-based Todo API. It covers error handling, the differences between structs and maps for responses, why and how we improved the handler, and why initial tests failed (with fixes). These lessons are beginner-friendly, with analogies to languages like JavaScript or Java where helpful.

## 1. Introduction
We started with a basic Todo API handler using Go, Chi router, and libraries like `go-playground/validator`. The goal was to make it more RESTful, consistent, and robust. Key improvements included:
- Better error handling with custom errors.
- Consistent response formats (e.g., using structs over maps).
- Dynamic validation messages.
- Proper HTTP status codes (e.g., 404 for "not found" instead of generic 400).

This led to test failures, which we debugged and fixed. The process taught core Go concepts like error management and data structures.

## 2. Error Handling in Go
Go handles errors differently from languages with exceptions (e.g., Java's `try-catch` or JS's `throw`). Errors are just values returned by functions.

### Key Concepts
- **Basic Error Checking**: Functions often return an `error` (e.g., `result, err := doSomething()`). Check `if err != nil { handle it }`. This encourages immediate handling, unlike deferred exceptions.
- **Custom Errors**: Define reusable errors in your `domain` package using `errors.New("message")`. Example:
  ```go
  var ErrNotFound = errors.New("todo not found")
  ```
  - Why? Allows type-safe checks without string matching (e.g., avoid `if err.Error() == "not found"`—fragile if messages change).
- **Checking Custom Errors**: Use `errors.Is(err, domain.ErrNotFound)` (like `instanceof` in JS/Java). This compares the error's identity.
- **RESTful Error Responses**:
  - Map errors to HTTP codes: 400 (bad input), 404 (not found), 500 (internal server error).
  - For security, use generic messages for 500 (e.g., "internal server error")—log the real details server-side.
  - Consistent format: Always return JSON like `{"error": "message"}` (via map or struct).

### Why It Matters
Proper handling makes APIs predictable (clients know what went wrong) and secure (no leaking internals). In our Todo app, we used custom errors like `ErrNotFound` to return 404 specifically, improving over generic 400s.

## 3. Struct vs. Map for API Responses
Both structs and maps group data (e.g., for JSON responses like `{"error": "message"}`), but they differ in structure and use.

### Quick Definitions
- **Map**: Like a JS object `{}` or Python dict. Flexible key-value pairs (e.g., `map[string]string{"error": "msg"}`).
- **Struct**: Like a Java class with fixed fields. Defined as a type (e.g., `type ErrorResponse struct { Error string json:"error" }`), then instantiated (e.g., `ErrorResponse{Error: "msg"}`).

### Pros and Cons Table

| Aspect              | Map (e.g., `map[string]string{"error": "msg"}`) | Struct (e.g., `ErrorResponse{Error: "msg"}`) |
|---------------------|------------------------------------------------|---------------------------------------------|
| **Type Safety**     | Con: No compile-time checks (typos like "eror" fail at runtime). | Pro: Go checks fields at compile time. |
| **Flexibility**     | Pro: Add keys dynamically (quick for prototypes). | Con: Fixed fields (requires updating the type to change). |
| **Extensibility**   | Con: Hard to add fields consistently across code. | Pro: Add a field once (e.g., `Code string`), and it works everywhere. |
| **Ease for Beginners** | Pro: Simple, no setup. | Con: Define type first, but reusable after. |
| **Consistency**     | Con: Easy to vary formats accidentally. | Pro: Enforces uniform structure. |
| **Performance**     | Neutral: Slightly faster, but negligible. | Neutral: Efficient in Go. |
| **When to Use**     | Quick scripts or one-offs. | Real apps/APIs for maintainability. |

### Why Struct is Often Better
In our app, we switched from maps to `domain.ErrorResponse` for consistency (all errors use the same type) and extensibility (easy to add fields like "code"). Both encode to the same JSON, but structs are more "Go-like" (idiomatic) for repeated structures.

## 4. Improving the UpdateTodo Handler
### What We Improved
- **Old Handler Issues**:
  - Hardcoded validation messages (e.g., fixed string for title errors).
  - Inconsistent error formats (mix of maps, plain strings).
  - Generic status codes (always 400 for service errors, even "not found").
  - No resource cleanup (e.g., not closing request body).
  - Exposed internal errors (risky for security).

- **New Handler Features**:
  - Dynamic validation: Use `err := validate.New().Struct(dto); if err != nil { return {Error: err.Error()} }`—pulls messages from DTO tags (e.g., "required,min=1").
  - Custom error checks: `if errors.Is(err, domain.ErrNotFound) { 404 } else { 500 with generic message }`.
  - Consistent structs: Use `domain.ErrorResponse` for all errors.
  - Added `defer r.Body.Close()` for cleanup.
  - Early returns for clean code.

### Why We Improved It
- **RESTful Principles**: Right status codes help clients (e.g., retry on 500, show "not found" UI on 404).
- **Consistency**: Easier for clients to parse (always check `response.error`).
- **Maintainability**: Dynamic messages avoid hardcoding; structs scale better.
- **Security/Best Practices**: Hide internals, prevent leaks.
- **Result**: Handler is more robust, like a professional API (e.g., aligns with standards in Express.js or Spring Boot).

Example Snippet (Improved Handler):
```go
if err := validate.New().Struct(todoDTO); err != nil {
    writeJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
    return
}
// ... service call with errors.Is checks
```

## 5. Understanding and Fixing Test Failures
### Why Tests Failed Initially
We used table-driven tests (structs defining scenarios) with mocks (simulating service) and assertions (checking status/body). Failures happened because handler changes (dynamic errors, custom checks) didn't match old expectations.

- **Invalid JSON**: Expected hardcoded "invalid JSON", but got actual `json.Decode` error ("unexpected EOF"). Why? Handler now returns `err.Error()` dynamically.
- **Invalid Data (Empty Title)**: Expected hardcoded message, but got validator's dynamic one (e.g., "failed on 'required' tag"). Why? Switched to `err.Error()` for flexibility.
- **Service Error**: Expected 404 with "not found", but got 500. Why? Mock returned `errors.New("not found")`, but handler checks `errors.Is(err, domain.ErrNotFound)`—messages match, but errors are different values (not identical).

Concept: Tests must match real behavior. Mocks need exact errors for `errors.Is`; dynamic messages require updating expectations.

### How We Fixed Them
- Updated `expectedBody` to match dynamic/real errors (e.g., validator's string for validation).
- Changed mock to return `domain.ErrNotFound` (exact custom error) for 404 case.
- Added cases (e.g., "title too long", internal error) for coverage.
- Kept `assert.JSONEq` for flexible JSON comparison (ignores whitespace).

Example Fix (Service Error Case):
```go
mockError:      domain.ErrNotFound, // Exact custom error
expectedStatus: http.StatusNotFound,
expectedBody:   `{"error":"todo not found"}`, // Matches ErrNotFound.Error()
```

## 6. Key Takeaways
- Error handling in Go is explicit and value-based—use custom errors for safety.
- Prefer structs over maps for consistent, extensible API responses.
- Improve handlers for RESTfulness: Dynamic validation, proper codes, cleanup.
- Tests fail when expectations don't match changes—update them to reflect reality.
- Practice: Apply these to other handlers/tests. Tools like Postman help verify.

This knowledge builds a solid foundation for Go backends. Refer back as you code!

