
# Complete Beginner's Guide: Go + PostgreSQL for Job Interviews

## Table of Contents
1. [Why Go + PostgreSQL?](#why-go--postgresql)
2. [Go Database Fundamentals](#go-database-fundamentals)
3. [Database Drivers & Libraries Explained](#database-drivers--libraries-explained)
4. [Connection Management Deep Dive](#connection-management-deep-dive)
5. [Understanding Query Patterns](#understanding-query-patterns)
6. [Error Handling in Go](#error-handling-in-go)
7. [Security Best Practices](#security-best-practices)
8. [Testing Your Database Code](#testing-your-database-code)
9. [Performance & Optimization](#performance--optimization)
10. [Advanced Concepts](#advanced-concepts)
11. [Interview Questions & Detailed Answers](#interview-questions--detailed-answers)
12. [Complete Working Example](#complete-working-example)

---

## 1. Why Go + PostgreSQL?

### Understanding the Match

**Go's Strengths:**
- **Compiled language** = Fast execution (unlike interpreted languages like Python)
- **Strong typing** = Catches errors at compile time, not runtime
- **Goroutines** = Handle thousands of concurrent database connections efficiently
- **Simple syntax** = Easy to read and maintain
- **Great standard library** = Built-in database support

**PostgreSQL's Strengths:**
- **ACID compliance** = Reliable transactions
- **Rich data types** = JSON, arrays, custom types
- **Advanced features** = Full-text search, GIS, window functions
- **Performance** = Excellent query optimizer

**Why They Work Together:**
```go
// Go's type system matches PostgreSQL's precision, and Go's simplicity makes it easy to work with PostgreSQL's rich data types.

type User struct {
    ID        int64     `json:"id" db:"id"`           // PostgreSQL BIGINT
    Email     string    `json:"email" db:"email"`     // PostgreSQL VARCHAR
    CreatedAt time.Time `json:"created_at" db:"created_at"` // PostgreSQL TIMESTAMP
}
```

---

## 2. Go Database Fundamentals

### 2.1 How Go Handles Databases

Unlike some languages that have built-in database support, Go uses a **driver-based approach**:
The difference between the driver-based approach and built-in database support is that the driver-based approach allows for more flexibility and compatibility with various databases, while built-in database support is more limited and specific to a particular database.

```go
import (
    "database/sql"           // Standard database interface
    _ "github.com/lib/pq"   // PostgreSQL driver (note the underscore!)
)
```

**The underscore import (`_`) explained:**
- We don't directly use functions from `lib/pq`
- We just need it to **register** the PostgreSQL driver with Go's `database/sql` package, (register means to add the driver to the list of available drivers that Go's `database/sql` package can use.)
- Go's `database/sql` package will use it automatically

### 2.2 The Database/SQL Interface Pattern

Go's `database/sql` package defines **interfaces**, not implementations:

```go
// This is what database/sql provides (simplified):
type DB interface {
    Query(query string, args ...interface{}) (*Rows, error)
    Exec(query string, args ...interface{}) (Result, error)
    // ... more methods
}
```

which implements the `DB` interface, and provides a consistent API for interacting with databases.

**What this means for you:**
- All PostgreSQL drivers implement the same interface
- You can switch databases without changing your code structure
- Testing becomes easier (you can mock the interface)

### 2.3 Understanding Go's Error Handling

Go doesn't have exceptions like Java or Python. Instead, functions return errors:

```go
// Instead of try-catch blocks, Go uses explicit error checking
result, err := db.Query("SELECT * FROM users")
if err != nil {
    // Handle error - this is MANDATORY in Go
    return err
}
// Continue with result...
```

**Why this is better for database code:**
- **Explicit** = You must handle every potential error
- **Predictable** = No hidden exceptions that crash your program
- **Debuggable** = Easy to trace where errors originate

---

## 3. Database Drivers & Libraries Explained

### 3.1 The Three-Layer Architecture

```
Your Application Code
        ↓
   Library (sqlx/GORM)
        ↓
   database/sql (Go standard)
        ↓
   Driver (lib/pq)
        ↓
   PostgreSQL Database
```

### 3.2 Detailed Library Comparison

#### Option 1: `database/sql` (Standard Library)

**What it is:** Go's built-in database package

```go
// Example with database/sql - VERBOSE but explicit
import "database/sql"

func GetUser(db *sql.DB, id int) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE id = $1"
    row := db.QueryRow(query, id)

    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)  // Manual field mapping
    if err != nil {
        return nil, err
    }

    return &user, nil
}
```

**Pros:**
- ✅ **No dependencies** - part of Go standard library
- ✅ **Stable** - won't break with Go updates
- ✅ **Explicit** - you see exactly what's happening

**Cons:**
- ❌ **Verbose** - lots of boilerplate code
- ❌ **Manual mapping** - you must map each field by hand
- ❌ **No advanced features** - no struct scanning, named parameters

#### Option 2: `sqlx`  // better choice

**What it is:** Extensions to `database/sql` that add convenience features

```go
// Same example with sqlx - MUCH cleaner
import "github.com/jmoiron/sqlx"

func GetUser(db *sqlx.DB, id int) (*User, error) {
    var user User
    query := "SELECT id, name, email FROM users WHERE id = $1"
    err := db.Get(&user, query, id)  // Automatic struct mapping!
    return &user, err
}
```

**Why this is The best choice:**

**Pros:**
- ✅ **Familiar** - same as database/sql but enhanced
- ✅ **Struct scanning** - automatically maps database rows to Go structs
- ✅ **Named parameters** - use `:name` instead of `$1, $2, $3`
- ✅ **Still explicit** - you write SQL, you control queries
- ✅ **Learning friendly** - teaches you SQL while providing Go convenience

**Cons:**
- ❌ **External dependency** - not part of Go standard library
- ❌ **Still requires SQL knowledge** - you write raw SQL queries

#### Option 3: `GORM` (Full ORM)

**What it is:** Object-Relational Mapping - hides SQL behind Go methods

```go
// GORM example - looks more like other languages
type User struct {
    gorm.Model
    Name  string
    Email string
}

func GetUser(db *gorm.DB, id int) (*User, error) {
    var user User
    result := db.First(&user, id)  // No SQL visible!
    return &user, result.Error
}
```

**When to use GORM:**
- ✅ **Rapid prototyping** - get started quickly
- ✅ **Simple CRUD operations** - basic create/read/update/delete
- ✅ **Familiar to ORM users** - similar to Django ORM, ActiveRecord, Java Spring Boot,

**When to avoid GORM:**
- ❌ **Complex queries** - harder to optimize performance
- ❌ **Learning SQL** - hides the SQL you need to know for interviews
- ❌ **Performance critical** - adds overhead and abstraction


**Cons:** by Uncle Bob

1. "ORMs are the Vietnam of Computer Science"
He famously said ORMs often become quagmires - they start simple but become increasingly complex and hard to escape from.

2. Violation of Dependency Inversion Principle
// ❌ What Uncle Bob warns against - business logic depending on ORM

3. Database as an Implementation Detail
Uncle Bob argues that your business logic should not know or care about the database:

"The database is a detail. The web is a detail. They are not the center of your system. They are plugins."

4. The "Object-Relational Impedance Mismatch"

What Uncle Bob Recommends Instead
1. Explicit SQL with Repository Pattern!!!

### 3.3 Why This Code Uses the Perfect Approach

Looking at the example code structure:

```go
type Store struct {
    queryTemplates map[string]*template.Template  // ✅ SQL templates for flexibility
    db             *sqlx.DB                       // ✅ sqlx for enhanced features
}
```

**This approach combines the best of both worlds:**
- ✅ **Learning-friendly** - you see and write real SQL
- ✅ **Maintainable** - SQL queries are separate from Go code
- ✅ **Flexible** - templates allow dynamic queries when needed
- ✅ **Performance** - no ORM overhead, direct SQL control
- ✅ **Interview-ready** - shows understanding of both Go and SQL

---

## 4. Connection Management Deep Dive

### 4.1 Understanding Database Connections

**What is a database connection?**
A connection is like a phone line between your Go application and PostgreSQL:

```go
// This creates a CONNECTION POOL, not a single connection
db, err := sqlx.Connect("postgres", "postgres://user:pass@localhost/mydb")
```

**Important:** `sqlx.Connect()` doesn't create one connection - it creates a **pool** of connections that your application can use.

### 4.2 Connection Pool Explained

Think of a connection pool like a taxi fleet:

```go
func setupDatabase() (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, err
    }

    // Configure the "taxi fleet"
    db.SetMaxOpenConns(25)    // Maximum 25 "taxis" can be working at once
    db.SetMaxIdleConns(10)    // Keep 10 "taxis" running even when not busy
    db.SetConnMaxLifetime(5 * time.Minute) // Replace old "taxis" every 5 minutes

    return db, nil
}
```

**Why connection pooling matters:**
- **Performance** - Don't create/destroy connections for every query
- **Resource management** - PostgreSQL has limits on concurrent connections
- **Reliability** - Reuse tested, working connections

### 4.3 DSN (Data Source Name) Breakdown

```go
// Full format breakdown
dsn := "postgres://username:password@hostname:port/database?param1=value1&param2=value2"

// Real example
dsn := "postgres://todouser:secret123@localhost:5432/todoapp?sslmode=disable&timezone=UTC"
```

**Each part explained:**
- `postgres://` - Database type (could be `mysql://`, etc.)
- `todouser:secret123` - Username and password
- `localhost:5432` - Host and port (5432 is PostgreSQL's default)
- `todoapp` - Database name
- `?sslmode=disable` - Connection parameters (disable SSL for development)

### 4.4 Environment-Based Configuration

**Never hardcode connection strings!**

```go
// ❌ BAD - credentials in source code
dsn := "postgres://user:password@localhost/db"

// ✅ GOOD - use environment variables
import "os"

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Database string
    SSLMode  string
}

func NewDatabaseConfig() DatabaseConfig {
    return DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnv("DB_PORT", "5432"),
        User:     getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", ""),
        Database: getEnv("DB_NAME", "todoapp"),
        SSLMode:  getEnv("DB_SSLMODE", "disable"),
    }
}

func (c DatabaseConfig) DSN() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
        c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

**In your `.env` file:**
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=todouser
DB_PASSWORD=secret123
DB_NAME=todoapp
DB_SSLMODE=disable
```

---

## 5. Understanding Query Patterns

### 5.1 The Five Essential Query Types

Every database application uses these five patterns. Let's understand each one:

#### Pattern 1: Single Row Read (`Get`, `QueryRow`)

**When to use:** Get one specific record by ID, email, etc.

```go
// Example: Get one todo by ID
func (s *Store) GetTodo(ctx context.Context, id int) (domain.Todo, error) {
    var todo domain.Todo
    query := `SELECT id, title, done, created_at FROM todos WHERE id = :id`

    err := s.db.GetContext(ctx, &todo, query, map[string]any{"id": id})
    if err == sql.ErrNoRows {
        return domain.Todo{}, ErrTodoNotFound  // Custom error for "not found"
    }
    if err != nil {
        return domain.Todo{}, err  // Database error
    }

    return todo, nil
}
```

**Key points:**
- Use `GetContext()` or `QueryRowContext()` for single rows
- Always handle `sql.ErrNoRows` (this means no record was found)
- Return domain-specific errors, not database errors

#### Pattern 2: Multiple Rows Read (`Select`, `Query`)

**When to use:** Get a list of records, search results, etc.

```go
// Example: Get all todos for a user
func (s *Store) ListTodos(ctx context.Context, userID int) ([]domain.Todo, error) {
    todos := make([]domain.Todo, 0)  // Initialize empty slice

    query := `SELECT id, title, done, created_at
              FROM todos
              WHERE user_id = :user_id
              ORDER BY created_at DESC`

    rows, err := s.db.NamedQueryContext(ctx, query, map[string]any{"user_id": userID})
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // ALWAYS close rows when done!

    for rows.Next() {
        var todo domain.Todo
        if err := rows.StructScan(&todo); err != nil {
            return nil, err
        }
        todos = append(todos, todo)
    }

    // Check if the loop ended due to an error
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return todos, nil
}
```

**Key points:**
- Always `defer rows.Close()` - this prevents memory leaks!
- Check `rows.Err()` after the loop - the loop might have ended due to an error
- Initialize slice with `make([]Type, 0)` for consistent JSON serialization

#### Pattern 3: Insert with Returning ID

**When to use:** Create a new record and get its generated ID

```go
// Example: Create a new todo and return its ID
func (s *Store) CreateTodo(ctx context.Context, todo domain.Todo) (int, error) {
    query := `INSERT INTO todos (title, user_id, done)
              VALUES (:title, :user_id, :done)
              RETURNING id`

    params := map[string]any{
        "title":   todo.Title,
        "user_id": todo.UserID,
        "done":    todo.Done,
    }

    var id int
    err := s.db.GetContext(ctx, &id, query, params)
    if err != nil {
        return 0, err
    }

    return id, nil
}
```

**Key points:**
- PostgreSQL's `RETURNING` clause lets you get the generated ID in one query
- Use `GetContext()` because `RETURNING` gives you exactly one row back
- Always check for constraint violations (covered in error handling section)

#### Pattern 4: Update with Affected Rows Check

**When to use:** Modify existing records

```go
// Example: Update a todo
func (s *Store) UpdateTodo(ctx context.Context, todo domain.Todo) error {
    query := `UPDATE todos
              SET title = :title, done = :done
              WHERE id = :id AND user_id = :user_id`  // Notice: user_id for security!

    params := map[string]any{
        "id":      todo.ID,
        "title":   todo.Title,
        "done":    todo.Done,
        "user_id": todo.UserID,  // Prevent users from updating others' todos
    }

    result, err := s.db.NamedExecContext(ctx, query, params)
    if err != nil {
        return err
    }

    // Check if the update actually modified a row
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrTodoNotFound  // No rows were updated
    }

    return nil
}
```

**Key points:**
- Always check `RowsAffected()` to know if the update actually happened
- Include user ownership in WHERE clause for security
- Use `NamedExecContext()` for operations that don't return data

#### Pattern 5: Delete with Affected Rows Check

**When to use:** Remove records

```go
// Example: Delete a todo
func (s *Store) DeleteTodo(ctx context.Context, id int, userID int) error {
    query := `DELETE FROM todos WHERE id = :id AND user_id = :user_id`

    params := map[string]any{
        "id":      id,
        "user_id": userID,  // Security: users can only delete their own todos
    }

    result, err := s.db.NamedExecContext(ctx, query, params)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrTodoNotFound
    }

    return nil
}
```

### 5.2 Parameter Binding Styles Explained

#### Style 1: Positional Parameters (database/sql style)

```go
// Parameters are numbered: $1, $2, $3...
query := "SELECT * FROM todos WHERE user_id = $1 AND done = $2"
rows, err := db.Query(query, userID, false)
```

**Pros:**
- ✅ Standard library compatible
- ✅ Slightly faster (less processing)

**Cons:**
- ❌ Hard to read with many parameters
- ❌ Easy to mix up parameter order
- ❌ Difficult to maintain

#### Style 2: Named Parameters (sqlx style) ⭐ RECOMMENDED

```go
// Parameters have names: :user_id, :done
query := "SELECT * FROM todos WHERE user_id = :user_id AND done = :done"
params := map[string]any{
    "user_id": userID,
    "done":    false,
}
rows, err := db.NamedQuery(query, params)
```

**Pros:**
- ✅ **Readable** - clear what each parameter is for
- ✅ **Maintainable** - easy to add/remove/reorder parameters
- ✅ **Less error-prone** - can't mix up parameter positions
- ✅ **Self-documenting** - parameter names explain their purpose

**Cons:**
- ❌ Requires sqlx library (not standard library)
- ❌ Slightly more processing overhead

### 5.3 The Template Pattern (Your Advanced Approach)

Looking at the example code, there's an advanced pattern being used:

```go
type Store struct {
    queryTemplates map[string]*template.Template  // SQL templates
    db             *sqlx.DB
}
```

**What this does:**
1. **Separates SQL from Go code** - SQL lives in `.sql` files
2. **Allows dynamic queries** - templates can change based on conditions
3. **Improves maintainability** - DBAs can review SQL separately from Go code
4. **Enables testing** - you can test SQL templates independently

**Example usage:**
```go
// In queries/get_todos.sql.tpl
SELECT id, title, done, created_at
FROM todos
WHERE user_id = :user_id
{{if .IncludeDeleted}}
  -- Include deleted todos if requested
{{else}}
  AND deleted_at IS NULL
{{end}}
ORDER BY created_at DESC
{{if .Limit}}LIMIT :limit{{end}}

// In Go code:
templateParams := map[string]any{
    "IncludeDeleted": false,
    "Limit":         true,
}
queryStr, err := prepareQuery(s.queryTemplates["get_todos"], templateParams)

queryParams := map[string]any{
    "user_id": userID,
    "limit":   50,
}
```

**Why this is advanced:**
- ✅ **Flexible** - same template can generate different queries
- ✅ **DRY** - don't repeat similar SQL queries
- ✅ **Secure** - template params are not user input (query params are)
- ✅ **Professional** - shows architectural thinking

---

## 6. Error Handling in Go

### 6.1 Go's Error Philosophy

Go's approach to error handling is different from most languages:

```go
// Other languages (Java, Python, C#):
try {
    user = database.getUser(id);
    // continue with user...
} catch (Exception e) {
    // handle error
}

// Go approach:
user, err := database.GetUser(id)
if err != nil {
    // handle error immediately
    return err
}
// continue with user...
```

**Why Go does this:**
- **Explicit** - you can't ignore errors accidentally
- **Local** - error handling happens right where the error occurs
- **Predictable** - no hidden control flow from exceptions

### 6.2 Database-Specific Errors

#### The Special Case: `sql.ErrNoRows`

This is the most common database error you'll encounter:

```go
import "database/sql"

func (s *Store) GetTodo(ctx context.Context, id int) (domain.Todo, error) {
    var todo domain.Todo
    err := s.db.GetContext(ctx, &todo, "SELECT * FROM todos WHERE id = :id",
                          map[string]any{"id": id})

    if err == sql.ErrNoRows {
        // This is NOT an error condition - it just means "not found"
        return domain.Todo{}, ErrTodoNotFound  // Return your own error type
    }

    if err != nil {
        // This IS an error - database problem, network issue, etc.
        return domain.Todo{}, fmt.Errorf("failed to get todo %d: %w", id, err)
    }

    return todo, nil
}
```

**Why `sql.ErrNoRows` is special:**
- It's not an error condition for your application
- It's just PostgreSQL saying "no rows matched your query"
- You should convert it to a domain-specific error

#### PostgreSQL Constraint Errors

PostgreSQL has specific error codes for constraint violations:

```go
import (
    "github.com/lib/pq"
)

func (s *Store) CreateUser(ctx context.Context, user domain.User) (int, error) {
    query := `INSERT INTO users (email, name) VALUES (:email, :name) RETURNING id`

    var id int
    err := s.db.GetContext(ctx, &id, query, map[string]any{
        "email": user.Email,
        "name":  user.Name,
    })

    if err != nil {
        // Check if it's a PostgreSQL error
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code {
            case "23505":  // unique_violation
                if strings.Contains(pqErr.Message, "email") {
                    return 0, ErrEmailAlreadyExists
                }
                return 0, ErrDuplicateData
            case "23502":  // not_null_violation
                return 0, ErrRequiredFieldMissing
            case "23503":  // foreign_key_violation
                return 0, ErrInvalidReference
            }
        }
        // If it's not a constraint error, it's a system error
        return 0, fmt.Errorf("failed to create user: %w", err)
    }

    return id, nil
}
```

**PostgreSQL Error Codes to Remember:**
- `23505` - Unique constraint violation (duplicate data)
- `23502` - NOT NULL constraint violation
- `23503` - Foreign key constraint violation
- `23514` - Check constraint violation

### 6.3 Custom Error Types

Create your own error types for better error handling:

```go
// Define your application errors
var (
    ErrTodoNotFound        = errors.New("todo not found")
    ErrEmailAlreadyExists  = errors.New("email already exists")
    ErrInvalidInput        = errors.New("invalid input")
    ErrUnauthorized        = errors.New("unauthorized")
)

// Or create error types with more information
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// Usage:
if user.Email == "" {
    return ValidationError{Field: "email", Message: "email is required"}
}
```

### 6.4 Error Wrapping and Context

Go 1.13+ introduced error wrapping for better error context:

```go
func (s *Store) UpdateTodo(ctx context.Context, todo domain.Todo) error {
    result, err := s.db.NamedExecContext(ctx, query, params)
    if err != nil {
        // Wrap the error with context
        return fmt.Errorf("failed to update todo %d: %w", todo.ID, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected for todo %d: %w", todo.ID, err)
    }

    if rowsAffected == 0 {
        return ErrTodoNotFound
    }

    return nil
}
```

**Benefits of error wrapping:**
- **Context** - you can see the chain of where errors occurred
- **Original error preserved** - you can still check for specific error types
- **Debugging** - much easier to trace problems

---

## 7. Security Best Practices

### 7.1 SQL Injection Prevention

**SQL injection** is when malicious SQL code is inserted into your queries. Here's how it works and how to prevent it:

#### The Attack

```go
// ❌ VULNERABLE CODE - NEVER DO THIS!
userInput := "1; DROP TABLE users; --"  // Malicious input
query := fmt.Sprintf("SELECT * FROM todos WHERE user_id = %s", userInput)
// Resulting query: SELECT * FROM todos WHERE user_id = 1; DROP TABLE users; --
// This would delete your entire users table!
```

#### The Prevention: Parameterized Queries

```go
// ✅ SAFE - Using parameterized queries
userInput := "1; DROP TABLE users; --"  // Same malicious input
query := "SELECT * FROM todos WHERE user_id = :user_id"
params := map[string]any{"user_id": userInput}

// PostgreSQL treats the entire string as a single parameter value
// The malicious SQL code becomes harmless data
```

**Why parameterized queries work:**
1. **Separation** - SQL structure is separate from data
2. **Escaping** - Database driver automatically escapes dangerous characters
3. **Type checking** - Parameters are validated against expected types

### 7.2 Understanding the Template vs Parameter Separation

Looking at the example code pattern:

```go
// ✅ CORRECT separation
templateParams := map[string]any{
    "tableName": "todos",      // Safe - developer-controlled
    "orderBy":   "created_at", // Safe - from predefined list
}

queryParams := map[string]any{
    "user_id": userInput,      // Safe - parameterized
    "status":  userStatus,     // Safe - parameterized
}
```

**Template parameters (templateParams):**
- Used for **SQL structure** (table names, column names, ORDER BY clauses)
- Should **NEVER** contain user input
- Used at query preparation time
- Examples: table names, column names, conditional SQL blocks

**Query parameters (queryParams):**
- Used for **data values** (WHERE conditions, INSERT values)
- Safe for user input (when properly parameterized)
- Used at query execution time
- Examples: user IDs, search terms, form data

### 7.3 Authentication and Authorization

#### Row-Level Security

Always include ownership checks in your queries:

```go
// ❌ BAD - any user can access any todo
func (s *Store) GetTodo(ctx context.Context, id int) (domain.Todo, error) {
    query := "SELECT * FROM todos WHERE id = :id"
    // Missing user ownership check!
}

// ✅ GOOD - users can only access their own todos
func (s *Store) GetTodo(ctx context.Context, id int, userID int) (domain.Todo, error) {
    query := "SELECT * FROM todos WHERE id = :id AND user_id = :user_id"
    params := map[string]any{
        "id":      id,
        "user_id": userID,  // Ensures user can only see their own todos
    }
}
```

#### Context-Based User Information

Use Go's context to pass user information through your application:

```go
type contextKey string

const UserIDKey contextKey = "user_id"

// Middleware to extract user from JWT and add to context
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract user ID from JWT token (simplified)
        userID := extractUserFromJWT(r.Header.Get("Authorization"))

        // Add user ID to context
        ctx := context.WithValue(r.Context(), UserIDKey, userID)

        // Continue with user in context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Extract user ID in your handlers
func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(UserIDKey).(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    todoID := extractTodoID(r)  // from URL path
    todo, err := h.store.GetTodo(r.Context(), todoID, userID)
    // ...
}
```

### 7.4 Input Validation

Always validate input before it reaches your database:

```go
type TodoCreateRequest struct {
    Title string `json:"title"`
}

func (r TodoCreateRequest) Validate() error {
    if strings.TrimSpace(r.Title) == "" {
        return ValidationError{Field: "title", Message: "title is required"}
    }

    if len(r.Title) > 255 {
        return ValidationError{Field: "title", Message: "title too long (max 255 characters)"}
    }

    // Check for potentially malicious content
    if containsSQLKeywords(r.Title) {
        return ValidationError{Field: "title", Message: "invalid characters in title"}
    }

    return nil
}

func containsSQLKeywords(s string) bool {
    dangerous := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "UNION"}
    upper := strings.ToUpper(s)
    for _, keyword := range dangerous {
        if strings.Contains(upper, keyword) {
            return true
        }
    }
    return false
}
```

---

## 8. Testing Your Database Code

### 8.1 Testing Strategies Overview

There are three main approaches to testing database code:

1. **Unit Tests with Mocks** - Fast, isolated, but don't test real SQL
2. **Integration Tests with Real Database** - Slower, but test everything
3. **Template Tests** - Test SQL templates separately (your advanced approach!)

### 8.2 Template Testing (Your Approach)

This is an advanced technique that tests your SQL templates without needing a database:

```go
func TestTemplateGetTodo(t *testing.T) {
    // Load your query templates
    queries, err := buildQueries("queries")
    require.NoError(t, err)

    // Test template compilation and rendering
    templateParams := map[string]any{
        "includeDeleted": false,
        "orderBy":       "created_at",
    }

    query, err := prepareQuery(queries["get_todo_query"], templateParams)
    require.NoError(t, err)

    // Verify the generated SQL looks correct
    expected := "SELECT id, title, done, created_at FROM todos WHERE id = :id AND deleted_at IS NULL ORDER BY created_at"
    assert.Contains(t, query, "SELECT id, title")
    assert.Contains(t, query, "WHERE id = :id")
    assert.NotContains(t, query, "deleted_at IS NOT NULL")  // deleted filtering should be applied

    t.Log("Generated query:", query)
}
```

**Benefits of template testing:**
- ✅ **Fast** - no database required
- ✅ **Validates SQL syntax** - catches template errors early
- ✅ **Tests different conditions** - verify templates generate correct SQL for different inputs
- ✅ **CI/CD friendly** - no external dependencies

### 8.3 Integration Testing with Real Database

```go
func TestCreateTodo_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    store := NewStore(db)

    // Test data
    todo := domain.Todo{
        Title:  "Test Todo",
        UserID: 1,
        Done:   false,
    }

    // Execute the test
    id, err := store.CreateTodo(context.Background(), todo)

    // Verify results
    require.NoError(t, err)
    assert.Greater(t, id, 0)

    // Verify it was actually saved
    savedTodo, err := store.GetTodo(context.Background(), id, 1)
    require.NoError(t, err)
    assert.Equal(t, todo.Title, savedTodo.Title)
    assert.Equal(t, todo.UserID, savedTodo.UserID)
    assert.Equal(t, todo.Done, savedTodo.Done)
}

func setupTestDB(t *testing.T) *sqlx.DB {
    // Use a separate test database
    dsn := "postgres://testuser:testpass@localhost/test_todoapp?sslmode=disable"
    db, err := sqlx.Connect("postgres", dsn)
    require.NoError(t, err)

    // Run migrations to set up schema
    err = runMigrations(db)
    require.NoError(t, err)

    return db
}

func cleanupTestDB(t *testing.T, db *sqlx.DB) {
    // Clean up test data
    db.Exec("TRUNCATE TABLE todos")
    db.Close()
}
```

### 8.4 Unit Testing with Mocks

Using `github.com/DATA-DOG/go-sqlmock`:

```go
func TestGetTodo_NotFound(t *testing.T) {
    // Create a mock database
    mockDB, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer mockDB.Close()

    // Set up expectation
    mock.ExpectQuery("SELECT (.+) FROM todos WHERE id = \\$1").
        WithArgs(999).
        WillReturnError(sql.ErrNoRows)

    // Create store with mock
    sqlxDB := sqlx.NewDb(mockDB, "postgres")
    store := &Store{db: sqlxDB}

    // Execute test
    _, err = store.GetTodo(context.Background(), 999)

    // Verify results
    assert.ErrorIs(t, err, ErrTodoNotFound)

    // Verify all expectations were met
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### 8.5 Test Database Best Practices

#### Use Test-Specific Database

```go
// ✅ GOOD - separate test database
const testDSN = "postgres://testuser:testpass@localhost/test_todoapp"

// ❌ BAD - using production/development database for tests
const testDSN = "postgres://user:pass@localhost/todoapp"  // Same as dev!
```

#### Parallel Test Safety

```go
func TestCreateTodo(t *testing.T) {
    t.Parallel()  // This test can run in parallel

    // Use unique data to avoid conflicts
    testID := fmt.Sprintf("test_%d_%d", time.Now().UnixNano(), rand.Int())
    todo := domain.Todo{
        Title: fmt.Sprintf("Test Todo %s", testID),
        // ...
    }
}
```

#### Test Transactions

```go
func TestUpdateTodo_Transaction(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    store := NewStore(db)

    // Start a transaction for this test
    tx, err := db.BeginTxx(context.Background(), nil)
    require.NoError(t, err)
    defer tx.Rollback()  // Always rollback test transactions

    // Create store with transaction
    txStore := &Store{db: tx}

    // Run your tests...
    // All changes will be rolled back automatically
}
```

---

## 9. Performance & Optimization

### 9.1 Understanding Database Performance

Database performance in Go applications typically bottlenecks at:

1. **Network latency** - Time to send query to database
2. **Query execution time** - How long PostgreSQL takes to run your query
3. **Result processing** - Time to scan results into Go structs
4. **Connection overhead** - Time to establish/manage database connections

### 9.2 Connection Pool Optimization

The connection pool is your first line of performance optimization:

```go
func OptimizeConnectionPool(db *sqlx.DB) {
    // Set maximum connections based on your PostgreSQL config
    // PostgreSQL default max_connections = 100
    // Leave room for other applications/admin connections
    db.SetMaxOpenConns(25)

    // Keep some connections "warm" to avoid connection overhead
    // But not too many (uses memory and database resources)
    db.SetMaxIdleConns(10)

    // Rotate connections periodically to avoid long-running connection issues
    db.SetConnMaxLifetime(5 * time.Minute)

    // How long to wait for a connection from the pool
    db.SetConnMaxIdleTime(10 * time.Minute)
}
```

**How to determine optimal values:**

```go
// Monitor your connection pool usage
func MonitorConnectionPool(db *sqlx.DB) {
    stats := db.Stats()

    log.Printf("Connection Pool Stats:")
    log.Printf("  Open connections: %d", stats.OpenConnections)
    log.Printf("  In use: %d", stats.InUse)
    log.Printf("  Idle: %d", stats.Idle)
    log.Printf("  Wait count: %d", stats.WaitCount)           // How often we waited for a connection
    log.Printf("  Wait duration: %v", stats.WaitDuration)     // Total time spent waiting
    log.Printf("  Max idle closed: %d", stats.MaxIdleClosed) // Connections closed due to SetMaxIdleConns
    log.Printf("  Max lifetime closed: %d", stats.MaxLifetimeClosed) // Connections closed due to SetConnMaxLifetime
}
```

### 9.3 Query Optimization Strategies

#### Index Your Queries

For every WHERE clause in your application, consider if you need an index:

```sql
-- If you frequently query by user_id
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- If you often filter by status and user together
CREATE INDEX idx_todos_user_status ON todos(user_id, done);

-- If you search by title
CREATE INDEX idx_todos_title ON todos USING gin(to_tsvector('english', title));

-- If you sort by created_at frequently
CREATE INDEX idx_todos_created_at ON todos(created_at);
```

#### Use EXPLAIN ANALYZE

```go
func (s *Store) AnalyzeQuery(query string, params map[string]any) {
    // Add EXPLAIN ANALYZE to see query performance
    explainQuery := "EXPLAIN ANALYZE " + query

    rows, err := s.db.NamedQuery(explainQuery, params)
    if err != nil {
        log.Printf("Error analyzing query: %v", err)
        return
    }
    defer rows.Close()

    log.Printf("Query analysis for: %s", query)
    for rows.Next() {
        var line string
        rows.Scan(&line)
        log.Printf("  %s", line)
    }
}
```

#### Limit Result Sets

```go
// ✅ GOOD - always limit results for list operations
func (s *Store) ListTodos(ctx context.Context, userID int, limit int) ([]domain.Todo, error) {
    if limit <= 0 || limit > 100 {
        limit = 50  // Reasonable default
    }

    query := `SELECT id, title, done, created_at
              FROM todos
              WHERE user_id = :user_id
              ORDER BY created_at DESC
              LIMIT :limit`

    params := map[string]any{
        "user_id": userID,
        "limit":   limit,
    }

    // ... execute query
}

// ❌ BAD - unlimited results can crash your application
func (s *Store) ListAllTodos(ctx context.Context, userID int) ([]domain.Todo, error) {
    query := "SELECT * FROM todos WHERE user_id = :user_id"  // Could return millions of rows!
}
```

### 9.4 Batch Operations

Instead of executing queries one by one:

```go
// ❌ SLOW - N+1 queries problem
func (s *Store) CreateTodosSlow(ctx context.Context, todos []domain.Todo) error {
    for _, todo := range todos {
        _, err := s.CreateTodo(ctx, todo)  // Each call is a separate database round-trip
        if err != nil {
            return err
        }
    }
    return nil
}

// ✅ FAST - single transaction with batch insert
func (s *Store) CreateTodosBatch(ctx context.Context, todos []domain.Todo) error {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Prepare the statement once
    query := "INSERT INTO todos (title, user_id, done) VALUES (:title, :user_id, :done)"
    stmt, err := tx.PrepareNamed(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute multiple times with different parameters
    for _, todo := range todos {
        _, err = stmt.ExecContext(ctx, map[string]any{
            "title":   todo.Title,
            "user_id": todo.UserID,
            "done":    todo.Done,
        })
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

### 9.5 Context and Timeouts

Always use context with timeouts for database operations:

```go
func (s *Store) GetTodoWithTimeout(id int) (domain.Todo, error) {
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var todo domain.Todo
    err := s.db.GetContext(ctx, &todo, "SELECT * FROM todos WHERE id = :id",
                          map[string]any{"id": id})

    if err != nil {
        // Check if it was a timeout
        if ctx.Err() == context.DeadlineExceeded {
            return domain.Todo{}, errors.New("database query timeout")
        }
        return domain.Todo{}, err
    }

    return todo, nil
}
```

### 9.6 Caching Strategies

For frequently accessed data that doesn't change often:

```go
import (
    "sync"
    "time"
)

type CachedStore struct {
    store *Store
    cache map[string]interface{}
    mutex sync.RWMutex
    ttl   time.Duration
}

type cacheItem struct {
    data    interface{}
    expires time.Time
}

func (cs *CachedStore) GetTodo(ctx context.Context, id int) (domain.Todo, error) {
    cacheKey := fmt.Sprintf("todo_%d", id)

    // Try cache first
    cs.mutex.RLock()
    if item, exists := cs.cache[cacheKey]; exists {
        if cached := item.(cacheItem); time.Now().Before(cached.expires) {
            cs.mutex.RUnlock()
            return cached.data.(domain.Todo), nil
        }
    }
    cs.mutex.RUnlock()

    // Cache miss - get from database
    todo, err := cs.store.GetTodo(ctx, id)
    if err != nil {
        return domain.Todo{}, err
    }

    // Store in cache
    cs.mutex.Lock()
    cs.cache[cacheKey] = cacheItem{
        data:    todo,
        expires: time.Now().Add(cs.ttl),
    }
    cs.mutex.Unlock()

    return todo, nil
}
```

**Note:** In production, use dedicated caching solutions like Redis instead of in-memory maps.

---

## 10. Advanced Concepts

### 10.1 Database Transactions

Transactions ensure that multiple database operations either all succeed or all fail:

```go
func (s *Store) TransferTodoOwnership(ctx context.Context, todoID, fromUserID, toUserID int) error {
    // Start a transaction
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to start transaction: %w", err)
    }

    // Important: ensure transaction is closed
    defer func() {
        if err != nil {
            tx.Rollback()  // Rollback on any error
        }
    }()

    // Step 1: Verify the todo belongs to fromUser
    var currentOwner int
    err = tx.GetContext(ctx, &currentOwner,
        "SELECT user_id FROM todos WHERE id = $1", todoID)
    if err != nil {
        return fmt.Errorf("failed to verify todo ownership: %w", err)
    }

    if currentOwner != fromUserID {
        return errors.New("todo doesn't belong to the specified user")
    }

    // Step 2: Update todo ownership
    result, err := tx.ExecContext(ctx,
        "UPDATE todos SET user_id = $1 WHERE id = $2", toUserID, todoID)
    if err != nil {
        return fmt.Errorf("failed to transfer todo: %w", err)
    }

    // Step 3: Verify the update worked
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to verify update: %w", err)
    }
    if rowsAffected == 0 {
        return errors.New("no rows were updated")
    }

    // Step 4: Log the transfer (audit trail)
    _, err = tx.ExecContext(ctx,
        "INSERT INTO todo_transfers (todo_id, from_user_id, to_user_id, transferred_at) VALUES ($1, $2, $3, NOW())",
        todoID, fromUserID, toUserID)
    if err != nil {
        return fmt.Errorf("failed to log transfer: %w", err)
    }

    // If we get here, everything succeeded - commit the transaction
    err = tx.Commit()
    if err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}
```

**Key transaction concepts:**
- **Atomicity** - All operations succeed or all fail
- **Consistency** - Database remains in valid state
- **Isolation** - Concurrent transactions don't interfere
- **Durability** - Committed changes are permanent

### 10.2 Database Migrations

Managing database schema changes over time:

```go
//go:embed migrations/*.sql
var migrationFS embed.FS

type Migration struct {
    Version int
    Name    string
    SQL     string
}

func GetMigrations() ([]Migration, error) {
    files, err := migrationFS.ReadDir("migrations")
    if err != nil {
        return nil, err
    }

    var migrations []Migration
    for _, file := range files {
        if strings.HasSuffix(file.Name(), ".sql") {
            content, err := migrationFS.ReadFile("migrations/" + file.Name())
            if err != nil {
                continue
            }

            // Parse version from filename: 001_create_users.sql
            parts := strings.Split(file.Name(), "_")
            version, err := strconv.Atoi(parts[0])
            if err != nil {
                continue
            }

            migrations = append(migrations, Migration{
                Version: version,
                Name:    strings.TrimSuffix(file.Name(), ".sql"),
                SQL:     string(content),
            })
        }
    }

    // Sort by version
    sort.Slice(migrations, func(i, j int) bool {
        return migrations[i].Version < migrations[j].Version
    })

    return migrations, nil
}

func RunMigrations(db *sqlx.DB) error {
    // Create migrations tracking table
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version INTEGER PRIMARY KEY,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        return fmt.Errorf("failed to create migrations table: %w", err)
    }

    // Get current migration version
    var currentVersion int
    err = db.Get(&currentVersion, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations")
    if err != nil {
        return fmt.Errorf("failed to get current migration version: %w", err)
    }

    // Get all migrations
    migrations, err := GetMigrations()
    if err != nil {
        return fmt.Errorf("failed to load migrations: %w", err)
    }

    // Apply pending migrations
    for _, migration := range migrations {
        if migration.Version <= currentVersion {
            continue  // Already applied
        }

        log.Printf("Applying migration %d: %s", migration.Version, migration.Name)

        // Run migration in transaction
        tx, err := db.Beginx()
        if err != nil {
            return fmt.Errorf("failed to start transaction for migration %d: %w", migration.Version, err)
        }

        // Execute migration SQL
        _, err = tx.Exec(migration.SQL)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
        }

        // Record migration as applied
        _, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
        }

        err = tx.Commit()
        if err != nil {
            return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
        }

        log.Printf("Successfully applied migration %d", migration.Version)
    }

    return nil
}
```

**Example migration files:**

```sql
-- migrations/001_create_users.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrations/002_create_todos.sql
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    done BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_todos_user_id ON todos(user_id);

-- migrations/003_add_todo_priority.sql
ALTER TABLE todos ADD COLUMN priority INTEGER DEFAULT 1;
CREATE INDEX idx_todos_priority ON todos(priority);
```

### 10.3 Connection Context and Cancellation

Understanding how Go's context works with database operations:

```go
func (s *Store) GetTodoWithCancellation(parentCtx context.Context, id int) (domain.Todo, error) {
    // Create a child context with timeout
    ctx, cancel := context.WithTimeout(parentCtx, 3*time.Second)
    defer cancel()

    var todo domain.Todo

    // This query will be cancelled if:
    // 1. The 3-second timeout is reached
    // 2. The parent context is cancelled
    // 3. cancel() is called manually
    err := s.db.GetContext(ctx, &todo, "SELECT * FROM todos WHERE id = :id",
                          map[string]any{"id": id})

    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return domain.Todo{}, errors.New("query timeout after 3 seconds")
        }
        if ctx.Err() == context.Canceled {
            return domain.Todo{}, errors.New("query was cancelled")
        }
        return domain.Todo{}, err
    }

    return todo, nil
}
```

**Real-world example with HTTP request context:**

```go
func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
    // The request context is automatically cancelled if:
    // - Client disconnects
    // - HTTP server shutdown
    // - Request timeout
    ctx := r.Context()

    todoID := extractTodoIDFromURL(r.URL.Path)
    userID := getUserIDFromContext(ctx)

    // Pass the request context to database operations
    // If the client disconnects, the database query will be cancelled
    todo, err := h.store.GetTodo(ctx, todoID, userID)
    if err != nil {
        if ctx.Err() == context.Canceled {
            // Client disconnected - no need to send response
            return
        }
        http.Error(w, "Failed to get todo", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(todo)
}
```

---

## 11. Interview Questions & Detailed Answers

### 11.1 Basic Concepts

**Q: What's the difference between `database/sql` and `sqlx`? Why would you choose one over the other?**

**A:** Both are Go libraries for working with SQL databases, but they serve different needs:

**`database/sql` (Standard Library):**
- **What it is:** Go's built-in database interface
- **Pros:** No external dependencies, guaranteed to be stable, explicit operations
- **Cons:** Verbose, requires manual field mapping, no named parameters

```go
// database/sql example - verbose but explicit
row := db.QueryRow("SELECT id, title FROM todos WHERE id = $1", todoID)
var todo Todo
err := row.Scan(&todo.ID, &todo.Title)  // Manual field mapping
```

**`sqlx` (Extension Library):**
- **What it is:** Adds convenience features on top of database/sql
- **Pros:** Struct scanning, named parameters, less boilerplate, still transparent
- **Cons:** External dependency, slightly more complexity

```go
// sqlx example - cleaner and more maintainable
var todo Todo
err := db.Get(&todo, "SELECT id, title FROM todos WHERE id = :id",
             map[string]any{"id": todoID})
```

**When to choose each:**
- **Choose `database/sql`** for: Libraries, simple applications, when you want zero dependencies
- **Choose `sqlx`** for: Most applications, when you want productivity without hiding SQL

---

**Q: How do you prevent SQL injection in Go?**

**A:** SQL injection is prevented by using **parameterized queries** that separate SQL structure from data:

**❌ Vulnerable approach (string concatenation):**
```go
// NEVER DO THIS - vulnerable to SQL injection
userInput := "1; DROP TABLE users; --"
query := fmt.Sprintf("SELECT * FROM todos WHERE user_id = %s", userInput)
db.Query(query)  // This could execute malicious SQL!
```

**✅ Safe approach (parameterized queries):**
```go
// ALWAYS DO THIS - safe from SQL injection
userInput := "1; DROP TABLE users; --"
query := "SELECT * FROM todos WHERE user_id = :user_id"
params := map[string]any{"user_id": userInput}
db.NamedQuery(query, params)  // userInput is treated as data, not SQL
```

**Why parameterized queries work:**
1. **Separation:** SQL structure is defined separately from data values
2. **Escaping:** Database driver automatically escapes dangerous characters
3. **Type safety:** Parameters are validated against expected types

**Additional security measures:**
```go
// Input validation
if len(userInput) > 255 {
    return errors.New("input too long")
}

// Whitelist validation for non-parameterizable parts (like column names)
allowedColumns := map[string]bool{"id": true, "title": true, "created_at": true}
if !allowedColumns[sortColumn] {
    return errors.New("invalid sort column")
}
```

---

**Q: Explain Go's error handling philosophy and how it applies to database operations.**

**A:** Go uses explicit error handling instead of exceptions, which is particularly valuable for database operations:

**Go's Error Philosophy:**
```go
// Go approach - explicit and predictable
result, err := db.Query("SELECT * FROM todos")
if err != nil {
    // Handle error immediately and explicitly
    return fmt.Errorf("failed to query todos: %w", err)
}
// Continue with result...
```

**Compared to exception-based languages:**
```java
// Java/C# approach - exceptions can be hidden
try {
    Result result = db.query("SELECT * FROM todos");
    // continue...
} catch (SQLException e) {
    // Handle error
}
```

**Benefits for database code:**
1. **Predictable control flow** - No hidden jumps from exceptions
2. **Forced error handling** - Compiler requires you to handle every error
3. **Local error context** - Handle errors where they occur
4. **Error wrapping** - Add context while preserving original error

**Database-specific error handling patterns:**
```go
func (s *Store) GetTodo(ctx context.Context, id int) (domain.Todo, error) {
    var todo domain.Todo
    err := s.db.GetContext(ctx, &todo, query, id)

    // Handle specific database errors
    if err == sql.ErrNoRows {
        return domain.Todo{}, ErrTodoNotFound  // Convert to domain error
    }

    if err != nil {
        // Wrap with context but preserve original error
        return domain.Todo{}, fmt.Errorf("failed to get todo %d: %w", id, err)
    }

    return todo, nil
}
```

---

### 11.2 Architecture and Design

**Q: How do you structure database access in a Go application? Explain the Repository pattern.**

**A:** The Repository pattern separates data access logic from business logic. Here's how to implement it properly:

**1. Define Domain Types (Independent of Database):**
```go
// domain/todo.go - Pure business logic, no database tags
package domain

type Todo struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Done      bool      `json:"done"`
    UserID    int       `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}
```

**2. Define Repository Interface (In Domain Layer):**
```go
// domain/repository.go - Abstract interface
package domain

type TodoRepository interface {
    Create(ctx context.Context, todo Todo) (int, error)
    GetByID(ctx context.Context, id int, userID int) (Todo, error)
    List(ctx context.Context, userID int, limit int) ([]Todo, error)
    Update(ctx context.Context, todo Todo) error
    Delete(ctx context.Context, id int, userID int) error
}
```

**3. Implement Repository (Infrastructure Layer):**
```go
// infrastructure/postgres/todo_store.go
package postgres

// DTO for database mapping
type TodoDTO struct {
    ID        int       `db:"id"`
    Title     string    `db:"title"`
    Done      bool      `db:"done"`
    UserID    int       `db:"user_id"`
    CreatedAt time.Time `db:"created_at"`
}

type TodoStore struct {
    db *sqlx.DB
}

// Ensure we implement the interface
var _ domain.TodoRepository = (*TodoStore)(nil)

func (s *TodoStore) Create(ctx context.Context, todo domain.Todo) (int, error) {
    query := `INSERT INTO todos (title, user_id, done)
              VALUES (:title, :user_id, :done)
              RETURNING id`

    params := map[string]any{
        "title":   todo.Title,
        "user_id": todo.UserID,
        "done":    todo.Done,
    }

    var id int
    err := s.db.GetContext(ctx, &id, query, params)
    return id, err
}

func (s *TodoStore) GetByID(ctx context.Context, id int, userID int) (domain.Todo, error) {
    var dto TodoDTO
    query := `SELECT id, title, done, user_id, created_at
              FROM todos
              WHERE id = :id AND user_id = :user_id`

    err := s.db.GetContext(ctx, &dto, query, map[string]any{
        "id": id,
        "user_id": userID,
    })

    if err == sql.ErrNoRows {
        return domain.Todo{}, domain.ErrTodoNotFound
    }
    if err != nil {
        return domain.Todo{}, err
    }

    // Convert DTO to domain object
    return domain.Todo{
        ID:        dto.ID,
        Title:     dto.Title,
        Done:      dto.Done,
        UserID:    dto.UserID,
        CreatedAt: dto.CreatedAt,
    }, nil
}
```

**4. Use Repository in Service Layer:**
```go
// service/todo_service.go
package service

type TodoService struct {
    repo domain.TodoRepository
}

func (s *TodoService) CreateTodo(ctx context.Context, userID int, title string) (domain.Todo, error) {
    // Business logic validation
    if strings.TrimSpace(title) == "" {
        return domain.Todo{}, errors.New("title cannot be empty")
    }

    todo := domain.Todo{
        Title:  strings.TrimSpace(title),
        UserID: userID,
        Done:   false,
    }

    id, err := s.repo.Create(ctx, todo)
    if err != nil {
        return domain.Todo{}, fmt.Errorf("failed to create todo: %w", err)
    }

    todo.ID = id
    return todo, nil
}
```

**Benefits of this structure:**
- ✅ **Testable** - Can mock repository interface
- ✅ **Database agnostic** - Can switch from PostgreSQL to MongoDB
- ✅ **Separation of concerns** - Business logic separate from data access
- ✅ **Clean dependencies** - Domain doesn't depend on infrastructure

---

**Q: How do you handle database configuration and connection management in production?**

**A:** Production database configuration requires careful attention to security, reliability, and performance:

**1. Environment-Based Configuration:**
```go
type DatabaseConfig struct {
    Host            string `env:"DB_HOST" envDefault:"localhost"`
    Port            int    `env:"DB_PORT" envDefault:"5432"`
    User            string `env:"DB_USER" envDefault:"postgres"`
    Password        string `env:"DB_PASSWORD"`
    Database        string `env:"DB_NAME" envDefault:"myapp"`
    SSLMode         string `env:"DB_SSLMODE" envDefault:"require"`
    MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
    MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
    ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME_MINUTES" envDefault:"5"`
}

func (c DatabaseConfig) DSN() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}
```

**2. Connection Pool Configuration:**
```go
func NewDatabase(config DatabaseConfig) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", config.DSN())
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Configure connection pool for production
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Minute)

    // Test the connection
    if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return db, nil
}
```

**3. Health Checks:**
```go
func (s *Store) HealthCheck(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    err := s.db.PingContext(ctx)
    if err != nil {
        return fmt.Errorf("database health check failed: %w", err)
    }

    return nil
}

// HTTP health endpoint
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    err := h.store.HealthCheck(r.Context())
    if err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "error": err.Error()})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
```

**4. Graceful Shutdown:**
```go
func main() {
    db, err := NewDatabase(config)
    if err != nil {
        log.Fatal(err)
    }

    // Setup graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-c
        log.Println("Shutting down gracefully...")

        // Close database connections
        if err := db.Close(); err != nil {
            log.Printf("Error closing database: %v", err)
        }

        os.Exit(0)
    }()

    // Start your application...
}
```

---

### 11.3 Performance and Optimization

**Q: How would you optimize a slow database query in a Go application?**

**A:** Query optimization is a systematic process. Here's my approach:

**1. Identify the Problem:**
```go
func (s *Store) SlowQuery(ctx context.Context, userID int) ([]domain.Todo, error) {
    // Add query logging to identify slow queries
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        if duration > 100*time.Millisecond {
            log.Printf("SLOW QUERY: %v - took %v", "SlowQuery", duration)
        }
    }()

    query := `SELECT t.*, u.name as user_name
              FROM todos t
              JOIN users u ON t.user_id = u.id
              WHERE t.user_id = :user_id
              AND t.created_at > NOW() - INTERVAL '30 days'
              ORDER BY t.created_at DESC`

    var todos []domain.Todo
    err := s.db.SelectContext(ctx, &todos, query, map[string]any{"user_id": userID})
    return todos, err
}
```

**2. Analyze with EXPLAIN:**
```go
func (s *Store) ExplainQuery(query string, params map[string]any) {
    explainQuery := "EXPLAIN ANALYZE " + query

    rows, err := s.db.NamedQuery(explainQuery, params)
    if err != nil {
        log.Printf("Error explaining query: %v", err)
        return
    }
    defer rows.Close()

    log.Printf("Query plan for: %s", query)
    for rows.Next() {
        var line string
        rows.Scan(&line)
        log.Printf("  %s", line)
    }
}
```

**3. Add Appropriate Indexes:**
```sql
-- For the slow query above, we need indexes on:
CREATE INDEX idx_todos_user_id_created_at ON todos(user_id, created_at DESC);
CREATE INDEX idx_users_id ON users(id);  -- Usually exists as PRIMARY KEY

-- For text search queries:
CREATE INDEX idx_todos_title_search ON todos USING gin(to_tsvector('english', title));
```

**4. Optimize the Query Structure:**
```go
// Before: N+1 query problem
func (s *Store) GetTodosWithUsersSlow(ctx context.Context, userID int) ([]TodoWithUser, error) {
    todos, err := s.GetTodos(ctx, userID)
    if err != nil {
        return nil, err
    }

    var result []TodoWithUser
    for _, todo := range todos {
        user, err := s.GetUser(ctx, todo.UserID)  // Separate query for each todo!
        if err != nil {
            return nil, err
        }
        result = append(result, TodoWithUser{Todo: todo, User: user})
    }
    return result, nil
}

// After: Single query with JOIN
func (s *Store) GetTodosWithUsersFast(ctx context.Context, userID int) ([]TodoWithUser, error) {
    query := `SELECT
                t.id, t.title, t.done, t.created_at,
                u.id as user_id, u.name as user_name, u.email as user_email
              FROM todos t
              JOIN users u ON t.user_id = u.id
              WHERE t.user_id = :user_id
              ORDER BY t.created_at DESC`

    rows, err := s.db.NamedQueryContext(ctx, query, map[string]any{"user_id": userID})
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []TodoWithUser
    for rows.Next() {
        var result TodoWithUser
        err := rows.Scan(
            &result.Todo.ID, &result.Todo.Title, &result.Todo.Done, &result.Todo.CreatedAt,
            &result.User.ID, &result.User.Name, &result.User.Email,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, result)
    }

    return results, nil
}
```

**5. Add Limits and Pagination:**
```go
func (s *Store) GetTodosPaginated(ctx context.Context, userID int, page, pageSize int) ([]domain.Todo, error) {
    if pageSize > 100 {
        pageSize = 100  // Prevent excessive memory usage
    }

    offset := (page - 1) * pageSize

    query := `SELECT id, title, done, created_at
              FROM todos
              WHERE user_id = :user_id
              ORDER BY created_at DESC
              LIMIT :limit OFFSET :offset`

    params := map[string]any{
        "user_id": userID,
        "limit":   pageSize,
        "offset":  offset,
    }

    var todos []domain.Todo
    err := s.db.SelectContext(ctx, &todos, query, params)
    return todos, err
}
```

**6. Monitor Performance:**
```go
func (s *Store) WithMetrics(operation string) func() {
    start := time.Now()
    return func() {
        duration := time.Since(start)

        // Log slow queries
        if duration > 100*time.Millisecond {
            log.Printf("SLOW: %s took %v", operation, duration)
        }

        // Send metrics to monitoring system (Prometheus, etc.)
        dbOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
    }
}

func (s *Store) GetTodo(ctx context.Context, id int) (domain.Todo, error) {
    defer s.WithMetrics("GetTodo")()

    // ... implementation
}
```

---

**Q: How do you handle database connections and connection pooling in a high-traffic Go application?**

**A:** Connection pooling is critical for performance in high-traffic applications:

**1. Understanding Connection Pool Mechanics:**
```go
func OptimizeConnectionPool(db *sqlx.DB) {
    // Maximum concurrent connections to PostgreSQL
    // Should be less than PostgreSQL's max_connections setting
    db.SetMaxOpenConns(100)  // Adjust based on your PostgreSQL config

    // Keep connections warm to avoid connection overhead
    // But not too many (uses memory and PostgreSQL resources)
    db.SetMaxIdleConns(25)

    // Rotate connections to avoid long-lived connection issues
    db.SetConnMaxLifetime(5 * time.Minute)

    // Close idle connections after this time
    db.SetConnMaxIdleTime(10 * time.Minute)
}
```

**2. Monitor Connection Pool Usage:**
```go
func MonitorConnectionPool(db *sqlx.DB) {
    ticker := time.NewTicker(30 * time.Second)
    go func() {
        for range ticker.C {
            stats := db.Stats()

            log.Printf("DB Pool Stats:")
            log.Printf("  Open: %d/%d", stats.OpenConnections, stats.MaxOpenConnections)
            log.Printf("  In Use: %d", stats.InUse)
            log.Printf("  Idle: %d", stats.Idle)

            // Alert if we're frequently waiting for connections
            if stats.WaitCount > 0 {
                log.Printf("  ⚠️  Wait Events: %d (total wait time: %v)",
                          stats.WaitCount, stats.WaitDuration)
            }

            // Alert if too many connections are being closed
            if stats.MaxLifetimeClosed > 10 {
                log.Printf("  ⚠️  Lifetime closures: %d", stats.MaxLifetimeClosed)
            }
        }
    }()
}
```

**3. Handle Connection Errors Gracefully:**
```go
func (s *Store) GetTodoWithRetry(ctx context.Context, id int) (domain.Todo, error) {
    var todo domain.Todo
    var lastErr error

    // Retry logic for transient connection errors
    for attempt := 1; attempt <= 3; attempt++ {
        err := s.db.GetContext(ctx, &todo, "SELECT * FROM todos WHERE id = :id",
                              map[string]any{"id": id})

        if err == nil {
            return todo, nil
        }

        lastErr = err

        // Check if it's a connection error worth retrying
        if isRetryableError(err) && attempt < 3 {
            backoff := time.Duration(attempt) * 100 * time.Millisecond
            time.Sleep(backoff)
            continue
        }

        break
    }

    return domain.Todo{}, lastErr
}

func isRetryableError(err error) bool {
    errStr := err.Error()
    return strings.Contains(errStr, "connection refused") ||
           strings.Contains(errStr, "timeout") ||
           strings.Contains(errStr, "temporary failure")
}
```

**4. Use Read Replicas for High Traffic:**
```go
type Store struct {
    writeDB *sqlx.DB  // Master database for writes
    readDB  *sqlx.DB  // Read replica for queries
}

func (s *Store) GetTodo(ctx context.Context, id int) (domain.Todo, error) {
    // Use read replica for queries
    var todo domain.Todo
    err := s.readDB.GetContext(ctx, &todo, "SELECT * FROM todos WHERE id = :id",
                              map[string]any{"id": id})
    return todo, err
}

func (s *Store) CreateTodo(ctx context.Context, todo domain.Todo) (int, error) {
    // Use master database for writes
    query := "INSERT INTO todos (title, user_id) VALUES (:title, :user_id) RETURNING id"
    var id int
    err := s.writeDB.GetContext(ctx, &id, query, map[string]any{
        "title":   todo.Title,
        "user_id": todo.UserID,
    })
    return id, err
}
```

**5. Circuit Breaker Pattern for Database Failures:**
```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFailure time.Time
    state       string // "closed", "open", "half-open"
    mutex       sync.Mutex
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = "half-open"
        } else {
            return errors.New("circuit breaker is open")
        }
    }

    err := operation()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()

        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }

    // Success - reset circuit breaker
    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

---

### 11.4 Real-World Scenarios

**Q: How would you implement pagination efficiently in PostgreSQL with Go?**

**A:** Efficient pagination is crucial for performance. Here are the main approaches:

**1. Offset-Based Pagination (Simple but has limitations):**
```go
type PaginationRequest struct {
    Page     int `json:"page" query:"page"`
    PageSize int `json:"page_size" query:"page_size"`
}

func (r *PaginationRequest) Validate() error {
    if r.Page < 1 {
        r.Page = 1
    }
    if r.PageSize < 1 || r.PageSize > 100 {
        r.PageSize = 20  // Default page size
    }
    return nil
}

func (s *Store) ListTodosWithPagination(ctx context.Context, userID int, req PaginationRequest) ([]domain.Todo, error) {
    req.Validate()

    offset := (req.Page - 1) * req.PageSize

    query := `SELECT id, title, done, created_at
              FROM todos
              WHERE user_id = :user_id
              ORDER BY created_at DESC
              LIMIT :limit OFFSET :offset`

    params := map[string]any{
        "user_id": userID,
        "limit":   req.PageSize,
        "offset":  offset,
    }

    var todos []domain.Todo
    err := s.db.SelectContext(ctx, &todos, query, params)
    return todos, err
}
```

**Problem with OFFSET:** As the offset increases, PostgreSQL still has to count all the skipped rows, making later pages very slow.

**2. Cursor-Based Pagination (More Efficient):**
```go
type CursorPaginationRequest struct {
    Cursor   *time.Time `json:"cursor,omitempty"`  // Last seen created_at
    PageSize int        `json:"page_size"`
}

func (s *Store) ListTodosWithCursor(ctx context.Context, userID int, req CursorPaginationRequest) ([]domain.Todo, *time.Time, error) {
    if req.PageSize < 1 || req.PageSize > 100 {
        req.PageSize = 20
    }

    var query string
    params := map[string]any{
        "user_id": userID,
        "limit":   req.PageSize,
    }

    if req.Cursor != nil {
        // Get todos created before the cursor
        query = `SELECT id, title, done, created_at
                 FROM todos
                 WHERE user_id = :user_id AND created_at < :cursor
                 ORDER BY created_at DESC
                 LIMIT :limit`
        params["cursor"] = *req.Cursor
    } else {
        // First page
        query = `SELECT id, title, done, created_at
                 FROM todos
                 WHERE user_id = :user_id
                 ORDER BY created_at DESC
                 LIMIT :limit`
    }

    var todos []domain.Todo
    err := s.db.SelectContext(ctx, &todos, query, params)
    if err != nil {
        return nil, nil, err
    }

    // Return next cursor (created_at of last item)
    var nextCursor *time.Time
    if len(todos) > 0 {
        lastCreatedAt := todos[len(todos)-1].CreatedAt
        nextCursor = &lastCreatedAt
    }

    return todos, nextCursor, nil
}
```

**3. Keyset Pagination (Most Efficient for Large Datasets):**
```go
type KeysetPaginationRequest struct {
    LastID        *int       `json:"last_id,omitempty"`
    LastCreatedAt *time.Time `json:"last_created_at,omitempty"`
    PageSize      int        `json:"page_size"`
}

func (s *Store) ListTodosWithKeyset(ctx context.Context, userID int, req KeysetPaginationRequest) ([]domain.Todo, error) {
    if req.PageSize < 1 || req.PageSize > 100 {
        req.PageSize = 20
    }

    var query string
    params := map[string]any{
        "user_id": userID,
        "limit":   req.PageSize,
    }

    if req.LastID != nil && req.LastCreatedAt != nil {
        // Continue from last seen record
        query = `SELECT id, title, done, created_at
                 FROM todos
                 WHERE user_id = :user_id
                 AND (created_at < :last_created_at
                      OR (created_at = :last_created_at
