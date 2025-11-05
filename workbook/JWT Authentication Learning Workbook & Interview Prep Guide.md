# JWT Authentication Learning Workbook & Interview Prep Guide

**Complete Guide to JWT Authentication in Go with Clean Architecture**

> 📚 **Purpose:** Master JWT authentication implementation in Go, understand clean architecture patterns, and prepare for technical interviews.

---

## 📑 Table of Contents

- [1. JWT Fundamentals](#1-jwt-fundamentals)
  - [1.1 What is JWT?](#11-what-is-jwt)
  - [1.2 JWT Structure Deep Dive](#12-jwt-structure-deep-dive)
  - [1.3 Standard JWT Claims](#13-standard-jwt-claims)
  - [1.4 JWT vs Session Authentication](#14-jwt-vs-session-authentication)
  - [1.5 How JWT Authentication Works](#15-how-jwt-authentication-works)
- [2. Clean Architecture Layers](#2-clean-architecture-layers)
  - [2.1 Layer Overview](#21-layer-overview)
  - [2.2 Responsibility Matrix](#22-responsibility-matrix)
  - [2.3 Data Flow - Login](#23-data-flow---login)
  - [2.4 Data Flow - Protected Routes](#24-data-flow---protected-routes)
- [3. Password Hashing with bcrypt](#3-password-hashing-with-bcrypt)
  - [3.1 Why Hash Passwords?](#31-why-hash-passwords)
  - [3.2 bcrypt Overview](#32-bcrypt-overview)
  - [3.3 Implementation Guide](#33-implementation-guide)
  - [3.4 Common Mistakes](#34-common-mistakes)
  - [3.5 Best Practices](#35-best-practices)
- [4. Complete Implementation Guide](#4-complete-implementation-guide)
  - [4.1 Domain Layer](#41-domain-layer)
  - [4.2 Service Layer](#42-service-layer)
  - [4.3 Delivery Layer - Auth Infrastructure](#43-delivery-layer---auth-infrastructure)
  - [4.4 Delivery Layer - Middleware](#44-delivery-layer---middleware)
  - [4.5 Delivery Layer - Handlers](#45-delivery-layer---handlers)
  - [4.6 Router Setup](#46-router-setup)
- [5. Security Best Practices](#5-security-best-practices)
  - [5.1 Secret Key Management](#51-secret-key-management)
  - [5.2 Token Expiration](#52-token-expiration)
  - [5.3 HTTPS in Production](#53-https-in-production)
  - [5.4 Token Storage](#54-token-storage)
  - [5.5 What NOT to Store in JWT](#55-what-not-to-store-in-jwt)
  - [5.6 Token Revocation Strategies](#56-token-revocation-strategies)
- [6. Common Pitfalls & Solutions](#6-common-pitfalls--solutions)
- [7. Testing JWT Authentication](#7-testing-jwt-authentication)
  - [7.1 Unit Testing Handlers](#71-unit-testing-handlers)
  - [7.2 Testing Middleware](#72-testing-middleware)
  - [7.3 Integration Tests](#73-integration-tests)
- [8. Interview Questions & Answers](#8-interview-questions--answers)
  - [8.1 Basic Questions](#81-basic-questions)
  - [8.2 Intermediate Questions](#82-intermediate-questions)
  - [8.3 Advanced Questions](#83-advanced-questions)
- [9. Practical Exercises](#9-practical-exercises)
- [10. Quick Reference Cheat Sheet](#10-quick-reference-cheat-sheet)

---

## 1. JWT Fundamentals

### 1.1 What is JWT?

**JWT (JSON Web Token)** is a compact, URL-safe means of representing claims to be transferred between two parties as a JSON object.

**Key Characteristics:**
- ✅ **Self-contained** - Contains all user information needed
- ✅ **Stateless** - Server doesn't store session data
- ✅ **Signed** - Can verify authenticity using cryptographic signature
- ⚠️ **Not encrypted** - Anyone can decode and read the payload (Base64)
- ✅ **Portable** - Works across different domains and services

**When to use JWT:**
- Stateless REST APIs
- Microservices architecture
- Mobile applications
- Single Page Applications (SPAs)
- Cross-domain authentication

---

### 1.2 JWT Structure Deep Dive

A JWT consists of three Base64Url-encoded parts separated by dots (`.`):


**Example JWT:**
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImVtYWlsIjoidXNlckBleGFtcGxlLmNvbSIsImV4cCI6MTczNTY4OTYwMH0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

---

#### **Part 1: Header (Red)**

Contains metadata about the token:

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

alg - Signing algorithm (HS256, RS256, etc.)
typ - Token type (always "JWT")

Base64Url encoded:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9


Part 2: Payload (Purple)
Contains the claims (user data):

{
  "user_id": 123,
  "email": "user@example.com",
  "name": "John Doe",
  "exp": 1735689600,
  "iat": 1735603200
}

Base64Url encoded:

eyJ1c2VyX2lkIjoxMjMsImVtYWlsIjoidXNlckBleGFtcGxlLmNvbSIsImV4cCI6MTczNTY4OTYwMH0
Part 3: Signature (Blue)
Ensures the token hasn't been tampered with:

scss
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret
)
Result:

SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
🔐 Important: Only the server with the secret key can create valid signatures!


| Claim | Full Name | Type | Description | Example |
|-------|-----------|------|-------------|---------|
| `iss` | Issuer | String | Who created the token | `"auth.example.com"` |
| `sub` | Subject | String | Subject identifier | `"user123"` |
| `aud` | Audience | String | Intended recipient | `"api.example.com"` |
| `exp` | Expiration Time | Number | Unix timestamp | `1735689600` |
| `nbf` | Not Before | Number | Token valid after this | `1735603200` |
| `iat` | Issued At | Number | Creation time | `1735603200` |
| `jti` | JWT ID | String | Unique identifier | `"abc123"` |


Custom Claims (Application-specific)

{
  "user_id": 123,
  "email": "user@example.com",
  "name": "John Doe",
  "role": "admin",
  "permissions": ["read", "write", "delete"]
}
⚠️ Important: Custom claims should NOT contain sensitive data (passwords, credit cards, etc.)

1.4 JWT vs Session Authentication
| Feature | JWT | Session |
|---------|-----|---------|
| **Storage Location** | Client (localStorage, cookie) | Server (database, Redis) |
| **Stateless** | ✅ Yes | ❌ No |
| **Scalability** | ✅ Easy (no shared state) | ⚠️ Harder (shared session store) |
| **Server Memory** | ✅ None | ⚠️ Required |
| **Revocation** | ❌ Difficult | ✅ Easy (delete session) |
| **Size** | ⚠️ Larger (sent every request) | ✅ Small (session ID only) |
| **Cross-domain** | ✅ Easy | ⚠️ Requires configuration |
| **Security** | ⚠️ XSS risk (localStorage) | ✅ Better (HttpOnly cookie) |


### When to use JWT:
Microservices architecture
Mobile applications
Stateless REST APIs
Need horizontal scaling
Cross-domain authentication
When to use Sessions:
Traditional web applications
Need immediate revocation
Smaller data transfer requirements
Single server or shared session store available

### 1.5 How JWT Authentication Works
┌─────────┐                                    ┌─────────┐
│ Client  │                                    │ Server  │
└────┬────┘                                    └────┬────┘
     │                                              │
     │  1. POST /login                             │
     │     Body: {                                 │
     │       "email": "user@example.com",          │
     │       "password": "Password123"             │
     │     }                                       │
     │─────────────────────────────────────────────>│
     │                                              │
     │                       2. Verify Credentials │
     │                          - Get user from DB │
     │                          - Compare password │
     │                            hash with bcrypt │
     │                                              │
     │                          3. Generate JWT    │
     │                          - Create claims    │
     │                          - Sign with secret │
     │                                              │
     │  4. Response:                               │
     │     {                                       │
     │       "token": "eyJhbG...",                 │
     │       "user": {...}                         │
     │     }                                       │
     │<─────────────────────────────────────────────│
     │                                              │
     │  5. Store Token                             │
     │     localStorage.setItem('token', ...)      │
     │                                              │
     │  6. GET /protected-resource                 │
     │     Headers:                                │
     │       Authorization: Bearer eyJhbG...       │
     │─────────────────────────────────────────────>│
     │                                              │
     │                          7. Validate Token  │
     │                          - Extract from hdr │
     │                          - Verify signature │
     │                          - Check expiration │
     │                          - Extract user info│
     │                                              │
     │  8. Response: { data: [...] }               │
     │<─────────────────────────────────────────────│
     │


     2. Clean Architecture Layers
     2.1 Layer Overview
     Clean Architecture separates concerns into distinct layers:
     ┌───────────────────────────────────────────────────────────────┐
     │                      DELIVERY LAYER                           │
     │  (HTTP Handlers, Middleware, Request/Response DTOs)           │
     │                                                                │
     │  Responsibilities:                                             │
     │  • Handle HTTP requests/responses                             │
     │  • Validate input (JSON parsing, validation tags)             │
     │  • Generate JWT tokens                                        │
     │  • Extract JWT from headers                                   │
     │  • Return appropriate HTTP status codes                       │
     │                                                                │
     │  Files: handlers/, middleware/, auth/                         │
     └───────────────────────┬───────────────────────────────────────┘
                             │
                             ▼
     ┌───────────────────────────────────────────────────────────────┐
     │                       SERVICE LAYER                           │
     │  (Business Logic, Use Cases)                                  │
     │                                                                │
     │  Responsibilities:                                             │
     │  • User registration logic                                    │
     │  • Password verification (bcrypt)                             │
     │  • Business rules validation                                  │
     │  • Orchestrate repository calls                              │
     │  • NO JWT knowledge (infrastructure concern)                  │
     │  • NO HTTP knowledge (delivery concern)                       │
     │                                                                │
     │  Files: service/                                              │
     └───────────────────────┬───────────────────────────────────────┘
                             │
                             ▼
     ┌───────────────────────────────────────────────────────────────┐
     │                       DOMAIN LAYER                            │
     │  (Entities, Interfaces, Business Errors)                      │
     │                                                                │
     │  Responsibilities:                                             │
     │  • Define core entities (User, Todo)                          │
     │  • Define service interfaces                                  │
     │  • Define domain errors                                       │
     │  • Business validation rules                                  │
     │  • NO external dependencies                                   │
     │                                                                │
     │  Files: domain/                                               │
     └───────────────────────┬───────────────────────────────────────┘
                             │
                             ▼
     ┌───────────────────────────────────────────────────────────────┐
     │                     REPOSITORY LAYER                          │
     │  (Database Access, Data Persistence)                          │
     │                                                                │
     │  Responsibilities:                                             │
     │  • CRUD operations                                            │
     │  • SQL queries                                                │
     │  • Data mapping (DB ↔ Domain)                                │
     │  • Transaction management                                     │
     │                                                                │
     │  Files: repository/                                           │
     └───────────────────────────────────────────────────────────────┘

     2.2 Responsibility Matrix

     | Layer | JWT | Password | HTTP | Database | Business Logic |
     |-------|-----|----------|------|----------|----------------|
     | **Domain** | ❌ No | ❌ No | ❌ No | ❌ No | ✅ Yes (rules) |
     | **Service** | ❌ No | ✅ Yes (hash/verify) | ❌ No | ❌ No | ✅ Yes (logic) |
     | **Delivery** | ✅ Yes (generate/validate) | ❌ No | ✅ Yes | ❌ No | ❌ No |
     | **Repository** | ❌ No | ❌ No | ❌ No | ✅ Yes | ❌ No |

Key Principle:

JWT is a delivery mechanism, not business logic!

Password hashing is business logic, not delivery concern!

2.3 Data Flow - Login

// ────────────────────────────────────────────────────────────────
// 1. CLIENT REQUEST
// ────────────────────────────────────────────────────────────────
POST /login HTTP/1.1
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password123"
}

// ────────────────────────────────────────────────────────────────
// 2. DELIVERY LAYER (Handler)
// ────────────────────────────────────────────────────────────────
```go
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
    // Parse and validate request
    var req LoginDTO
    json.NewDecoder(r.Body).Decode(&req)
    validate.Struct(req) // Email format, required fields

    // Call service (business logic)
    user, err := h.Service.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        // Handle errors (401, 500, etc.)
        return
    }
```

    // ✅ GENERATE JWT (delivery concern - here!)
    claims := auth.NewUserClaims(user, 24*time.Hour)
    _, tokenString, _ := h.TokenAuth.Encode(claims.ToMap())

    // Return response
    utils.WriteJSON(w, http.StatusOK, LoginResponseDTO{
        Token: tokenString,
        User:  UserResponseDTO{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
        },
    })
}

// ────────────────────────────────────────────────────────────────
// 3. SERVICE LAYER (Business Logic)
// ────────────────────────────────────────────────────────────────
func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
    // Fetch user from database
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, domain.ErrInvalidCredentials
    }

    // ✅ VERIFY PASSWORD (business logic - here!)
    err = bcrypt.CompareHashAndPassword(
        []byte(user.Password), // Hashed password from DB
        []byte(password),      // Plain password from user
    )
    if err != nil {
        return nil, domain.ErrInvalidCredentials
    }

    // Return authenticated user (NO JWT knowledge!)
    return user, nil
}
```go
// ────────────────────────────────────────────────────────────────
// 4. REPOSITORY LAYER (Database)
// ────────────────────────────────────────────────────────────────
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User

    query := "SELECT id, name, email, password FROM users WHERE email = ?"
    err := r.DB.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password, // Hashed password
    )

    if err != nil {
        return nil, err
    }

    return &user, nil
}

2.4 Data Flow - Protected Routes

// ────────────────────────────────────────────────────────────────
// 1. CLIENT REQUEST
// ────────────────────────────────────────────────────────────────
GET /todos HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

// ────────────────────────────────────────────────────────────────
// 2. MIDDLEWARE CHAIN
// ────────────────────────────────────────────────────────────────

// 2a. Verifier (jwtauth library) - Extracts and validates token
func Verifier(tokenAuth *JWTAuth) func(http.Handler) http.Handler {
    // • Extract token from Authorization header
    // • Verify signature with secret key
    // • Validate expiration (exp claim)
    // • Store token in request context
    // • Continue even if invalid (stores error in context)
}

// 2b. Authenticator (custom) - Checks validation result
func Authenticator(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, _, err := jwtauth.FromContext(r.Context())

        if err != nil || token == nil {
            http.Error(w, "Unauthorized", 401)
            return // Stop here!
        }

        // Validate claims structure
        claims, err := auth.FromTokenClaims(token.PrivateClaims())
        if err != nil {
            http.Error(w, "Invalid token claims", 401)
            return // Stop here!
        }

        next.ServeHTTP(w, r) // Continue
    })
}

// 2c. UserContext (custom) - Extracts user info
func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, _, _ := jwtauth.FromContext(r.Context())

        // Extract claims
        claims, _ := auth.FromTokenClaims(token.PrivateClaims())

        // Create UserContext
        userCtx := auth.NewUserContext(claims)

        // Add to request context
        ctx := userCtx.AddToContext(r.Context())

        // Pass updated context to handler
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// ────────────────────────────────────────────────────────────────
// 3. HANDLER (uses authenticated user)
// ────────────────────────────────────────────────────────────────
func (h *TodoHandlers) ListTodos(w http.ResponseWriter, r *http.Request) {
    // Get user from context (set by middleware)
    userCtx, ok := auth.FromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", 401)
        return
    }

    // Call service with user ID
    todos, err := h.Service.ListTodos(r.Context(), userCtx.ID)

    utils.WriteJSON(w, http.StatusOK, todos)
}

// ────────────────────────────────────────────────────────────────
// 4. SERVICE LAYER
// ────────────────────────────────────────────────────────────────
func (s *TodoService) ListTodos(ctx context.Context, userID int64) ([]*domain.Todo, error) {
    // Business rule: user can only see their own todos
    return s.Repo.GetByUserID(ctx, userID)
}

3. Password Hashing with bcrypt
3.1 Why Hash Passwords?

❌ NEVER store passwords in plain text!
// ❌ TERRIBLE - Anyone with DB access sees passwords
user := User{
    Email:    "user@example.com",
    Password: "MyPassword123", // Plain text!
}
db.Insert(user)


Risks of plain text passwords:

🔓 Database breach exposes all passwords
🔓 Admins/DBAs can see passwords
🔓 Users reuse passwords across sites
🔓 Password hints/reset becomes dangerous
✅ Always hash passwords before storing:


// ✅ GOOD - Password is hashed
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("MyPassword123"), bcrypt.DefaultCost)
user := User{
    Email:    "user@example.com",
    Password: string(hashedPassword), // "$2a$10$N9qo8uLOick..."
}
db.Insert(user)

3.2 bcrypt Overview
bcrypt is a password-hashing function designed to be slow and resistant to brute-force attacks.

Key Features:

✅ Salted - Random salt added to each password (prevents rainbow tables)
✅ Adaptive - Cost factor can increase as computers get faster
✅ Slow by design - Makes brute force attacks impractical
✅ One-way - Cannot decrypt back to original password
How bcrypt works:

Password: "MyPassword123"
    ↓
bcrypt.GenerateFromPassword(password, cost)
    ↓
1. Generate random salt (16 bytes)
2. Combine salt + password
3. Apply bcrypt algorithm (2^cost iterations)
    ↓
Hashed: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
         └─┬─┘└┬┘└────────┬────────┘└──────────────┬──────────────┘
          alg cost   salt (22 chars)    hash (31 chars)

Parts of a bcrypt hash:

$2a$ - bcrypt algorithm version
10 - Cost factor (2^10 = 1024 iterations)
N9qo8uLOickgx2ZMRZoMye - Salt (22 characters)
IjZAgcfl7p92ldGxad68LJZdL17lhWy - Hash (31 characters)

3.3 Implementation Guide
Registration - Hash Password

package service

import (
    "context"
    "github.com/macesz/todo-go/domain"
    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    Repo domain.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
    // Check if email already exists
    existing, _ := s.Repo.GetByEmail(ctx, email)
    if existing != nil {
        return nil, domain.ErrDuplicate
    }

    // ✅ Hash password with bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(password),      // Plain text password
        bcrypt.DefaultCost,    // Cost factor (10)
    )
    if err != nil {
        return nil, err
    }

    // Create user with hashed password
    user := &domain.User{
        Name:     name,
        Email:    email,
        Password: string(hashedPassword), // Store hashed password
    }

    // Save to database
    if err := s.Repo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

Cost Factors:
| Cost | Iterations | Time (approx) | Use Case |
|------|-----------|---------------|----------|
| 4 | 16 | ~0.001s | Testing only |
| 8 | 256 | ~0.02s | Very weak |
| 10 | 1,024 | ~0.1s | Default (recommended) |
| 12 | 4,096 | ~0.4s | High security |
| 14 | 16,384 | ~1.6s | Very high security |
Recommendation: Use bcrypt.DefaultCost (10) for most applications.


Login - Verify Password

func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
    // Fetch user from database
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        // Don't reveal if email exists or not
        return nil, domain.ErrInvalidCredentials
    }

    // ✅ Compare password with bcrypt
    err = bcrypt.CompareHashAndPassword(
        []byte(user.Password), // Hashed password from database
        []byte(password),      // Plain text password from login request
    )
    if err != nil {
        // Password doesn't match
        return nil, domain.ErrInvalidCredentials
    }

    // Password is correct!
    return user, nil
}


How CompareHashAndPassword works:

Input:
  - hashedPassword: "$2a$10$N9qo8uLO..."
  - plainPassword: "MyPassword123"

Process:
  1. Extract salt from hashedPassword
  2. Hash plainPassword with extracted salt and same cost
  3. Compare resulting hash with stored hash

Result:
  - If hashes match → return nil (success)
  - If hashes don't match → return error

Important Notes:

✅ CompareHashAndPassword is constant-time (prevents timing attacks)
✅ Always returns error for invalid password (don't check error message)
✅ Same password + same salt = same hash (deterministic)
✅ Different passwords = different hashes


3.4 Common Mistakes
❌ Mistake 1: Comparing hashes directly

// ❌ WRONG - Hashes are different every time!
if user.Password == password {
    // This will NEVER work!
}

// ❌ WRONG - Even hashing again doesn't work
hashedInput, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if user.Password == string(hashedInput) {
    // Still wrong! Salt is different each time
}

// ✅ CORRECT - Use CompareHashAndPassword
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
if err != nil {
    // Password is wrong
}

❌ Mistake 2: Hashing in handler layer

// ❌ WRONG - Business logic in delivery layer
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserDTO
    json.NewDecoder(r.Body).Decode(&req)

    // Hashing here is wrong layer!
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

    user, _ := h.Service.CreateUser(r.Context(), req.Name, req.Email, string(hashedPassword))
    // ...
}

// ✅ CORRECT - Hashing in service layer
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*User, error) {
    // Hash here - it's business logic!
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    // ...
}

❌ Mistake 3: Revealing if email exists

// ❌ WRONG - Tells attacker if email is registered
func (s *UserService) Login(ctx context.Context, email, password string) (*User, error) {
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, errors.New("Email not found") // Reveals email doesn't exist
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, errors.New("Invalid password") // Reveals email exists
    }

    return user, nil
}

// ✅ CORRECT - Generic error message
func (s *UserService) Login(ctx context.Context, email, password string) (*User, error) {
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, domain.ErrInvalidCredentials // Generic message
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, domain.ErrInvalidCredentials // Same generic message
    }

    return user, nil
}

❌ Mistake 4: Not handling bcrypt errors
// ❌ WRONG - Ignoring errors
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// What if this fails?

// ✅ CORRECT - Handle errors
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
    return nil, fmt.Errorf("failed to hash password: %w", err)
}

3.5 Best Practices

// ✅ GOOD - Use default cost (10)
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    bcrypt.DefaultCost, // Cost = 10
)

// ⚠️ Only increase for high-security scenarios
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    12, // Higher cost = slower but more secure
)

✅ 2. Return generic error messages


// ✅ GOOD - Don't reveal details
//
```g
var ErrInvalidCredentials = errors.New("invalid email or password")

func (s *UserService) Login(ctx context.Context, email, password string) (*User, error) {
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, ErrInvalidCredentials // Generic
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, ErrInvalidCredentials // Same message
    }

    return user, nil
}
```

✅ 3. Validate password strength before hashing

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*User, error) {
    // Validate password strength BEFORE hashing
    if len(password) < 8 {
        return nil, domain.ErrWeakPassword
    }

    // Could also check for:
    // - Uppercase letters
    // - Numbers
    // - Special characters

    // Now hash
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    // ...
}

✅ 4. Test with known passwords

func TestLogin_Success(t *testing.T) {
    // Create test user with known password
    plainPassword := "Password123"
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

    testUser := &domain.User{
        ID:       1,
        Email:    "test@example.com",
        Password: string(hashedPassword),
    }

    // Mock repository
    mockRepo := new(mocks.MockUserRepository)
    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(testUser, nil)

    service := NewUserService(mockRepo)

    // Test login with correct password
    user, err := service.Login(context.Background(), "test@example.com", plainPassword)

    assert.NoError(t, err)
    assert.Equal(t, testUser.ID, user.ID)
}

✅ 5. Never log passwords


// ❌ WRONG - Logging password
log.Printf("Creating user with email %s and password %s", email, password)

// ❌ WRONG - Logging hashed password
log.Printf("Password hash: %s", hashedPassword)

// ✅ CORRECT - Don't log sensitive data
log.Printf("Creating user with email %s", email)

4. Complete Implementation Guide
4.1 Domain Layer
Purpose: Define core business entities and rules

domain/user.go

package domain

import (
    "context"
    "errors"
)

// User - Core business entity
type User struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"` // Hashed password, never in JSON response
}

// UserService - Business logic interface
type UserService interface {
    CreateUser(ctx context.Context, name, email, password string) (*User, error)
    Login(ctx context.Context, email, password string) (*User, error)
    GetUser(ctx context.Context, id int64) (*User, error)
}

// UserRepository - Data access interface
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int64) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
}

4.2 Service Layer
Purpose: Implement business logic
service/user_service.go
package service

import (
    "context"
    "fmt"

    "github.com/macesz/todo-go/domain"
    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    Repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
    return &UserService{Repo: repo}
}

// CreateUser - Register new user
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
    // Business rule: check if email already exists
    existing, err := s.Repo.GetByEmail(ctx, email)
    if err == nil && existing != nil {
        return nil, domain.ErrDuplicate
    }

    // Business rule: hash password before storing
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(password),
        bcrypt.DefaultCost, // Cost factor: 10
    )
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    // Create user entity
    user := &domain.User{
        Name:     name,
        Email:    email,
        Password: string(hashedPassword),
    }

    // Persist to database
    if err := s.Repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}

// Login - Authenticate user
```go
func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
    // Fetch user from database
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        // Don't reveal if email exists
        return nil, domain.ErrInvalidCredentials
    }

    // Business rule: verify password with bcrypt
    err = bcrypt.CompareHashAndPassword(
        []byte(user.Password), // Hashed password from DB
        []byte(password),      // Plain text password from request
    )
    if err != nil {
        // Password doesn't match
        return nil, domain.ErrInvalidCredentials
    }

    // Successfully authenticated
    return user, nil
}

// GetUser - Fetch user by ID
func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
    user, err := s.Repo.GetByID(ctx, id)
    if err != nil {
        return nil, domain.ErrUserNotFound
    }
    return user, nil
}
```
Key Points:

✅ Password hashing in service layer (business logic)
✅ Generic error messages (don't reveal email existence)
✅ No JWT knowledge
✅ No HTTP knowledge
✅ Returns domain entities

### 4.3 Delivery Layer - Auth Infrastructure
Purpose: JWT claims and context management

delivery/web/auth/claims.go

package auth

import (
    "errors"
    "time"

    "github.com/macesz/todo-go/domain"
)

// UserClaims - JWT payload structure
type UserClaims struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
    Exp    int64  `json:"exp"`
}

// NewUserClaims - Convert domain.User to JWT claims
func NewUserClaims(user *domain.User, expiresIn time.Duration) UserClaims {
    return UserClaims{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
        Exp:    time.Now().Add(expiresIn).Unix(),
    }
}

// ToMap - Convert to map for jwtauth library
func (c UserClaims) ToMap() map[string]interface{} {
    return map[string]interface{}{
        "user_id": c.UserID,
        "email":   c.Email,
        "name":    c.Name,
        "exp":     c.Exp,
    }
}

// FromTokenClaims - Extract and validate claims from JWT
// IMPORTANT: JWT stores numbers as float64, not int64!
```go
func FromTokenClaims(claims map[string]interface{}) (*UserClaims, error) {
    // Extract user_id (JWT stores as float64)
    userID, ok := claims["user_id"].(float64)
    if !ok {
        return nil, errors.New("missing or invalid user_id claim")
    }

    // Extract email
    email, ok := claims["email"].(string)
    if !ok {
        return nil, errors.New("missing or invalid email claim")
    }

    // Extract name
    name, ok := claims["name"].(string)
    if !ok {
        return nil, errors.New("missing or invalid name claim")
    }

    return &UserClaims{
        UserID: int64(userID), // Convert float64 to int64
        Email:  email,
        Name:   name,
    }, nil
}
```
delivery/web/auth/context.go

```go
package auth

import "context"

type contextKey string

const userContextKey contextKey = "userContext"

// UserContext - User information in request context
type UserContext struct {
    ID    int64
    Email string
    Name  string
}

// NewUserContext - Create from JWT claims
func NewUserContext(claims *UserClaims) *UserContext {
    return &UserContext{
        ID:    claims.UserID,
        Email: claims.Email,
        Name:  claims.Name,
    }
}

// AddToContext - Store in request context
func (u *UserContext) AddToContext(ctx context.Context) context.Context {
    return context.WithValue(ctx, userContextKey, u)
}

// FromContext - Retrieve from request context
func FromContext(ctx context.Context) (*UserContext, bool) {
    userCtx, ok := ctx.Value(userContextKey).(*UserContext)
    return userCtx, ok
}
```
### 4.4 Delivery Layer - Middleware
Purpose: Extract, validate, and process JWT tokens

delivery/web/middlewares/auth.go

```go
package middlewares

import (
    "net/http"

    "github.com/go-chi/jwtauth/v5"
    "github.com/macesz/todo-go/delivery/web/auth"
    "github.com/macesz/todo-go/domain"
)

// Authenticator - Validates token and claims structure
// Runs AFTER jwtauth.Verifier
func Authenticator(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get token from context (set by Verifier)
        token, _, err := jwtauth.FromContext(r.Context())

        // Check if Verifier found issues
        if err != nil {
            writeError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
            return
        }

        if token == nil {
            writeError(w, http.StatusUnauthorized, "Missing authentication token")
            return
        }

        // Validate claims structure
        privateClaims := token.PrivateClaims()
        _, err = auth.FromTokenClaims(privateClaims)
        if err != nil {
            writeError(w, http.StatusUnauthorized, "Invalid token claims: "+err.Error())
            return
        }

        // Token is valid, continue to next handler
        next.ServeHTTP(w, r)
    })
}

// UserContext - Extracts user info and adds to request context
// Runs AFTER Authenticator
func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get token from context
        token, _, err := jwtauth.FromContext(r.Context())
        if err != nil || token == nil {
            writeError(w, http.StatusUnauthorized, "Unauthorized")
            return
        }

        // Extract claims
        privateClaims := token.PrivateClaims()
        claims, err := auth.FromTokenClaims(privateClaims)
        if err != nil {
            writeError(w, http.StatusUnauthorized, "Invalid token claims")
            return
        }

        // Create UserContext and add to request context
        userContext := auth.NewUserContext(claims)
        ctx := userContext.AddToContext(r.Context())

        // Pass updated context to next handler
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Helper function to write JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write([]byte(`{"error":"` + message + `"}`))
}
```

### Middleware Chain Explanation:

Request
   ↓
[1] jwtauth.Verifier (library)
   • Extracts "Bearer token" from Authorization header
   • Verifies signature with secret key
   • Checks expiration (exp claim)
   • Stores token + any errors in context
   • ALWAYS continues (even if invalid)
   ↓
[2] Authenticator (custom)
   • Retrieves token from context
   • Checks if Verifier found errors
   • Validates claims structure (user_id, email, name)
   • Returns 401 if invalid
   • Continues if valid
   ↓
[3] UserContext (custom)
   • Extracts claims from token
   • Creates UserContext struct
   • Adds UserContext to request context
   • Handler can now call auth.FromContext()
   ↓
Handler
   • Gets user info from context
   • Uses user.ID for business logic

### 4.5 Delivery Layer - Handlers
delivery/web/handlers/user.go

```go
package handlers

import (
    "encoding/json"
    "errors"
    "net/http"
    "time"

    "github.com/go-chi/jwtauth/v5"
    "github.com/go-playground/validator/v10"
    "github.com/macesz/todo-go/delivery/web/auth"
    "github.com/macesz/todo-go/delivery/web/utils"
    "github.com/macesz/todo-go/domain"
)

type UserHandlers struct {
    Service   domain.UserService
    TokenAuth *jwtauth.JWTAuth
    Validator *validator.Validate
}

func NewUserHandlers(service domain.UserService, tokenAuth *jwtauth.JWTAuth) *UserHandlers {
    return &UserHandlers{
        Service:   service,
        TokenAuth: tokenAuth,
        Validator: validator.New(),
    }
}

// CreateUser - Register new user
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    // Parse request
    var req domain.CreateUserDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{
            Error: "Invalid request body",
        })
        return
    }

    // Validate input
    if err := h.Validator.Struct(req); err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{
            Error: translateValidationError(err),
        })
        return
    }

    // Create user (business logic)
    user, err := h.Service.CreateUser(r.Context(), req.Name, req.Email, req.Password)
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrDuplicate):
            utils.WriteJSON(w, http.StatusConflict, domain.ErrorResponse{
                Error: "Email already exists",
            })
            return
        default:
            utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
                Error: "Internal server error",
            })
            return
        }
    }

    // Return response (without password)
    respUser := domain.UserResponseDTO{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }

    utils.WriteJSON(w, http.StatusCreated, respUser)
}

// Login - Authenticate and return JWT
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    // Parse request
    var req domain.LoginDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{
            Error: "Invalid request body",
        })
        return
    }

    // Validate input
    if err := h.Validator.Struct(req); err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{
            Error: translateValidationError(err),
        })
        return
    }

    // Authenticate user (business logic)
    user, err := h.Service.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        if errors.Is(err, domain.ErrInvalidCredentials) {
            utils.WriteJSON(w, http.StatusUnauthorized, domain.ErrorResponse{
                Error: "Invalid email or password",
            })
            return
        }
        utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
            Error: "Internal server error",
        })
        return
    }

    // Generate JWT (delivery concern)
    claims := auth.NewUserClaims(user, 24*time.Hour)
    _, tokenString, err := h.TokenAuth.Encode(claims.ToMap())
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
            Error: "Failed to generate token",
        })
        return
    }

    // Return response
    utils.WriteJSON(w, http.StatusOK, domain.LoginResponseDTO{
        Token: tokenString,
        User: domain.UserResponseDTO{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
        },
    })
}

// GetUser - Get current authenticated user (protected route)
func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
    // Get authenticated user from context (set by middleware)
    userCtx, ok := auth.FromContext(r.Context())
    if !ok {
        utils.WriteJSON(w, http.StatusUnauthorized, domain.ErrorResponse{
            Error: "Unauthorized",
        })
        return
    }

    // Fetch user details from service
    user, err := h.Service.GetUser(r.Context(), userCtx.ID)
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
            Error: "Failed to fetch user",
        })
        return
    }

    // Return response (without password)
    utils.WriteJSON(w, http.StatusOK, domain.UserResponseDTO{
        ID:    user.ID,
        Email: user.Email,
        Name:  user.Name,
    })
}

// Helper function to translate validation errors
func translateValidationError(err error) string {
    validationErrs, ok := err.(validator.ValidationErrors)
    if !ok {
        return "Validation error"
    }

    var messages []string
    for _, e := range validationErrs {
        field := e.Field()
        tag := e.Tag()

        var message string
        switch tag {
        case "required":
            message = field + " is required"
        case "email":
            message = field + " must be a valid email address"
        case "min":
            message = field + " must be at least " + e.Param() + " characters"
        case "max":
            message = field + " must be at most " + e.Param() + " characters"
        case "containsany":
            if e.Param() == "0123456789" {
                message = field + " must contain at least one digit"
            } else if e.Param() == "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
                message = field + " must contain at least one uppercase letter"
            } else {
                message = field + " must contain specific characters"
            }
        default:
            message = field + " is invalid"
        }

        messages = append(messages, message)
    }

    return strings.Join(messages, "; ")
}
```

delivery/web/utils/response.go

```go
package utils

import (
    "encoding/json"
    "net/http"
)

// WriteJSON - Helper to write JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

### 4.6 Router Setup

```go
package web

import (
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/jwtauth/v5"
    "github.com/macesz/todo-go/delivery/web/handlers"
    "github.com/macesz/todo-go/delivery/web/middlewares"
)

func SetupRouter(userHandlers *handlers.UserHandlers, todoHandlers *handlers.TodoHandlers) http.Handler {
    r := chi.NewRouter()

    // JWT configuration
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        panic("JWT_SECRET environment variable is required")
    }
    tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

    // Set TokenAuth in handlers
    userHandlers.TokenAuth = tokenAuth

    // Global middlewares
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)

    // ────────────────────────────────────────────────────────────
    // PUBLIC ROUTES (no authentication required)
    // ────────────────────────────────────────────────────────────
    r.Group(func(r chi.Router) {
        r.Post("/register", userHandlers.CreateUser)
        r.Post("/login", userHandlers.Login)
    })

    // ────────────────────────────────────────────────────────────
    // PROTECTED ROUTES (JWT authentication required)
    // ────────────────────────────────────────────────────────────
    r.Group(func(r chi.Router) {
        // JWT Middleware Chain
        r.Use(jwtauth.Verifier(tokenAuth))     // 1. Extract & verify token
        r.Use(middlewares.Authenticator)        // 2. Check validation result
        r.Use(middlewares.UserContext)          // 3. Add user to context

        // User routes
        r.Get("/user", userHandlers.GetUser)

        // Todo routes
        r.Route("/todos", func(r chi.Router) {
            r.Get("/", todoHandlers.ListTodos)
            r.Post("/", todoHandlers.CreateTodo)
            r.Get("/{id}", todoHandlers.GetTodo)
            r.Put("/{id}", todoHandlers.UpdateTodo)
            r.Delete("/{id}", todoHandlers.DeleteTodo)
        })
    })

    return r
}
```

cmd/server/main.go

```go
package main

import (
    "log"
    "net/http"

    "github.com/macesz/todo-go/delivery/web"
    "github.com/macesz/todo-go/delivery/web/handlers"
    "github.com/macesz/todo-go/repository"
    "github.com/macesz/todo-go/service"
)

func main() {
    // Setup database connection
    db, err := setupDatabase()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Create repositories
    userRepo := repository.NewUserRepository(db)
    todoRepo := repository.NewTodoRepository(db)

    // Create services
    userService := service.NewUserService(userRepo)
    todoService := service.NewTodoService(todoRepo)

    // Create handlers
    userHandlers := handlers.NewUserHandlers(userService, nil)
    todoHandlers := handlers.NewTodoHandlers(todoService)

    // Setup router
    router := web.SetupRouter(userHandlers, todoHandlers)

    // Start server
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}
```


### 5. Security Best Practices
5.1 Secret Key Management
❌ NEVER hardcode secrets

// ❌ TERRIBLE - Secret in code
tokenAuth := jwtauth.New("HS256", []byte("mysecret123"), nil)

// ❌ BAD - Secret in config file committed to Git
config.yaml:
  jwt_secret: "mysecret123"

✅ Use environment variables
// ✅ GOOD - Load from environment
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

.env file (NOT committed to Git):
bash
JWT_SECRET=8xK7vZ2qT9nR4wL6pY3mJ5hF8dS1cX0bV7nQ9rP4kM2e

.gitignore:
bash
.env
*.env

Generate Strong Secrets
Minimum Requirements:

✅ At least 256 bits (32 bytes) for HS256
✅ Cryptographically random
✅ Different for each environment

Generate secret:
# Using openssl
openssl rand -base64 32

# Using Go
```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

func main() {
    b := make([]byte, 32)
    rand.Read(b)
    fmt.Println(base64.StdEncoding.EncodeToString(b))
}
```

Ok I need you to to go through step by step on all the thing we develop for JWT i didnt had time to learn since last week, so I need refreshing, and I refresh you as well with the chi JWT authentication middleware for Go HTTP services, that Verifier DOES validate expiration and evryting we need

I would like a to get a detaild explanation to learn the concepts and also a step by step guide for all the implementatiuon
Claude Sonnet 4.5
Gondolkodási folyamat



Complete Guide: JWT Authentication in Go with Chi
A comprehensive, step-by-step tutorial on implementing JWT authentication for REST APIs.

Part 1: Understanding JWT Concepts
What is JWT?
JWT (JSON Web Token) is a compact, URL-safe token format for securely transmitting information between parties.

Structure: Three parts separated by dots (.)

css
header.payload.signature
Example:

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6ImpvaG5AZXhhbXBsZS5jb20ifQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
JWT Parts Explained
1. Header (Algorithm & Token Type)
json
{
  "alg": "HS256",
  "typ": "JWT"
}
alg: Signing algorithm (HMAC SHA256)
typ: Token type
2. Payload (Claims - Your Data)
json
{
  "user_id": 1,
  "email": "john@example.com",
  "name": "John Doe",
  "exp": 1735689600
}
Standard Claims:

exp (expiration): When token expires (Unix timestamp)
iat (issued at): When token was created
nbf (not before): Token not valid before this time
Custom Claims:

user_id, email, name: Your application data
3. Signature (Verification)
scss
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  your-secret-key
)
Purpose: Ensures token hasn't been tampered with

How JWT Authentication Works
arduino
┌─────────────┐                                  ┌─────────────┐
│   Client    │                                  │   Server    │
└──────┬──────┘                                  └──────┬──────┘
       │                                                │
       │  1. POST /login (email, password)             │
       │─────────────────────────────────────────────>│
       │                                                │
       │                     2. Verify credentials     │
       │                        Create JWT token       │
       │                                                │
       │  3. Response: { "token": "eyJ..." }           │
       │<─────────────────────────────────────────────│
       │                                                │
       │  Store token (localStorage, memory, etc.)     │
       │                                                │
       │  4. GET /todos                                │
       │     Authorization: Bearer eyJ...              │
       │─────────────────────────────────────────────>│
       │                                                │
       │                     5. Verify token signature │
       │                        Check expiration       │
       │                        Extract user_id        │
       │                                                │
       │  6. Response: [{ "id": 1, "title": "..." }]   │
       │<─────────────────────────────────────────────│
       │                                                │
Key Points:

Stateless: Server doesn't store sessions
Self-contained: Token contains all user info
Secure: Signature prevents tampering
Expirable: Automatic timeout via exp claim
Part 2: How go-chi/jwtauth Works
The Middleware Chain
go
r.Use(jwtauth.Verifier(tokenAuth))    // Step 1: Extract & validate token
r.Use(middlewares.Authenticator)       // Step 2: Check validation result
r.Use(middlewares.UserContext)         // Step 3: Add user info to context
What Each Middleware Does
1. jwtauth.Verifier (Built-in)
Job: Extract and validate JWT token

go
func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. Extract token from request
            token := extractFromHeader(r)  // Authorization: Bearer <token>

            // 2. Verify signature using secret key
            decodedToken, err := ja.Decode(token)

            // 3. Validate expiration, not-before, issued-at
            if err := jwt.Validate(decodedToken); err != nil {
                // Token invalid, but continue (stores error in context)
            }

            // 4. Store token + error in context
            ctx := jwtauth.NewContext(r.Context(), decodedToken, err)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
What it validates:

✅ Signature: Token signed with correct secret?
✅ Expiration (exp): Is token expired?
✅ Not-Before (nbf): Is token valid yet?
✅ Issued-At (iat): Was token issued in the past?
Important: Even if validation fails, it continues to next middleware (stores error in context)

2. Authenticator (Custom)
Job: Check validation result and reject invalid requests

go
func Authenticator(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get token + error from context (set by Verifier)
        token, _, err := jwtauth.FromContext(r.Context())

        // If Verifier found any issues, reject request
        if err != nil || token == nil {
            http.Error(w, "Unauthorized", 401)
            return
        }

        // Validate claims structure
        claims, err := auth.FromTokenClaims(token.PrivateClaims())
        if err != nil {
            http.Error(w, "Invalid token claims", 401)
            return
        }

        // Token is valid, continue
        next.ServeHTTP(w, r)
    })
}
Why separate from Verifier?

Custom error responses (JSON instead of plain text)
Custom claim validation
Business logic checks
3. UserContext (Custom)
Job: Extract user info and add to request context

go
func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, _, _ := jwtauth.FromContext(r.Context())

        // Extract claims
        claims, _ := auth.FromTokenClaims(token.PrivateClaims())

        // Create user context
        userContext := auth.UserContext{
            ID:    claims.UserID,
            Email: claims.Email,
            Name:  claims.Name,
        }

        // Add to request context
        ctx := userContext.AddToContext(r.Context())
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
Why needed?

Handlers need user ID to fetch user-specific data
Avoids parsing token in every handler
Type-safe access to user info
Part 3: Complete Implementation
Project Structure
bash
todo-go/
├── domain/
│   ├── user.go           # Business entity
│   ├── todo.go
│   ├── errors.go
│   └── config.go
│
├── delivery/
│   └── web/
│       ├── auth/
│       │   ├── claims.go     # JWT claims structure
│       │   └── context.go    # User context helpers
│       ├── handlers/
│       │   ├── user.go       # Login, Register handlers
│       │   └── todo.go       # Todo CRUD handlers
│       ├── middlewares/
│       │   └── auth.go       # Custom auth middlewares
│       ├── utils/
│       │   └── response.go   # JSON response helpers
│       └── server.go         # Router setup
│
├── service/
│   └── user_service.go   # Business logic
│
└── repository/
    └── user_repo.go      # Database access
Step 1: Domain Layer (Business Entities)
domain/user.go - Pure business logic, no JWT/HTTP knowledge

go
package domain

import "context"

// User - Core business entity
type User struct {
	ID       int64
	Name     string
	Email    string
	Password string // Hashed password
}

// Business validation
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if len(u.Password) < 8 {
		return ErrWeakPassword
	}
	return nil
}

// UserService - Business logic interface
type UserService interface {
	CreateUser(ctx context.Context, name, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
}
domain/errors.go - Business errors

go
package domain

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
)
domain/config.go - Application configuration

go
package domain

type Config struct {
	JWTSecret string
	DBPath    string
	Port      string
}
Step 2: Auth Layer (JWT Infrastructure)
delivery/web/auth/claims.go - JWT claims structure

go
package auth

import (
	"errors"
	"time"

	"github.com/macesz/todo-go/domain"
)

// UserClaims - JWT payload structure
type UserClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Exp    int64  `json:"exp"`
}

// NewUserClaims - Convert domain.User to JWT claims
func NewUserClaims(user *domain.User, expiresIn time.Duration) UserClaims {
	return UserClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Exp:    time.Now().Add(expiresIn).Unix(),
	}
}

// ToMap - Convert to map for jwtauth library
func (c UserClaims) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id": c.UserID,
		"email":   c.Email,
		"name":    c.Name,
		"exp":     c.Exp,
	}
}

// FromTokenClaims - Extract and validate claims from JWT token
// IMPORTANT: JWT stores numbers as float64, not int64!
func FromTokenClaims(claims map[string]interface{}) (*UserClaims, error) {
	// user_id validation
	userID, ok := claims["user_id"].(float64) // Must be float64!
	if !ok {
		return nil, errors.New("missing or invalid user_id claim")
	}

	// email validation
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("missing or invalid email claim")
	}

	// name validation
	name, ok := claims["name"].(string)
	if !ok {
		return nil, errors.New("missing or invalid name claim")
	}

	return &UserClaims{
		UserID: int64(userID), // Convert float64 to int64
		Email:  email,
		Name:   name,
	}, nil
}
Why float64 not int64?

JSON specification doesn't distinguish integer/float
All numbers in JSON (and JWT) are float64
Must convert: int64(userID)
delivery/web/auth/context.go - HTTP context helpers

go
package auth

import "context"

type contextKey string

const userContextKey contextKey = "userContext"

// UserContext - User information in request context
type UserContext struct {
	ID    int64
	Email string
	Name  string
}

// NewUserContext - Create from JWT claims
func NewUserContext(claims *UserClaims) *UserContext {
	return &UserContext{
		ID:    claims.UserID,
		Email: claims.Email,
		Name:  claims.Name,
	}
}

// AddToContext - Store user context in request context
func (u *UserContext) AddToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// FromContext - Retrieve user context from request context
func FromContext(ctx context.Context) (*UserContext, bool) {
	userCtx, ok := ctx.Value(userContextKey).(*UserContext)
	return userCtx, ok
}
Step 3: Middlewares
delivery/web/middlewares/auth.go - Authentication middlewares

go
package middlewares

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/delivery/web/utils"
)

// Authenticator - Validates token and claims structure
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from context (set by Verifier)
		token, _, err := jwtauth.FromContext(r.Context())

		// Check if Verifier found issues (expired, invalid signature, etc.)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized - " + err.Error(),
			})
			return
		}

		if token == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Missing authentication token",
			})
			return
		}

		// Validate claims structure (user_id, email, name exist and correct type)
		privateClaims := token.PrivateClaims()
		_, err = auth.FromTokenClaims(privateClaims)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid token claims - " + err.Error(),
			})
			return
		}

		// Token is valid, continue to next middleware
		next.ServeHTTP(w, r)
	})
}

// UserContext - Extracts user info and adds to request context
func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from context
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil || token == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized",
			})
			return
		}

		// Extract claims (already validated by Authenticator)
		privateClaims := token.PrivateClaims()
		claims, err := auth.FromTokenClaims(privateClaims)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid token claims",
			})
			return
		}

		// Create UserContext and add to request context
		userContext := auth.NewUserContext(claims)
		ctx := userContext.AddToContext(r.Context())

		// Pass request with updated context to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
Step 4: Handlers
delivery/web/handlers/user.go - User handlers

go
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/delivery/web/utils"
	"github.com/macesz/todo-go/domain"
)

type UserHandlers struct {
	Service   domain.UserService
	TokenAuth *jwtauth.JWTAuth
}

// Login - Authenticate user and return JWT token
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}
	defer r.Body.Close()

	// Validate input
	if req.Email == "" || req.Password == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
		return
	}

	// Authenticate user (business logic)
	user, err := h.Service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid email or password",
			})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
		return
	}

	// Create JWT claims (24 hour expiration)
	claims := auth.NewUserClaims(user, 24*time.Hour)

	// Generate JWT token
	_, tokenString, err := h.TokenAuth.Encode(claims.ToMap())
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
		return
	}

	// Return token and user info
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// CreateUser - Register new user
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}
	defer r.Body.Close()

	// Create user (business logic)
	user, err := h.Service.CreateUser(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			utils.WriteJSON(w, http.StatusConflict, map[string]string{
				"error": "Email already exists",
			})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
		return
	}

	// Return created user (without password)
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

// GetUser - Get current user info (protected route)
func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by UserContext middleware)
	userCtx, ok := auth.FromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Fetch full user details
	user, err := h.Service.GetUser(r.Context(), userCtx.ID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}
delivery/web/handlers/todo.go - Todo handlers (using auth context)

go
package handlers

import (
	"net/http"

	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/delivery/web/utils"
	"github.com/macesz/todo-go/domain"
)

type TodoHandlers struct {
	Service domain.TodoService
}

// ListTodos - Get all todos for authenticated user
func (h *TodoHandlers) ListTodos(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userCtx, ok := auth.FromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Fetch user's todos
	todos, err := h.Service.ListTodos(r.Context(), userCtx.ID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to list todos",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, todos)
}

// CreateTodo - Create new todo for authenticated user
func (h *TodoHandlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := auth.FromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}
	defer r.Body.Close()

	// Create todo for authenticated user
	todo, err := h.Service.CreateTodo(r.Context(), userCtx.ID, req.Title, req.Description)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create todo",
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, todo)
}
Step 5: Utilities
delivery/web/utils/response.go - JSON response helpers

go
package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSON - Helper to write JSON responses
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// JsonError - Convert error to JSON string (for http.Error compatibility)
func JsonError(err error) string {
	data, _ := json.Marshal(map[string]string{"error": err.Error()})
	return string(data)
}
Step 6: Server Setup
delivery/web/server.go - Router and middleware setup

go
package web

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/handlers"
	"github.com/macesz/todo-go/delivery/web/middlewares"
	"github.com/macesz/todo-go/domain"
)

type Handlers struct {
	Todo *handlers.TodoHandlers
	User *handlers.UserHandlers
}

func StartServer(ctx context.Context, conf domain.Config, hdlrs *Handlers) {
	// Create Chi router
	r := chi.NewRouter()

	// JWT Auth setup (HS256 with secret from config)
	tokenAuth := jwtauth.New("HS256", []byte(conf.JWTSecret), nil)

	// Pass tokenAuth to user handlers (needed for Login)
	hdlrs.User.TokenAuth = tokenAuth

	// Global middlewares
	r.Use(middleware.RequestID)  // Add unique request ID
	r.Use(middleware.RealIP)     // Set real client IP
	r.Use(middleware.Logger)     // Log requests
	r.Use(middleware.Recoverer)  // Recover from panics

	// ============================================
	// PUBLIC ROUTES (No authentication required)
	// ============================================
	r.Group(func(r chi.Router) {
		r.Post("/user", hdlrs.User.CreateUser)  // Register
		r.Post("/login", hdlrs.User.Login)      // Login
	})

	// ============================================
	// PROTECTED ROUTES (JWT authentication required)
	// ============================================
	r.Group(func(r chi.Router) {
		// JWT Middleware Chain
		r.Use(jwtauth.Verifier(tokenAuth))     // Extract & validate token
		r.Use(middlewares.Authenticator)        // Check validation result
		r.Use(middlewares.UserContext)          // Add user info to context

		// Content type validation
		r.Use(middleware.AllowContentType("application/json"))

		// Todo routes
		r.Route("/todos", func(r chi.Router) {
			r.Get("/", hdlrs.Todo.ListTodos)         // GET /todos
			r.Get("/{id}", hdlrs.Todo.GetTodo)       // GET /todos/123
			r.Post("/", hdlrs.Todo.CreateTodo)       // POST /todos
			r.Put("/{id}", hdlrs.Todo.UpdateTodo)    // PUT /todos/123
			r.Delete("/{id}", hdlrs.Todo.DeleteTodo) // DELETE /todos/123
		})

		// User routes
		r.Route("/user", func(r chi.Router) {
			r.Get("/", hdlrs.User.GetUser)           // GET /user (current user)
			r.Delete("/{id}", hdlrs.User.DeleteUser) // DELETE /user/123
		})
	})

	// Start server
	log.Printf("Server listening on :%s\n", conf.Port)
	if err := http.ListenAndServe(":"+conf.Port, r); err != nil {
		log.Fatal(err)
	}
}
Part 4: Testing the Implementation
1. Register a New User
bash
curl -X POST http://localhost:3000/user \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
  }'
Response:

json
{
  "id": 1,
  "email": "john@example.com",
  "name": "John Doe"
}
2. Login (Get JWT Token)
bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
Response:

json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6ImpvaG5AZXhhbXBsZS5jb20iLCJuYW1lIjoiSm9obiBEb2UiLCJleHAiOjE3MzU3NzYwMDB9.abc123xyz",
  "user": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe"
  }
}
Save the token! You'll need it for protected routes.

3. Access Protected Route (Create Todo)
bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X POST http://localhost:3000/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread"
  }'
Response:

json
{
  "id": 1,
  "user_id": 1,
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "completed": false,
  "created_at": "2025-01-02T10:30:00Z"
}
4. Test Without Token (Should Fail)
bash
curl -X GET http://localhost:3000/todos
Response:

json
{
  "error": "Missing authentication token"
}
Status: 401 Unauthorized

5. Test With Expired Token (Should Fail)
After 24 hours, the token expires automatically.

bash
curl -X GET http://localhost:3000/todos \
  -H "Authorization: Bearer <expired-token>"
Response:

json
{
  "error": "Unauthorized - exp not satisfied"
}
Status: 401 Unauthorized

Part 5: Key Concepts Summary
Middleware Flow Visualization
sql
Request: GET /todos
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

        ↓
┌─────────────────────────────────────────┐
│  1. jwtauth.Verifier                    │
│     • Extract token from header         │
│     • Verify signature with secret      │
│     • Validate expiration (exp)         │
│     • Validate not-before (nbf)         │
│     • Store token + error in context    │
└─────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────┐
│  2. middlewares.Authenticator           │
│     • Get token from context            │
│     • Check if error != nil             │
│     • Validate claims structure         │
│     • Reject if invalid                 │
└─────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────┐
│  3. middlewares.UserContext             │
│     • Extract claims from token         │
│     • Create UserContext                │
│     • Add to request context            │
└─────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────┐
│  4. Handler (ListTodos)                 │
│     • Get user from context             │
│     • Fetch todos for user.ID           │
│     • Return JSON response              │
└─────────────────────────────────────────┘
Context Flow
go
// Request starts
ctx := r.Context()

// After Verifier
ctx = context.WithValue(ctx, "jwt.token", token)
ctx = context.WithValue(ctx, "jwt.error", err)

// After UserContext
ctx = context.WithValue(ctx, "userContext", &UserContext{
    ID:    1,
    Email: "john@example.com",
    Name:  "John Doe",
})

// In Handler
userCtx, _ := auth.FromContext(r.Context())
// userCtx.ID = 1
// userCtx.Email = "john@example.com"
Why Three Middlewares?
Single Responsibility Principle:

Táblázat
Middleware	Responsibility
Verifier	Extract & cryptographically validate token
Authenticator	Business logic validation (claims structure)
UserContext	Extract user info for handlers
Benefits:

✅ Easy to test individually
✅ Reusable components
✅ Custom error handling per layer
✅ Clear separation of concerns
Part 6: Common Pitfalls & Solutions
❌ Pitfall 1: Type Assertion for Numbers
go
// WRONG - Will always fail!
userID, ok := claims["user_id"].(int64)  // ok = false

// CORRECT - JWT uses float64
userID, ok := claims["user_id"].(float64)  // ok = true
intUserID := int64(userID)
Why: JSON doesn't distinguish integer vs float

❌ Pitfall 2: Forgetting to Return After Error
go
// WRONG - Continues after error!
if err != nil {
    http.Error(w, "Unauthorized", 401)
    // Missing return!
}
next.ServeHTTP(w, r)  // Still executes!

// CORRECT
if err != nil {
    http.Error(w, "Unauthorized", 401)
    return  // Stop here
}
next.ServeHTTP(w, r)
❌ Pitfall 3: Not Using WithContext
go
// WRONG - Context not updated!
func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := addUserToContext(r.Context())
        next.ServeHTTP(w, r)  // Uses old context!
    })
}

// CORRECT
func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := addUserToContext(r.Context())
        next.ServeHTTP(w, r.WithContext(ctx))  // Pass updated context
    })
}
❌ Pitfall 4: Storing Sensitive Data in JWT
go
// WRONG - JWT is NOT encrypted, only signed!
claims := map[string]interface{}{
    "user_id": user.ID,
    "password": user.Password,  // Visible to anyone!
    "credit_card": "1234-5678", // Very bad!
}

// CORRECT - Only non-sensitive data
claims := map[string]interface{}{
    "user_id": user.ID,
    "email": user.Email,
    "name": user.Name,
}
JWT is Base64 encoded, NOT encrypted! Anyone can decode and read it.

Part 7: Security Best Practices
1. Use Strong Secret Keys
go
// ❌ BAD
JWTSecret: "secret"

// ✅ GOOD (minimum 256 bits / 32 bytes for HS256)
JWTSecret: "a7f8e9d2c1b4a6e5f3d8c7b9a2e1f4d6c8b5a3e2f1d9c7b6a4e3f2d1c9b8a7e6"
Generate strong secret:

bash
openssl rand -base64 32
2. Set Reasonable Expiration Times
go
// Short-lived tokens (more secure)
claims := auth.NewUserClaims(user, 1*time.Hour)    // 1 hour
claims := auth.NewUserClaims(user, 15*time.Minute) // 15 minutes

// Long-lived tokens (convenience vs security trade-off)
claims := auth.NewUserClaims(user, 24*time.Hour)   // 24 hours
claims := auth.NewUserClaims(user, 7*24*time.Hour) // 7 days
Best practice: Short-lived access tokens + refresh tokens

3. Use HTTPS in Production
go
// Development
http.ListenAndServe(":3000", r)

// Production (HTTPS)
http.ListenAndServeTLS(":443", r, "cert.pem", "key.pem")
Why: JWT in HTTP headers can be intercepted (man-in-the-middle attack)

4. Validate Token on Every Request
go
// ✅ CORRECT - Middleware runs on every protected route
r.Group(func(r chi.Router) {
    r.Use(jwtauth.Verifier(tokenAuth))
    r.Use(middlewares.Authenticator)
    // Protected routes
})

// ❌ WRONG - Validating once and storing result
// Don't cache validation results!
5. Don't Store Tokens in localStorage (XSS Risk)
Better alternatives:

Memory: Store in JavaScript variable (lost on page refresh)
HttpOnly Cookie: Cannot be accessed by JavaScript (prevents XSS)
Secure Cookie: Only sent over HTTPS
Conclusion
You now have a complete, production-ready JWT authentication system with:

✅ Clean architecture - Domain/Delivery/Infrastructure separation
✅ Type-safe claims - No magic strings
✅ Proper validation - Signature, expiration, claims structure
✅ Middleware chain - Verifier → Authenticator → UserContext
✅ Security best practices - Strong secrets, expiration, HTTPS
✅ Testable code - Clear responsibilities per layer

Key Takeaway: jwtauth.Verifier handles all cryptographic validation (signature + expiration). Your custom middlewares add business logic validation and user context extraction.

Ready to deploy! 🚀

I try to fix my tests for handler


// CreateUser creates a new HTTP handler for creating a new user.
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqUser domain.CreateUserDTO // Empty User struct to decode into

	// Decode the JSON body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		// domain.ErrorResponse{Error: err.Error() for dynamic error message
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	if err := validate.New().Struct(reqUser); err != nil {
		useErr := translateValidationError(err)
		// Dynamic message, e.g., "Name is required; Email is required"
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: useErr})
		return
	}

	// Create the user using the service
	user, err := h.Service.CreateUser(r.Context(), reqUser.Name, reqUser.Email, reqUser.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmail):
			utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
			return
		case errors.Is(err, domain.ErrInvalidPassword):
			utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
			return
		case errors.Is(err, domain.ErrDuplicate):
			utils.WriteJSON(w, http.StatusConflict, domain.ErrorResponse{Error: err.Error()})
			return
		default:
			utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"})
			return
		}
	}

	respUser := domain.UserResponseDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	utils.WriteJSON(w, http.StatusCreated, respUser)
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

--- FAIL: TestCreateUser (0.00s)
    --- FAIL: TestCreateUser/Valid_input (0.00s)
        handler_test.go:121:
            	Error Trace:	/Users/orsi/code/go/go-todo-chi/delivery/web/user/handler_test.go:121
            	Error:      	Not equal:
            	            	expected: 201
            	            	actual  : 400
            	Test:       	TestCreateUser/Valid_input
        mock_user_service.go:23: FAIL:	CreateUser(string,string,string,string)
            		at: [/Users/orsi/code/go/go-todo-chi/delivery/web/user/handler_test.go:103]
        mock_user_service.go:23: FAIL: 0 out of 1 expectation(s) were met.
            	The code you are testing needs to make 1 more call(s).
            	at: [/Users/orsi/code/go/go-todo-chi/delivery/web/user/mocks/mock_user_service.go:23 /opt/homebrew/Cellar/go/1.25.2/libexec/src/testing/testing.go:1308 /opt/homebrew/Cellar/go/1.25.2/libexec/src/testing/testing.go:1572 /opt/homebrew/Cellar/go/1.25.2/libexec/src/testing/testing.go:1928 /opt/homebrew/Cellar/go/1.25.2/libexec/src/runtime/panic.go:615 /opt/homebrew/Cellar/go/1.25.2/libexec/src/testing/testing.go:1013 /Users/orsi/code/go/go-todo-chi/delivery/web/user/handler_test.go:121]
    --- FAIL: TestCreateUser/Internal_server_error (0.00s)
        handler_test.go:121:
            	Error Trace:	/Users/orsi/code/go/go-todo-chi/delivery/web/user/handler_test.go:121
            	Error:      	Not equal:
            	            	expected: 500
            	            	actual  : 400
            	Test:       	TestCreateUser/Internal_server_error
        mock_user_service.go:23: FAIL:	CreateUser(string,string,string,string)
            		at: [/Users/orsi/code/go/go-todo-chi/delivery/web/user/handler_test.go:103]
        mock_user_service.go:23: FAIL: 0 out of 1 expectation(s) were met.
            	The code you are testing needs to make 1 more call(s).
            	at: [/Users/orsi/code/go/go-todo-chi/delivery/web/user/mocks/mock_user_service.go:23 /opt/homebrew/Cellar/go/1.25.2/libexec/src/testing/testing.go:1308 /opt/homebrew/Cellar/go/1.25.2/libexec/src/testing/testing.go:1572 /opt/homebrew/Cellar/go/1.25.2/libexec/src/testing/testing.go:1928 /opt/homebrew/Cellar/go/1.25.2/libexec/src/runtime/panic.go:615 /opt/homebrew/Cellar/go/1.25.2/libexec/src/testing/testing.go:1013 /Users/orsi/code/go/go-todo-chi/delivery/web/user/handler_test.go:121]
FAIL
exit status 1
FAIL	github.com/macesz/to
Claude Sonnet 4.5
Gondolkodási folyamat











Looking at the test failures, the issue is clear: validation is failing before the service is even called. The tests expect 201/500 but get 400, and the mock says "0 out of 1 expectation(s) were met."

The problem is in your validation step. Let me help you debug this:

Step 1: Check Your DTO Validation Tags
Show me your domain.CreateUserDTO struct. It probably looks like this:

go
type CreateUserDTO struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
The validation tags are probably stricter than your test data!

Step 2: Debug - Print What's Failing
Add temporary debug output to see what's failing:

go
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    var reqUser domain.CreateUserDTO

    if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
        return
    }

    // 🔍 DEBUG: Print what was decoded
    fmt.Printf("DEBUG: Decoded user: %+v\n", reqUser)

    if err := validate.New().Struct(reqUser); err != nil {
        useErr := translateValidationError(err)
        // 🔍 DEBUG: Print validation error
        fmt.Printf("DEBUG: Validation error: %s\n", useErr)

        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: useErr})
        return
    }

    // ... rest of handler
}
Run your test again and check the output!

Step 3: Common Validation Issues
Issue 1: Email Validation
If your DTO has validate:"email", then "test@example.com" should work, but check if there's extra whitespace or encoding issues.

Issue 2: Password Length
If you have validate:"min=8", but your test uses "password" (8 chars), it should pass. But if you have validate:"min=9", it will fail.

Issue 3: Context in Mock
Your mock shows 4 string arguments, but you're passing context.Context + 3 strings. Make sure your mock signature matches:

go
// Mock definition should be:
func (m *MockUserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
    // Implementation
}
Step 4: Fixed Test Structure
Here's the corrected test with proper debugging:

go
func TestCreateUser(t *testing.T) {
    tests := []struct {
        name           string
        inputBody      string
        setupMock      func(*mocks.MockUserService) // Function to setup expectations
        expectedStatus int
        expectedBody   string
    }{
        {
            name:      "Valid input",
            inputBody: `{"name":"Test User","email":"test@example.com","password":"password123"}`,
            setupMock: func(m *mocks.MockUserService) {
                m.On("CreateUser",
                    mock.Anything,           // context
                    "Test User",             // name
                    "test@example.com",      // email
                    "password123",           // password
                ).Return(&domain.User{
                    ID:    1,
                    Name:  "Test User",
                    Email: "test@example.com",
                }, nil).Once()
            },
            expectedStatus: http.StatusCreated,
            expectedBody:   `{"id":1,"name":"Test User","email":"test@example.com"}`,
        },
        {
            name:           "Invalid JSON",
            inputBody:      `{"name":"Test User"`, // Malformed
            setupMock:      func(m *mocks.MockUserService) {}, // No mock call expected
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"unexpected EOF"}`,
        },
        {
            name:      "Internal server error",
            inputBody: `{"name":"Test User","email":"test@example.com","password":"password123"}`,
            setupMock: func(m *mocks.MockUserService) {
                m.On("CreateUser",
                    mock.Anything,
                    "Test User",
                    "test@example.com",
                    "password123",
                ).Return(nil, errors.New("database failure")).Once()
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   `{"error":"internal server error"}`,
        },
        {
            name:           "Missing name",
            inputBody:      `{"email":"test@example.com","password":"password123"}`,
            setupMock:      func(m *mocks.MockUserService) {}, // No mock call
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Name is required"}`,
        },
        {
            name:           "Missing email",
            inputBody:      `{"name":"Test User","password":"password123"}`,
            setupMock:      func(m *mocks.MockUserService) {}, // No mock call
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Email is required"}`,
        },
        {
            name:           "Missing password",
            inputBody:      `{"name":"Test User","email":"test@example.com"}`,
            setupMock:      func(m *mocks.MockUserService) {}, // No mock call
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Password is required"}`,
        },
        {
            name:           "Invalid email format",
            inputBody:      `{"name":"Test User","email":"invalid-email","password":"password123"}`,
            setupMock:      func(m *mocks.MockUserService) {}, // No mock call
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Email must be a valid email address"}`, // Adjust based on your translateValidationError
        },
        {
            name:           "Password too short",
            inputBody:      `{"name":"Test User","email":"test@example.com","password":"short"}`,
            setupMock:      func(m *mocks.MockUserService) {}, // No mock call
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Password must be at least 8 characters"}`, // Adjust
        },
        {
            name:      "Duplicate email",
            inputBody: `{"name":"Test User","email":"test@example.com","password":"password123"}`,
            setupMock: func(m *mocks.MockUserService) {
                m.On("CreateUser",
                    mock.Anything,
                    "Test User",
                    "test@example.com",
                    "password123",
                ).Return(nil, domain.ErrDuplicate).Once()
            },
            expectedStatus: http.StatusConflict,
            expectedBody:   `{"error":"` + domain.ErrDuplicate.Error() + `"}`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock service
            mockService := new(mocks.MockUserService)

            // Setup expectations
            tt.setupMock(mockService)

            // Create handler
            handler := &UserHandlers{
                Service: mockService,
            }

            // Create request
            req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tt.inputBody))
            req.Header.Set("Content-Type", "application/json")

            // Create response recorder
            rr := httptest.NewRecorder()

            // Call handler
            handler.CreateUser(rr, req)

            // Assert status code
            assert.Equal(t, tt.expectedStatus, rr.Code, "Status code mismatch")

            // Assert response body (normalize JSON for comparison)
            var expectedJSON, actualJSON map[string]interface{}

            if err := json.Unmarshal([]byte(tt.expectedBody), &expectedJSON); err != nil {
                t.Fatalf("Failed to parse expected JSON: %v", err)
            }

            if err := json.Unmarshal(rr.Body.Bytes(), &actualJSON); err != nil {
                t.Fatalf("Failed to parse actual JSON: %v\nBody: %s", err, rr.Body.String())
            }

            assert.Equal(t, expectedJSON, actualJSON, "Response body mismatch")

            // Assert all mock expectations were met
            mockService.AssertExpectations(t)
        })
    }
}
Step 5: Check Your translateValidationError Function
Make sure it returns the exact error messages you're testing for:

go
func translateValidationError(err error) string {
    var messages []string

    for _, err := range err.(validator.ValidationErrors) {
        field := err.Field()
        tag := err.Tag()

        switch tag {
        case "required":
            messages = append(messages, field + " is required")
        case "email":
            messages = append(messages, field + " must be a valid email address")
        case "min":
            messages = append(messages, fmt.Sprintf("%s must be at least %s characters", field, err.Param()))
        default:
            messages = append(messages, field + " is invalid")
        }
    }

    return strings.Join(messages, "; ")
}
Step 6: Verify Your DTO Struct
Make sure your DTO matches your test expectations:

go
package domain

type CreateUserDTO struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type UserResponseDTO struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
Quick Debug Command
Run a single test with verbose output:

bash
go test -v -run TestCreateUser/Valid_input ./delivery/web/user
This will show you exactly what's being validated and why it's failing.

Most Likely Issue
Based on the error "0 out of 1 expectation(s) were met", your validation is failing before the service is called. This means:

✅ JSON decoding succeeds
❌ Validation fails here ← Your issue
❌ Service is never called
Solution: Either:

Fix your test data to match validation rules
Or relax your validation rules
Or check if translateValidationError is returning the expected error message
Show me:

Your CreateUserDTO struct with validation tags
Your translateValidationError function
The output when you add the debug prints
Then I can give you the exact fix! 🔍

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=5,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=255,containsany=0123456789,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}
Claude Sonnet 4.5
Gondolkodási folyamat




Found the problem! 🎯

Your password validation is very strict:

go
Password string `json:"password" validate:"required,min=6,max=255,containsany=0123456789,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
This requires:

✅ required - Must exist
✅ min=6 - At least 6 characters
❌ containsany=0123456789 - Must contain at least one digit
❌ containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ - Must contain at least one uppercase letter
But your test uses: "password" which:

✅ Has 8 characters (meets min=6)
❌ Has NO digits (fails containsany=0123456789)
❌ Has NO uppercase letters (fails containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ)
Solution: Update Test Passwords
Change all test passwords to meet validation requirements:

go
func TestCreateUser(t *testing.T) {
    tests := []struct {
        name           string
        inputBody      string
        setupMock      func(*mocks.MockUserService)
        expectedStatus int
        expectedBody   string
    }{
        {
            name:      "Valid input",
            inputBody: `{"name":"Test User","email":"test@example.com","password":"Password123"}`, // ✅ Has uppercase + digit
            setupMock: func(m *mocks.MockUserService) {
                m.On("CreateUser",
                    mock.Anything,
                    "Test User",
                    "test@example.com",
                    "Password123", // ← Update here too
                ).Return(&domain.User{
                    ID:    1,
                    Name:  "Test User",
                    Email: "test@example.com",
                }, nil).Once()
            },
            expectedStatus: http.StatusCreated,
            expectedBody:   `{"id":1,"name":"Test User","email":"test@example.com"}`,
        },
        {
            name:           "Invalid JSON",
            inputBody:      `{"name":"Test User"`,
            setupMock:      func(m *mocks.MockUserService) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"unexpected EOF"}`,
        },
        {
            name:      "Internal server error",
            inputBody: `{"name":"Test User","email":"test@example.com","password":"Password123"}`,
            setupMock: func(m *mocks.MockUserService) {
                m.On("CreateUser",
                    mock.Anything,
                    "Test User",
                    "test@example.com",
                    "Password123", // ← Update
                ).Return(nil, errors.New("database failure")).Once()
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   `{"error":"internal server error"}`,
        },
        {
            name:           "Missing name",
            inputBody:      `{"email":"test@example.com","password":"Password123"}`, // ← Update
            setupMock:      func(m *mocks.MockUserService) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Name is required"}`, // Check exact message from translateValidationError
        },
        {
            name:           "Missing email",
            inputBody:      `{"name":"Test User","password":"Password123"}`, // ← Update
            setupMock:      func(m *mocks.MockUserService) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Email is required"}`,
        },
        {
            name:           "Missing password",
            inputBody:      `{"name":"Test User","email":"test@example.com"}`,
            setupMock:      func(m *mocks.MockUserService) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Password is required"}`,
        },
        {
            name:           "Name too short",
            inputBody:      `{"name":"Test","email":"test@example.com","password":"Password123"}`, // "Test" is 4 chars, needs 5
            setupMock:      func(m *mocks.MockUserService) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Name must be at least 5 characters"}`, // Adjust based on translateValidationError
        },
        {
            name:           "Invalid email format",
            inputBody:      `{"name":"Test User","email":"invalid-email","password":"Password123"}`,
            setupMock:      func(m *mocks.MockUserService) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Email must be a valid email address"}`,
        },
        {
            name:           "Password too short",
            inputBody:      `{"name":"Test User","email":"test@example.com","password":"Pass1"}`, // 5 chars, needs 6
            setupMock:      func(m *mocks.MockUserService) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Password must be at least 6 characters"}`,
        },
        {
            name:           "Password missing digit",
            inputBody:      `{"name":"Test User","email":"test@example.com","password":"Password"}`, // No digit
            setupMock:      func(m *mocks.MockUserService) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Password must contain at least one digit"}`, // Adjust
        },
        {
            name:           "Password missing uppercase",
            inputBody:      `{"name":"Test User","email":"test@example.com","password":"password123"}`, // No uppercase
            setupMock:      func(m *mocks.MockUserService) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Password must contain at least one uppercase letter"}`, // Adjust
        },
        {
            name:      "Duplicate email",
            inputBody: `{"name":"Test User","email":"test@example.com","password":"Password123"}`,
            setupMock: func(m *mocks.MockUserService) {
                m.On("CreateUser",
                    mock.Anything,
                    "Test User",
                    "test@example.com",
                    "Password123", // ← Update
                ).Return(nil, domain.ErrDuplicate).Once()
            },
            expectedStatus: http.StatusConflict,
            expectedBody:   `{"error":"` + domain.ErrDuplicate.Error() + `"}`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockService := new(mocks.MockUserService)
            tt.setupMock(mockService)

            handler := &UserHandlers{
                Service: mockService,
            }

            req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tt.inputBody))
            req.Header.Set("Content-Type", "application/json")

            rr := httptest.NewRecorder()

            handler.CreateUser(rr, req)

            assert.Equal(t, tt.expectedStatus, rr.Code, "Status code mismatch")

            // Normalize JSON comparison
            var expectedJSON, actualJSON map[string]interface{}

            if err := json.Unmarshal([]byte(tt.expectedBody), &expectedJSON); err != nil {
                t.Fatalf("Failed to parse expected JSON: %v", err)
            }

            if err := json.Unmarshal(rr.Body.Bytes(), &actualJSON); err != nil {
                t.Fatalf("Failed to parse actual JSON: %v\nBody: %s", err, rr.Body.String())
            }

            assert.Equal(t, expectedJSON, actualJSON, "Response body mismatch")

            mockService.AssertExpectations(t)
        })
    }
}
Update Your translateValidationError Function
Make sure it handles all your validation tags correctly:

go
func translateValidationError(err error) string {
    validationErrs, ok := err.(validator.ValidationErrors)
    if !ok {
        return "Validation error"
    }

    var messages []string

    for _, e := range validationErrs {
        field := e.Field()
        tag := e.Tag()
        param := e.Param()

        var message string
        switch tag {
        case "required":
            message = field + " is required"
        case "email":
            message = field + " must be a valid email address"
        case "min":
            if field == "Password" || field == "Name" {
                message = fmt.Sprintf("%s must be at least %s characters", field, param)
            } else {
                message = fmt.Sprintf("%s must be at least %s", field, param)
            }
        case "max":
            message = fmt.Sprintf("%s must be at most %s characters", field, param)
        case "containsany":
            // Detect what character set is missing
            if param == "0123456789" {
                message = field + " must contain at least one digit"
            } else if param == "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
                message = field + " must contain at least one uppercase letter"
            } else {
                message = fmt.Sprintf("%s must contain at least one of: %s", field, param)
            }
        default:
            message = fmt.Sprintf("%s is invalid (%s)", field, tag)
        }

        messages = append(messages, message)
    }

    return strings.Join(messages, "; ")
}
Quick Fix Summary
Main issue: Test password "password" doesn't meet validation rules.

Solution: Use "Password123" or similar that has:

✅ At least 6 characters
✅ At least one digit (0-9)
✅ At least one uppercase letter (A-Z)
Valid password examples:

"Password123"
"SecureP@ss1"
"MyPass99"
"Test1234"
Invalid password examples:

❌ "password" - no uppercase, no digit
❌ "Password" - no digit
❌ "password123" - no uppercase
❌ "Pass1" - too short (only 5 chars)
Run Updated Test
bash
go test -v -run TestCreateUser ./delivery/web/user
This should now pass! 🎉

Bonus: Test Helper for Valid Passwords
Create a helper to generate valid test passwords:

go
// test_helpers.go
package user

import "fmt"

// ValidTestPassword returns a password that meets all validation requirements
func ValidTestPassword(suffix string) string {
    return fmt.Sprintf("Password%s123", suffix)
}

// Then in tests:
inputBody: `{"name":"Test User","email":"test@example.com","password":"` + ValidTestPassword("") + `"}`,
This makes it clear that you're using a valid password format! 🔐

now I would like to make a test for Login
Claude Sonnet 4.5
Gondolkodási folyamat




Complete Login Handler Tests
Let me create comprehensive tests for your Login handler. First, let's make sure your handler is set up correctly, then write all the test cases.

Step 1: Login Handler Implementation
delivery/web/user/handler.go

go
// Login authenticates a user and returns a JWT token
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqLogin struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&reqLogin); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	// Validate input
	if err := validate.New().Struct(reqLogin); err != nil {
		useErr := translateValidationError(err)
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: useErr})
		return
	}

	// Authenticate user
	user, err := h.Service.Login(r.Context(), reqLogin.Email, reqLogin.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			utils.WriteJSON(w, http.StatusUnauthorized, domain.ErrorResponse{Error: "Invalid email or password"})
			return
		default:
			utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"})
			return
		}
	}

	// Create JWT claims (24 hour expiration)
	claims := auth.NewUserClaims(user, 24*time.Hour)

	// Generate JWT token
	_, tokenString, err := h.TokenAuth.Encode(claims.ToMap())
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	// Return token and user info
	respLogin := domain.LoginResponseDTO{
		Token: tokenString,
		User: domain.UserResponseDTO{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}

	utils.WriteJSON(w, http.StatusOK, respLogin)
}
Step 2: Add DTOs
domain/user.go - Add response DTOs

go
package domain

// LoginResponseDTO - Response for successful login
type LoginResponseDTO struct {
	Token string          `json:"token"`
	User  UserResponseDTO `json:"user"`
}

// UserResponseDTO - User information (without password)
type UserResponseDTO struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
Step 3: Complete Login Tests
delivery/web/user/handler_test.go

go
package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/macesz/todo-go/delivery/web/user/mocks"
	"github.com/macesz/todo-go/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      string
		setupMock      func(*mocks.MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "Valid credentials",
			inputBody: `{"email":"test@example.com","password":"Password123"}`,
			setupMock: func(m *mocks.MockUserService) {
				m.On("Login",
					mock.Anything, // context
					"test@example.com",
					"Password123",
				).Return(&domain.User{
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

				// Check token exists and is not empty
				assert.NotEmpty(t, response.Token, "Token should not be empty")

				// Check user info
				assert.Equal(t, int64(1), response.User.ID)
				assert.Equal(t, "test@example.com", response.User.Email)
				assert.Equal(t, "Test User", response.User.Name)

				// Verify token is valid JWT (optional but recommended)
				assert.True(t, strings.Count(response.Token, ".") == 2, "Token should have 3 parts separated by dots")
			},
		},
		{
			name:      "Invalid credentials - wrong password",
			inputBody: `{"email":"test@example.com","password":"WrongPassword123"}`,
			setupMock: func(m *mocks.MockUserService) {
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
				assert.Equal(t, "Invalid email or password", response.Error)
			},
		},
		{
			name:      "Invalid credentials - user not found",
			inputBody: `{"email":"nonexistent@example.com","password":"Password123"}`,
			setupMock: func(m *mocks.MockUserService) {
				m.On("Login",
					mock.Anything,
					"nonexistent@example.com",
					"Password123",
				).Return(nil, domain.ErrInvalidCredentials).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid email or password", response.Error)
			},
		},
		{
			name:           "Invalid JSON",
			inputBody:      `{"email":"test@example.com"`, // Malformed JSON
			setupMock:      func(m *mocks.MockUserService) {}, // No service call
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "unexpected EOF", response.Error)
			},
		},
		{
			name:           "Missing email",
			inputBody:      `{"password":"Password123"}`,
			setupMock:      func(m *mocks.MockUserService) {}, // No service call
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, "Email is required")
			},
		},
		{
			name:           "Missing password",
			inputBody:      `{"email":"test@example.com"}`,
			setupMock:      func(m *mocks.MockUserService) {}, // No service call
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, "Password is required")
			},
		},
		{
			name:           "Invalid email format",
			inputBody:      `{"email":"invalid-email","password":"Password123"}`,
			setupMock:      func(m *mocks.MockUserService) {}, // No service call
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, "Email must be a valid email address")
			},
		},
		{
			name:           "Empty email",
			inputBody:      `{"email":"","password":"Password123"}`,
			setupMock:      func(m *mocks.MockUserService) {}, // No service call
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, "Email")
			},
		},
		{
			name:           "Empty password",
			inputBody:      `{"email":"test@example.com","password":""}`,
			setupMock:      func(m *mocks.MockUserService) {}, // No service call
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, "Password is required")
			},
		},
		{
			name:      "Internal server error",
			inputBody: `{"email":"test@example.com","password":"Password123"}`,
			setupMock: func(m *mocks.MockUserService) {
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
			// Create mock service
			mockService := new(mocks.MockUserService)
			tt.setupMock(mockService)

			// Create JWT auth (with test secret)
			tokenAuth := jwtauth.New("HS256", []byte("test-secret-key-for-testing"), nil)

			// Create handler
			handler := &UserHandlers{
				Service:   mockService,
				TokenAuth: tokenAuth,
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Login(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code mismatch")

			// Run custom response checks
			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}

			// Assert all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
Step 4: Advanced Test - Verify JWT Token Claims
Add a test that actually decodes and verifies the JWT token:

go
func TestLogin_TokenValidation(t *testing.T) {
	// Setup
	mockService := new(mocks.MockUserService)
	mockService.On("Login",
		mock.Anything,
		"test@example.com",
		"Password123",
	).Return(&domain.User{
		ID:    42,
		Name:  "John Doe",
		Email: "test@example.com",
	}, nil).Once()

	// Create JWT auth with known secret
	testSecret := []byte("test-secret-key-for-testing")
	tokenAuth := jwtauth.New("HS256", testSecret, nil)

	handler := &UserHandlers{
		Service:   mockService,
		TokenAuth: tokenAuth,
	}

	// Create request
	body := `{"email":"test@example.com","password":"Password123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Execute
	handler.Login(rr, req)

	// Assert status
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response domain.LoginResponseDTO
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)

	// Decode and verify token
	token, err := jwt.Parse([]byte(response.Token), jwt.WithKey(jwa.HS256, testSecret))
	assert.NoError(t, err, "Token should be valid")

	// Verify claims
	claims := token.PrivateClaims()

	userID, ok := claims["user_id"].(float64)
	assert.True(t, ok, "user_id claim should exist")
	assert.Equal(t, int64(42), int64(userID), "user_id should match")

	email, ok := claims["email"].(string)
	assert.True(t, ok, "email claim should exist")
	assert.Equal(t, "test@example.com", email, "email should match")

	name, ok := claims["name"].(string)
	assert.True(t, ok, "name claim should exist")
	assert.Equal(t, "John Doe", name, "name should match")

	// Verify expiration is set (24 hours in future)
	exp, ok := claims["exp"].(float64)
	assert.True(t, ok, "exp claim should exist")

	expectedExp := time.Now().Add(24 * time.Hour).Unix()
	// Allow 5 second tolerance for test execution time
	assert.InDelta(t, expectedExp, int64(exp), 5, "Expiration should be ~24 hours from now")

	mockService.AssertExpectations(t)
}
Step 5: Table-Driven Helper Functions
Create helpers to make tests more readable:

delivery/web/user/test_helpers.go

go
package user

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/user/mocks"
)

// TestSetup holds common test dependencies
type TestSetup struct {
	MockService *mocks.MockUserService
	Handler     *UserHandlers
	TokenAuth   *jwtauth.JWTAuth
}

// NewTestSetup creates a new test setup with initialized dependencies
func NewTestSetup() *TestSetup {
	mockService := new(mocks.MockUserService)
	tokenAuth := jwtauth.New("HS256", []byte("test-secret-key"), nil)

	return &TestSetup{
		MockService: mockService,
		Handler: &UserHandlers{
			Service:   mockService,
			TokenAuth: tokenAuth,
		},
		TokenAuth: tokenAuth,
	}
}

// ExecuteRequest creates a request and executes it through the handler
func (ts *TestSetup) ExecuteRequest(t *testing.T, method, path, body string, handlerFunc func(http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlerFunc(rr, req)

	return rr
}
Step 6: Simplified Tests Using Helpers
Updated test using helpers:

go
func TestLogin_WithHelpers(t *testing.T) {
	t.Run("Valid credentials", func(t *testing.T) {
		// Setup
		ts := NewTestSetup()
		ts.MockService.On("Login",
			mock.Anything,
			"test@example.com",
			"Password123",
		).Return(&domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		}, nil).Once()

		// Execute
		body := `{"email":"test@example.com","password":"Password123"}`
		rr := ts.ExecuteRequest(t, "POST", "/login", body, ts.Handler.Login)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)

		var response domain.LoginResponseDTO
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NotEmpty(t, response.Token)
		assert.Equal(t, int64(1), response.User.ID)
		assert.Equal(t, "test@example.com", response.User.Email)

		ts.MockService.AssertExpectations(t)
	})
}
Step 7: Integration Test (Optional)
Test the full flow without mocks (requires test database):

go
func TestLogin_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup real dependencies
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	tokenAuth := jwtauth.New("HS256", []byte("test-secret"), nil)

	handler := &UserHandlers{
		Service:   service,
		TokenAuth: tokenAuth,
	}

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	_, err := db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		"Test User", "test@example.com", string(hashedPassword))
	assert.NoError(t, err)

	// Test login
	body := `{"email":"test@example.com","password":"Password123"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response domain.LoginResponseDTO
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "test@example.com", response.User.Email)
}
Step 8: Run Tests
bash
# Run all login tests
go test -v -run TestLogin ./delivery/web/user

# Run with coverage
go test -cover -run TestLogin ./delivery/web/user

# Run specific test
go test -v -run TestLogin/Valid_credentials ./delivery/web/user

# Skip integration tests
go test -short -v ./delivery/web/user
Summary
Your Login tests now cover:

✅ Happy path - Valid credentials return token
✅ Invalid credentials - Wrong password/email returns 401
✅ Validation errors - Missing/invalid fields return 400
✅ JSON parsing errors - Malformed JSON returns 400
✅ Internal errors - Database errors return 500
✅ Token validation - JWT is properly formatted and contains correct claims
✅ Token expiration - Token has correct expiration time

Test coverage includes:

Request validation
Service interaction
JWT token generation
Error handling
Response formatting
All tests use mocks for fast, isolated testing! 🚀

checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response domain.LoginResponseDTO
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err, "Should parse response JSON")

				// Check token exists and is not empty
				assert.NotEmpty(t, response.Token, "Token should not be empty")

				// Check user info
				assert.Equal(t, int64(1), response.User.ID)
				assert.Equal(t, "test@example.com", response.User.Email)
				assert.Equal(t, "Test User", response.User.Name)

				// Verify token is valid JWT (optional but recommended)
				assert.True(t, strings.Count(response.Token, ".") == 2, "Token should have 3 parts separated by dots")
			},

why is it in evry case and not in the loop?
Claude Sonnet 4.5
Gondolkodási folyamat





Great observation! 🎯

The checkResponse function is different for each test case because each test has different expectations. Look closely:

Each Test Case Has Its Own Custom Validation
go
{
    name: "Valid credentials",
    checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
        var response domain.LoginResponseDTO  // ← Expects LoginResponseDTO
        // ... check token and user info
    },
},
{
    name: "Invalid credentials",
    checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
        var response domain.ErrorResponse  // ← Expects ErrorResponse
        // ... check error message
    },
},
{
    name: "Missing email",
    checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
        var response domain.ErrorResponse  // ← Different error message
        // ... check "Email is required"
    },
},
Why This Pattern?
Each test case validates different things:

Táblázat
Test Case	Response Type	What to Verify
Valid credentials	LoginResponseDTO	Token exists, user info correct, JWT format
Invalid credentials	ErrorResponse	Error = "Invalid email or password"
Missing email	ErrorResponse	Error contains "Email is required"
Invalid JSON	ErrorResponse	Error = "unexpected EOF"
Internal error	ErrorResponse	Error = "internal server error"
The Loop Runs Shared Logic
The loop handles common operations:

go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // 1. Setup (same for all) ✅
        mockService := new(mocks.MockUserService)
        tt.setupMock(mockService)

        // 2. Create handler (same for all) ✅
        handler := &UserHandlers{...}

        // 3. Execute request (same for all) ✅
        req := httptest.NewRequest(...)
        rr := httptest.NewRecorder()
        handler.Login(rr, req)

        // 4. Check status code (same for all) ✅
        assert.Equal(t, tt.expectedStatus, rr.Code)

        // 5. Custom validation (DIFFERENT for each) ⚠️
        if tt.checkResponse != nil {
            tt.checkResponse(t, rr)  // ← Calls test-specific function
        }

        // 6. Verify mocks (same for all) ✅
        mockService.AssertExpectations(t)
    })
}
Comparison: Without Custom Functions
If we didn't use checkResponse, we'd need ugly switch statements:

go
// ❌ BAD: Messy switch logic
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // ... setup and execute ...

        // Check response based on test name 😱
        switch tt.name {
        case "Valid credentials":
            var response domain.LoginResponseDTO
            json.Unmarshal(rr.Body.Bytes(), &response)
            assert.NotEmpty(t, response.Token)
            assert.Equal(t, int64(1), response.User.ID)
            // ... more checks ...

        case "Invalid credentials":
            var response domain.ErrorResponse
            json.Unmarshal(rr.Body.Bytes(), &response)
            assert.Equal(t, "Invalid email or password", response.Error)

        case "Missing email":
            var response domain.ErrorResponse
            json.Unmarshal(rr.Body.Bytes(), &response)
            assert.Contains(t, response.Error, "Email is required")

        // ... 10 more cases ...
        }
    })
}
Problems:

❌ Huge switch statement
❌ Mixes setup logic with validation
❌ Hard to read
❌ Easy to forget to add new cases
Current Pattern: Clean & Flexible
go
// ✅ GOOD: Each test declares its own validation
{
    name: "Valid credentials",
    checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
        // Validation logic right here, next to test definition
        var response domain.LoginResponseDTO
        json.Unmarshal(rr.Body.Bytes(), &response)
        assert.NotEmpty(t, response.Token)
        assert.Equal(t, int64(1), response.User.ID)
    },
},
Benefits:

✅ Test definition and validation are together
✅ Each test is self-contained
✅ Easy to add new tests with custom validation
✅ No giant switch statement
When You WOULD Extract to Loop
If validation is identical for multiple tests, then extract it:

go
// Example: If all success cases check the same things
func checkSuccessResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedUserID int64) {
    var response domain.LoginResponseDTO
    json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NotEmpty(t, response.Token)
    assert.Equal(t, expectedUserID, response.User.ID)
}

// Then in tests:
{
    name: "Valid credentials",
    checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
        checkSuccessResponse(t, rr, 1)  // Reuse helper
    },
},
Your Login Tests Have Different Validations
Look at the actual differences:

Test 1: Valid Credentials
go
checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
    var response domain.LoginResponseDTO  // ← Different type
    // ... checks token, user.ID, user.Email, user.Name
    assert.Equal(t, int64(1), response.User.ID)  // ← Specific value
}
Test 2: Invalid Credentials
go
checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
    var response domain.ErrorResponse  // ← Different type
    // ... checks error message
    assert.Equal(t, "Invalid email or password", response.Error)  // ← Different check
}
Test 3: Missing Email
go
checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
    var response domain.ErrorResponse  // ← Same type as test 2
    // ... but different error message
    assert.Contains(t, response.Error, "Email is required")  // ← Different message
}
Each needs different validation logic, so each gets its own function!

Summary
Why checkResponse is in each test case:

✅ Each test validates different response types (LoginResponseDTO vs ErrorResponse)
✅ Each test checks different fields/values
✅ Keeps validation logic close to test definition (easier to read)
✅ Avoids messy switch statements in the loop
The loop handles:

✅ Common setup (mock, handler, request)
✅ Common execution (calling handler)
✅ Common assertion (status code)
✅ Then delegates custom validation to each test's checkResponse
This is a clean table-driven test pattern where:

Shared logic → in the loop
Custom logic → in each test case
Does this make sense? It's like having a template (the loop) with customizable parts (checkResponse)! 🎨

Ok now I need a learning workbook about jwt and how to handle each layer
also for interview prep in an md file
Claude Sonnet 4.5
Gondolkodási folyamat


JWT Authentication Learning Workbook & Interview Prep
Complete Guide to JWT Authentication in Go with Clean Architecture

Table of Contents
JWT Fundamentals
Clean Architecture Layers
Implementation Guide
Security Best Practices
Common Pitfalls
Interview Questions & Answers
Practical Exercises
Cheat Sheet
1. JWT Fundamentals
1.1 What is JWT?
JWT (JSON Web Token) is a compact, URL-safe means of representing claims to be transferred between two parties.

Structure:

xxxxx.yyyyy.zzzzz
Three Parts:

Header - Algorithm & token type
Payload - Claims (data)
Signature - Verification
1.2 JWT Structure Deep Dive
Header
json
{
  "alg": "HS256",
  "typ": "JWT"
}
Base64Url encoded: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9

Payload (Claims)
json
{
  "user_id": 123,
  "email": "user@example.com",
  "name": "John Doe",
  "exp": 1735689600,
  "iat": 1735603200
}
Base64Url encoded: eyJ1c2VyX2lkIjoxMjMsImVtYWlsIjoidXNlckBleGFtcGxlLmNvbSJ9

Signature
scss
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret
)
Result: SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

1.3 Standard JWT Claims
Táblázat
Claim	Full Name	Type	Description
iss	Issuer	String	Who created the token
sub	Subject	String	Who the token is about
aud	Audience	String	Who should accept the token
exp	Expiration Time	Number	When token expires (Unix timestamp)
nbf	Not Before	Number	Token not valid before this time
iat	Issued At	Number	When token was created
jti	JWT ID	String	Unique identifier for the token
Custom Claims:

user_id - Application-specific user identifier
email - User's email address
role - User's role (admin, user, etc.)
Any other application data
1.4 JWT vs Session Authentication
Táblázat
Feature	JWT	Session
Storage	Client-side (localStorage, cookie)	Server-side (database, Redis)
Stateless	✅ Yes (server doesn't store anything)	❌ No (server stores session data)
Scalability	✅ Easy (no shared state)	⚠️ Harder (need session store)
Revocation	❌ Difficult (token is valid until expiry)	✅ Easy (delete session)
Size	⚠️ Larger (sent with every request)	✅ Small (just session ID)
Security	⚠️ XSS risk if in localStorage	✅ Better if HttpOnly cookie
1.5 How JWT Authentication Works
sql
┌─────────┐                                      ┌─────────┐
│ Client  │                                      │ Server  │
└────┬────┘                                      └────┬────┘
     │                                                │
     │  1. POST /login                               │
     │    { email, password }                        │
     │──────────────────────────────────────────────>│
     │                                                │
     │                          2. Verify credentials│
     │                             Hash password     │
     │                             Compare with DB   │
     │                                                │
     │                          3. Generate JWT      │
     │                             Create claims     │
     │                             Sign with secret  │
     │                                                │
     │  4. Response: { token: "eyJ..." }             │
     │<──────────────────────────────────────────────│
     │                                                │
     │  5. Store token                               │
     │     (localStorage/cookie)                     │
     │                                                │
     │  6. GET /protected-resource                   │
     │     Authorization: Bearer eyJ...              │
     │──────────────────────────────────────────────>│
     │                                                │
     │                          7. Extract token     │
     │                             Verify signature  │
     │                             Check expiration  │
     │                             Extract user info │
     │                                                │
     │  8. Response: { data: [...] }                 │
     │<──────────────────────────────────────────────│
     │                                                │
2. Clean Architecture Layers
2.1 Layer Overview
scss
┌─────────────────────────────────────────────────────────┐
│                    Delivery Layer                        │
│  (HTTP Handlers, Middleware, Request/Response DTOs)      │
│  • Handles HTTP requests                                 │
│  • Validates input                                       │
│  • Generates JWT tokens                                  │
│  • Returns HTTP responses                                │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                    Service Layer                         │
│  (Business Logic, Use Cases)                             │
│  • User registration                                     │
│  • User authentication (password verification)           │
│  • Business rules validation                             │
│  • NO JWT knowledge (that's delivery concern)            │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                    Domain Layer                          │
│  (Entities, Interfaces, Business Errors)                 │
│  • User entity                                           │
│  • Service interfaces                                    │
│  • Domain errors                                         │
│  • NO dependencies on other layers                       │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                  Repository Layer                        │
│  (Database Access, Data Persistence)                     │
│  • CRUD operations                                       │
│  • Database queries                                      │
│  • Data mapping                                          │
└─────────────────────────────────────────────────────────┘
2.2 Responsibility Matrix
Táblázat
Layer	JWT Responsibility	Why?
Domain	None	Pure business logic, no infrastructure
Service	None	Business logic only, infrastructure-agnostic
Delivery	Generate & Validate	HTTP/REST concern, not business logic
Repository	None	Database concern only
Key Principle: JWT is a delivery mechanism, not business logic!

2.3 Data Flow
Login Flow
go
// 1. CLIENT REQUEST
POST /login
{
  "email": "user@example.com",
  "password": "Password123"
}

// 2. DELIVERY LAYER (Handler)
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
    // Decode JSON
    var req LoginDTO
    json.NewDecoder(r.Body).Decode(&req)

    // Validate input
    validate.Struct(req)

    // Call service (business logic)
    user, err := h.Service.Login(ctx, req.Email, req.Password)

    // Generate JWT (delivery concern)
    claims := auth.NewUserClaims(user, 24*time.Hour)
    tokenString, _ := h.TokenAuth.Encode(claims.ToMap())

    // Return response
    WriteJSON(w, 200, LoginResponseDTO{
        Token: tokenString,
        User:  user,
    })
}

// 3. SERVICE LAYER (Business Logic)
func (s *UserService) Login(ctx context.Context, email, password string) (*User, error) {
    // Fetch user from database
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    // Verify password (business logic)
    if !bcrypt.CompareHashAndPassword(user.Password, []byte(password)) {
        return nil, ErrInvalidCredentials
    }

    // Return user (NO JWT here!)
    return user, nil
}

// 4. REPOSITORY LAYER (Database)
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := r.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email).
        Scan(&user.ID, &user.Name, &user.Email, &user.Password)
    return &user, err
}
Protected Route Flow
go
// 1. CLIENT REQUEST
GET /todos
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

// 2. MIDDLEWARE CHAIN
// 2a. Verifier (jwtauth library)
func Verifier(tokenAuth *JWTAuth) func(http.Handler) http.Handler {
    // Extract token from header
    // Verify signature
    // Validate expiration
    // Store token in context
}

// 2b. Authenticator (custom)
func Authenticator(next http.Handler) http.Handler {
    // Get token from context
    // Check if validation passed
    // Validate claims structure
    // Reject if invalid
}

// 2c. UserContext (custom)
func UserContext(next http.Handler) http.Handler {
    // Extract claims from token
    // Create UserContext
    // Add to request context
}

// 3. HANDLER (uses authenticated user)
func (h *TodoHandlers) ListTodos(w http.ResponseWriter, r *http.Request) {
    // Get user from context (set by middleware)
    userCtx, _ := auth.FromContext(r.Context())

    // Call service with user ID
    todos, err := h.Service.ListTodos(r.Context(), userCtx.ID)

    // Return response
    WriteJSON(w, 200, todos)
}

// 4. SERVICE LAYER (uses user ID)
func (s *TodoService) ListTodos(ctx context.Context, userID int64) ([]*Todo, error) {
    // Business logic: user can only see their own todos
    return s.Repo.GetByUserID(ctx, userID)
}
3. Implementation Guide
3.1 Domain Layer
Purpose: Define business entities and rules

domain/user.go

go
package domain

import (
	"context"
	"errors"
)

// User - Core business entity
type User struct {
	ID       int64
	Name     string
	Email    string
	Password string // Hashed
}

// Business validation
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if len(u.Password) < 8 {
		return ErrWeakPassword
	}
	return nil
}

// UserService - Business logic interface
type UserService interface {
	CreateUser(ctx context.Context, name, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	GetUser(ctx context.Context, id int64) (*User, error)
}
domain/errors.go

go
package domain

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicate          = errors.New("email already exists")
)
Key Points:

✅ No JWT knowledge
✅ No HTTP knowledge
✅ Pure business logic
✅ Infrastructure-agnostic
3.2 Service Layer
Purpose: Implement business logic

service/user_service.go

go
package service

import (
	"context"

	"github.com/macesz/todo-go/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// CreateUser - Register new user
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	// Business rule: check if email exists
	existing, err := s.Repo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, domain.ErrDuplicate
	}

	// Business rule: hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Validate business rules
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Persist
	if err := s.Repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login - Authenticate user
func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	// Fetch user
	user, err := s.Repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Business rule: verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Return authenticated user (NO JWT here!)
	return user, nil
}

// GetUser - Fetch user by ID
func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return s.Repo.GetByID(ctx, id)
}
Key Points:

✅ No JWT knowledge
✅ No HTTP knowledge
✅ Only business logic (password hashing, validation)
✅ Returns domain entities
3.3 Delivery Layer - Auth Infrastructure
Purpose: JWT claims structure and helpers

delivery/web/auth/claims.go

go
package auth

import (
	"errors"
	"time"

	"github.com/macesz/todo-go/domain"
)

// UserClaims - JWT payload structure
type UserClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Exp    int64  `json:"exp"`
}

// NewUserClaims - Convert domain.User to JWT claims
func NewUserClaims(user *domain.User, expiresIn time.Duration) UserClaims {
	return UserClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Exp:    time.Now().Add(expiresIn).Unix(),
	}
}

// ToMap - Convert to map for jwtauth library
func (c UserClaims) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id": c.UserID,
		"email":   c.Email,
		"name":    c.Name,
		"exp":     c.Exp,
	}
}

// FromTokenClaims - Extract and validate claims from JWT
// IMPORTANT: JWT stores numbers as float64!
func FromTokenClaims(claims map[string]interface{}) (*UserClaims, error) {
	// Extract user_id (stored as float64 in JWT)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("missing or invalid user_id claim")
	}

	// Extract email
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("missing or invalid email claim")
	}

	// Extract name
	name, ok := claims["name"].(string)
	if !ok {
		return nil, errors.New("missing or invalid name claim")
	}

	return &UserClaims{
		UserID: int64(userID), // Convert float64 to int64
		Email:  email,
		Name:   name,
	}, nil
}
delivery/web/auth/context.go

go
package auth

import "context"

type contextKey string

const userContextKey contextKey = "userContext"

// UserContext - User information in request context
type UserContext struct {
	ID    int64
	Email string
	Name  string
}

// NewUserContext - Create from JWT claims
func NewUserContext(claims *UserClaims) *UserContext {
	return &UserContext{
		ID:    claims.UserID,
		Email: claims.Email,
		Name:  claims.Name,
	}
}

// AddToContext - Store in request context
func (u *UserContext) AddToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// FromContext - Retrieve from request context
func FromContext(ctx context.Context) (*UserContext, bool) {
	userCtx, ok := ctx.Value(userContextKey).(*UserContext)
	return userCtx, ok
}
3.4 Delivery Layer - Middlewares
Purpose: Extract, validate, and process JWT tokens

delivery/web/middlewares/auth.go

go
package middlewares

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/delivery/web/utils"
)

// Authenticator - Validates token and claims structure
// Runs AFTER jwtauth.Verifier
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from context (set by Verifier)
		token, _, err := jwtauth.FromContext(r.Context())

		// Check if Verifier found issues
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized - " + err.Error(),
			})
			return
		}

		if token == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Missing authentication token",
			})
			return
		}

		// Validate claims structure
		privateClaims := token.PrivateClaims()
		_, err = auth.FromTokenClaims(privateClaims)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid token claims - " + err.Error(),
			})
			return
		}

		// Token is valid, continue
		next.ServeHTTP(w, r)
	})
}

// UserContext - Extracts user info and adds to request context
// Runs AFTER Authenticator
func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from context
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil || token == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized",
			})
			return
		}

		// Extract claims
		privateClaims := token.PrivateClaims()
		claims, err := auth.FromTokenClaims(privateClaims)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid token claims",
			})
			return
		}

		// Create UserContext and add to request context
		userContext := auth.NewUserContext(claims)
		ctx := userContext.AddToContext(r.Context())

		// Pass updated context to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
Middleware Chain:

sql
Request → Verifier → Authenticator → UserContext → Handler
          ↓          ↓               ↓              ↓
          Extract    Check result    Extract user   Use user info
          Validate   Validate claims Add to context
3.5 Delivery Layer - Handlers
delivery/web/handlers/user.go

go
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/delivery/web/utils"
	"github.com/macesz/todo-go/domain"
)

type UserHandlers struct {
	Service   domain.UserService
	TokenAuth *jwtauth.JWTAuth
}

// Login - Authenticate and return JWT
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parse request
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Authenticate (business logic)
	user, err := h.Service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid email or password",
			})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
		return
	}

	// Generate JWT (delivery concern)
	claims := auth.NewUserClaims(user, 24*time.Hour)
	_, tokenString, err := h.TokenAuth.Encode(claims.ToMap())
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
		return
	}

	// Return response
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// GetUser - Get current user info (protected route)
func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userCtx, ok := auth.FromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Fetch user details
	user, err := h.Service.GetUser(r.Context(), userCtx.ID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}
3.6 Router Setup
delivery/web/server.go

go
package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/handlers"
	"github.com/macesz/todo-go/delivery/web/middlewares"
)

func SetupRouter(userHandlers *handlers.UserHandlers, todoHandlers *handlers.TodoHandlers, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	// JWT setup
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
	userHandlers.TokenAuth = tokenAuth

	// Global middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// PUBLIC ROUTES
	r.Group(func(r chi.Router) {
		r.Post("/register", userHandlers.CreateUser)
		r.Post("/login", userHandlers.Login)
	})

	// PROTECTED ROUTES
	r.Group(func(r chi.Router) {
		// JWT Middleware Chain
		r.Use(jwtauth.Verifier(tokenAuth))  // 1. Extract & validate
		r.Use(middlewares.Authenticator)     // 2. Check validation result
		r.Use(middlewares.UserContext)       // 3. Add user to context

		// User routes
		r.Get("/user", userHandlers.GetUser)

		// Todo routes
		r.Get("/todos", todoHandlers.ListTodos)
		r.Post("/todos", todoHandlers.CreateTodo)
		r.Get("/todos/{id}", todoHandlers.GetTodo)
		r.Put("/todos/{id}", todoHandlers.UpdateTodo)
		r.Delete("/todos/{id}", todoHandlers.DeleteTodo)
	})

	return r
}
4. Security Best Practices
4.1 Secret Key Management
❌ NEVER DO THIS:

go
// Hardcoded secret in code
tokenAuth := jwtauth.New("HS256", []byte("secret123"), nil)
✅ DO THIS:

go
// Load from environment variable
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
Generate strong secret:

bash
# Generate 256-bit (32 bytes) random secret
openssl rand -base64 32

# Result: 8xK7vZ2qT9nR4wL6pY3mJ5hF8dS1cX0bV7nQ9rP4kM2e=
Minimum Requirements:

✅ At least 256 bits (32 bytes) for HS256
✅ Cryptographically random
✅ Never committed to version control
✅ Different for each environment (dev, staging, prod)
4.2 Token Expiration
Short-lived tokens are more secure:

go
// ❌ Too long - security risk if token is stolen
claims := auth.NewUserClaims(user, 30*24*time.Hour) // 30 days

// ✅ Short-lived access token
claims := auth.NewUserClaims(user, 15*time.Minute) // 15 minutes

// ✅ Reasonable for development
claims := auth.NewUserClaims(user, 1*time.Hour) // 1 hour

// ✅ Balance between security and UX
claims := auth.NewUserClaims(user, 24*time.Hour) // 24 hours
Best Practice: Access + Refresh Tokens

go
// Access token - short-lived, used for API requests
accessToken := auth.NewUserClaims(user, 15*time.Minute)

// Refresh token - long-lived, used to get new access token
refreshToken := auth.NewRefreshToken(user, 7*24*time.Hour)
4.3 HTTPS in Production
❌ NEVER use HTTP in production:

go
// Development only
http.ListenAndServe(":3000", router)
✅ Always use HTTPS:

go
// Production
http.ListenAndServeTLS(":443", "cert.pem", "key.pem", router)
Why? JWT is sent in Authorization header. HTTP is unencrypted → JWT can be intercepted!

4.4 Token Storage
Client-side options:

Táblázat
Storage	Security	Pros	Cons
localStorage	❌ Low	Easy to use	Vulnerable to XSS
sessionStorage	❌ Low	Cleared on tab close	Vulnerable to XSS
Memory (JS variable)	⚠️ Medium	Not accessible via document	Lost on refresh
HttpOnly Cookie	✅ High	Not accessible via JS	Requires CSRF protection
Secure Cookie	✅ Highest	HTTPS only + HttpOnly	Requires proper setup
Best Practice: HttpOnly + Secure Cookie

go
// Server sets JWT in cookie
http.SetCookie(w, &http.Cookie{
    Name:     "jwt",
    Value:    tokenString,
    HttpOnly: true,  // Not accessible via JavaScript
    Secure:   true,  // HTTPS only
    SameSite: http.SameSiteStrictMode, // CSRF protection
    MaxAge:   86400, // 24 hours
    Path:     "/",
})
4.5 What NOT to Store in JWT
❌ NEVER store sensitive data:

go
// ❌ BAD - Anyone can decode JWT!
claims := map[string]interface{}{
    "user_id":     user.ID,
    "password":    user.Password,      // NEVER!
    "credit_card": user.CreditCard,    // NEVER!
    "ssn":         user.SSN,           // NEVER!
    "api_key":     user.APIKey,        // NEVER!
}
✅ Only store non-sensitive identifiers:

go
// ✅ GOOD - Public or identifiable data only
claims := map[string]interface{}{
    "user_id": user.ID,
    "email":   user.Email,
    "name":    user.Name,
    "role":    user.Role,
}
Remember: JWT is signed, not encrypted. Anyone can decode and read it!

bash
# Anyone can decode JWT
echo "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ..." | base64 -d
# Reveals: {"user_id":1,"password":"secret123"}
4.6 Token Revocation
Problem: JWT is valid until expiration, even if user logs out or password changes.

Solutions:

1. Blacklist (Token Revocation List)
go
// Store revoked tokens in Redis/Database
func (s *AuthService) RevokeToken(tokenID string, expiry time.Time) error {
    return s.Redis.Set("revoked:"+tokenID, "1", time.Until(expiry))
}

// Check in middleware
func CheckRevoked(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, _ := jwtauth.FromContext(r.Context())
        tokenID := token.Claims["jti"].(string)

        if s.Redis.Exists("revoked:" + tokenID) {
            http.Error(w, "Token revoked", 401)
            return
        }

        next.ServeHTTP(w, r)
    })
}
2. Short Expiration + Refresh Tokens
go
// Short-lived access token (15 minutes)
accessToken := generateToken(user, 15*time.Minute)

// Long-lived refresh token (7 days) stored in database
refreshToken := generateRefreshToken(user, 7*24*time.Hour)
s.DB.StoreRefreshToken(user.ID, refreshToken)
3. Version in Claims
go
// Add user version to claims
claims := map[string]interface{}{
    "user_id": user.ID,
    "version": user.TokenVersion, // Incremented on password change
}

// Verify version in middleware
func VerifyVersion(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, _ := jwtauth.FromContext(r.Context())
        claimVersion := token.Claims["version"].(float64)

        user := s.DB.GetUser(userID)
        if int(claimVersion) != user.TokenVersion {
            http.Error(w, "Token invalid", 401)
            return
        }

        next.ServeHTTP(w, r)
    })
}
5. Common Pitfalls
5.1 Type Assertion Issues
Problem: JWT stores numbers as float64, not int64

❌ WRONG:

go
userID, ok := claims["user_id"].(int64)
// ok is always false! JWT uses float64
✅ CORRECT:

go
userID, ok := claims["user_id"].(float64)
if !ok {
    return errors.New("invalid user_id")
}
intUserID := int64(userID) // Convert to int64
5.2 Forgetting to Return After Error
❌ WRONG:

go
if err != nil {
    http.Error(w, "Unauthorized", 401)
    // Missing return - handler continues!
}
next.ServeHTTP(w, r) // Still executes!
✅ CORRECT:

go
if err != nil {
    http.Error(w, "Unauthorized", 401)
    return // Stop here
}
next.ServeHTTP(w, r)
5.3 Not Using WithContext
❌ WRONG:

go
func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := addUserToContext(r.Context())
        next.ServeHTTP(w, r) // Old context!
    })
}
✅ CORRECT:

go
func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := addUserToContext(r.Context())
        next.ServeHTTP(w, r.WithContext(ctx)) // Updated context
    })
}
5.4 Validating Token in Wrong Layer
❌ WRONG:

go
// Service layer validating JWT
func (s *UserService) GetUser(ctx context.Context, tokenString string) (*User, error) {
    // Parsing JWT in business logic!
    token, err := jwt.Parse(tokenString, ...)
    userID := token.Claims["user_id"]

    return s.Repo.GetByID(ctx, userID)
}
✅ CORRECT:

go
// Middleware validates JWT, service uses user ID
func UserContext(next http.Handler) http.Handler {
    // Validate JWT here
}

func (s *UserService) GetUser(ctx context.Context, userID int64) (*User, error) {
    // Pure business logic, no JWT knowledge
    return s.Repo.GetByID(ctx, userID)
}
5.5 Exposing Password in Response
❌ WRONG:

go
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
    user, _ := h.Service.Login(...)

    // Returns hashed password!
    json.NewEncoder(w).Encode(user)
}
✅ CORRECT:

go
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
    user, _ := h.Service.Login(...)

    // Use DTO without password
    response := UserResponseDTO{
        ID:    user.ID,
        Email: user.Email,
        Name:  user.Name,
        // Password excluded
    }

    json.NewEncoder(w).Encode(response)
}
6. Interview Questions & Answers
6.1 Basic Questions
Q1: What is JWT and how does it work?
Answer:

JWT (JSON Web Token) is a compact, URL-safe token format for securely transmitting information between parties. It consists of three parts:

Header - Algorithm and token type
Payload - Claims (data)
Signature - Verification using a secret key
How it works:

User logs in with credentials
Server verifies credentials
Server creates JWT with user information (claims)
Server signs JWT with secret key
Server returns JWT to client
Client includes JWT in Authorization header for subsequent requests
Server verifies JWT signature and extracts user information
Server processes request with authenticated user context
Key advantage: Stateless - server doesn't need to store sessions.

Q2: What's the difference between authentication and authorization?
Answer:

Authentication: Verifying who you are (identity)
Example: Login with email/password proves you're John Doe
Authorization: Verifying what you can do (permissions)
Example: Can John Doe delete this todo? Can he access admin panel?
In JWT context:

JWT helps with authentication by proving identity
JWT can include role/permissions for authorization
go
// Authentication (JWT middleware)
userCtx := auth.FromContext(r.Context())
if userCtx == nil {
    return errors.New("not authenticated")
}

// Authorization (business logic)
if userCtx.Role != "admin" {
    return errors.New("not authorized")
}
Q3: Why use JWT instead of sessions?
Answer:

Táblázat
Feature	JWT	Session
Stateless	Yes - server doesn't store anything	No - server stores session data
Scalability	Easy - no shared state between servers	Harder - need shared session store (Redis)
Performance	Faster - no database lookup	Slower - need to fetch session
Revocation	Difficult - token valid until expiry	Easy - delete session from store
Size	Larger - sent with every request	Smaller - just session ID
When to use JWT:

Microservices architecture
Mobile apps
Stateless APIs
Distributed systems
When to use sessions:

Traditional web apps
Need immediate revocation
Smaller data transfer requirements
Q4: What are JWT claims?
Answer:

Claims are statements about an entity (typically, the user) and additional data. There are three types:

1. Registered Claims (Standard):

iss (issuer) - Who created the token
sub (subject) - Who the token is about
aud (audience) - Who should accept the token
exp (expiration) - When token expires
iat (issued at) - When token was created
nbf (not before) - Token not valid before this time
2. Public Claims:

Defined in IANA JWT Registry
Collision-resistant names
3. Private Claims (Custom):

Your application-specific data
Examples: user_id, email, role
go
claims := map[string]interface{}{
    // Registered
    "exp": time.Now().Add(24 * time.Hour).Unix(),
    "iat": time.Now().Unix(),

    // Custom
    "user_id": 123,
    "email":   "user@example.com",
    "role":    "admin",
}
6.2 Intermediate Questions
Q5: How do you handle JWT expiration?
Answer:

Three strategies:

1. Simple Approach - Re-login:

go
// Token expires, user must login again
if token.IsExpired() {
    return errors.New("Token expired, please login again")
}
2. Refresh Token Pattern (Most Common):

go
// Two tokens:
// - Short-lived access token (15 min)
// - Long-lived refresh token (7 days)

type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// When access token expires, use refresh token to get new one
func (h *AuthHandlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
    refreshToken := extractToken(r)

    // Verify refresh token (stored in database)
    user, err := h.Service.VerifyRefreshToken(refreshToken)
    if err != nil {
        return errors.New("Invalid refresh token")
    }

    // Generate new access token
    newAccessToken := generateAccessToken(user)

    json.NewEncoder(w).Encode(map[string]string{
        "access_token": newAccessToken,
    })
}
3. Sliding Window:

go
// If token is close to expiry, automatically refresh
func RefreshIfNeeded(token *jwt.Token) string {
    exp := token.Claims["exp"].(float64)
    expiresAt := time.Unix(int64(exp), 0)

    // If expires in less than 5 minutes
    if time.Until(expiresAt) < 5*time.Minute {
        return generateNewToken(token.Claims)
    }

    return token.Raw
}
Best Practice: Use refresh tokens with short-lived access tokens.

Q6: How do you secure JWT tokens?
Answer:

1. Use Strong Secret Keys:

go
// Minimum 256 bits for HS256
secret := make([]byte, 32)
rand.Read(secret)
2. Use HTTPS Only:

go
// Never send JWT over HTTP
http.ListenAndServeTLS(":443", "cert.pem", "key.pem", router)
3. Set Appropriate Expiration:

go
// Short-lived tokens
claims := NewUserClaims(user, 15*time.Minute)
4. Store in HttpOnly Cookies:

go
http.SetCookie(w, &http.Cookie{
    Name:     "jwt",
    Value:    token,
    HttpOnly: true,  // No JS access
    Secure:   true,  // HTTPS only
    SameSite: http.SameSiteStrictMode,
})
5. Don't Store Sensitive Data:

go
// ❌ Don't include passwords, credit cards, etc.
// ✅ Only non-sensitive identifiers
claims := map[string]interface{}{
    "user_id": user.ID,
    "email":   user.Email,
}
6. Validate All Claims:

go
// Check signature, expiration, issuer, audience
if !token.Valid() {
    return errors.New("Invalid token")
}
7. Implement Token Revocation:

go
// Blacklist or refresh token rotation
Q7: What's the difference between symmetric (HS256) and asymmetric (RS256) signing?
Answer:

Táblázat
Feature	HS256 (Symmetric)	RS256 (Asymmetric)
Algorithm	HMAC with SHA-256	RSA with SHA-256
Keys	Same key for sign and verify	Private key to sign, public key to verify
Speed	✅ Faster	⚠️ Slower
Key Distribution	⚠️ Secret must be shared	✅ Public key can be shared
Use Case	Single server/service	Multiple services, public verification
HS256 (Symmetric):

go
// Same secret for both operations
tokenAuth := jwtauth.New("HS256", []byte("secret-key"), nil)

// Server signs
token, _ := tokenAuth.Encode(claims)

// Server verifies (same secret)
verified, _ := tokenAuth.Decode(token)
RS256 (Asymmetric):

go
// Private key to sign
privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
tokenAuth := jwtauth.New("RS256", privateKey, nil)

// Server signs with private key
token, _ := tokenAuth.Encode(claims)

// Anyone can verify with public key
publicKey := &privateKey.PublicKey
verified := jwt.Verify(token, publicKey)
When to use:

HS256: Single server, microservices with shared secret
RS256: Multiple services, third-party verification, public APIs
Q8: How would you implement role-based access control (RBAC) with JWT?
Answer:

Step 1: Add role to JWT claims

go
type UserClaims struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"` // "admin", "user", "moderator"
}

claims := UserClaims{
    UserID: user.ID,
    Email:  user.Email,
    Role:   user.Role,
}
Step 2: Create authorization middleware

go
func RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get user from context
            userCtx, _ := auth.FromContext(r.Context())

            // Check if user has required role
            hasRole := false
            for _, role := range roles {
                if userCtx.Role == role {
                    hasRole = true
                    break
                }
            }

            if !hasRole {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
Step 3: Apply to routes

go
// Admin only routes
r.Group(func(r chi.Router) {
    r.Use(jwtauth.Verifier(tokenAuth))
    r.Use(middlewares.Authenticator)
    r.Use(middlewares.UserContext)
    r.Use(middlewares.RequireRole("admin")) // Admin only

    r.Delete("/users/{id}", adminHandlers.DeleteUser)
    r.Get("/admin/stats", adminHandlers.GetStats)
})

// Admin or moderator routes
r.Group(func(r chi.Router) {
    r.Use(jwtauth.Verifier(tokenAuth))
    r.Use(middlewares.Authenticator)
    r.Use(middlewares.UserContext)
    r.Use(middlewares.RequireRole("admin", "moderator"))

    r.Delete("/posts/{id}", modHandlers.DeletePost)
})
Advanced: Permission-based (more granular)

go
type UserClaims struct {
    UserID      int64    `json:"user_id"`
    Permissions []string `json:"permissions"`
    // ["posts:read", "posts:write", "users:delete"]
}

func RequirePermission(permission string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userCtx, _ := auth.FromContext(r.Context())

            hasPermission := false
            for _, p := range userCtx.Permissions {
                if p == permission {
                    hasPermission = true
                    break
                }
            }

            if !hasPermission {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

// Usage
r.Delete("/posts/{id}", RequirePermission("posts:delete")(handler))
6.3 Advanced Questions
Q9: Explain the middleware chain in your JWT implementation
Answer:

Three-stage middleware chain:

css
Request → [1] Verifier → [2] Authenticator → [3] UserContext → Handler
[1] jwtauth.Verifier (Library-provided):

go
// What it does:
// - Extracts token from Authorization header
// - Verifies signature with secret key
// - Validates expiration (exp)
// - Validates not-before (nbf)
// - Validates issued-at (iat)
// - Stores token + error in context

// IMPORTANT: Even if invalid, it continues (stores error in context)
[2] Authenticator (Custom):

go
// What it does:
// - Retrieves token from context (set by Verifier)
// - Checks if Verifier found any errors
// - Validates custom claims structure (user_id, email, name exist and correct types)
// - Rejects request if invalid (returns 401)
// - Continues if valid

func Authenticator(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, _, err := jwtauth.FromContext(r.Context())

        if err != nil || token == nil {
            http.Error(w, "Unauthorized", 401)
            return
        }

        // Validate claims structure
        _, err = auth.FromTokenClaims(token.PrivateClaims())
        if err != nil {
            http.Error(w, "Invalid claims", 401)
            return
        }

        next.ServeHTTP(w, r)
    })
}
[3] UserContext (Custom):

go
// What it does:
// - Extracts claims from token
// - Creates UserContext struct
// - Adds UserContext to request context
// - Handlers can now access user info

func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, _, _ := jwtauth.FromContext(r.Context())

        claims, _ := auth.FromTokenClaims(token.PrivateClaims())

        userContext := auth.NewUserContext(claims)
        ctx := userContext.AddToContext(r.Context())

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
Why three stages?

Single Responsibility: Each middleware has one job
Separation of Concerns: Crypto validation vs business validation vs data extraction
Flexibility: Can customize Authenticator and UserContext without touching Verifier
Testability: Each can be tested independently
Q10: How do you test JWT handlers?
Answer:

Three levels of testing:

1. Unit Tests (Mock Service):

go
func TestLogin(t *testing.T) {
    // Setup
    mockService := new(mocks.MockUserService)
    mockService.On("Login",
        mock.Anything,
        "test@example.com",
        "Password123",
    ).Return(&domain.User{
        ID:    1,
        Name:  "Test User",
        Email: "test@example.com",
    }, nil)

    tokenAuth := jwtauth.New("HS256", []byte("test-secret"), nil)

    handler := &UserHandlers{
        Service:   mockService,
        TokenAuth: tokenAuth,
    }

    // Execute
    body := `{"email":"test@example.com","password":"Password123"}`
    req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
    rr := httptest.NewRecorder()

    handler.Login(rr, req)

    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)

    var response LoginResponseDTO
    json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NotEmpty(t, response.Token)
    assert.Equal(t, int64(1), response.User.ID)

    mockService.AssertExpectations(t)
}
2. Middleware Tests:

go
func TestAuthenticator(t *testing.T) {
    tokenAuth := jwtauth.New("HS256", []byte("test-secret"), nil)

    // Create valid token
    claims := map[string]interface{}{
        "user_id": float64(1),
        "email":   "test@example.com",
        "name":    "Test User",
        "exp":     time.Now().Add(1 * time.Hour).Unix(),
    }
    _, tokenString, _ := tokenAuth.Encode(claims)

    // Test with valid token
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+tokenString)

    // Apply middleware chain
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    chain := jwtauth.Verifier(tokenAuth)(
        middlewares.Authenticator(handler),
    )

    rr := httptest.NewRecorder()
    chain.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
}
3. Integration Tests:

go
func TestLogin_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Setup real dependencies
    db := setupTestDB(t)
    defer db.Close()

    repo := repository.NewUserRepository(db)
    service := service.NewUserService(repo)
    tokenAuth := jwtauth.New("HS256", []byte("test-secret"), nil)

    handler := &UserHandlers{
        Service:   service,
        TokenAuth: tokenAuth,
    }

    // Create test user in DB
    hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
    db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
        "Test User", "test@example.com", string(hashedPwd))

    // Test login
    body := `{"email":"test@example.com","password":"Password123"}`
    req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
    rr := httptest.NewRecorder()

    handler.Login(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var response LoginResponseDTO
    json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NotEmpty(t, response.Token)
}
Q11: How would you implement JWT refresh tokens?
Answer:

Complete implementation:

Step 1: Add refresh token to domain

go
package domain

type RefreshToken struct {
    Token     string
    UserID    int64
    ExpiresAt time.Time
    CreatedAt time.Time
}

type AuthService interface {
    Login(ctx context.Context, email, password string) (*User, string, string, error)
    // Returns: user, accessToken, refreshToken, error

    RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
    // Returns: new accessToken, error

    RevokeRefreshToken(ctx context.Context, refreshToken string) error
}
Step 2: Implement service

go
func (s *AuthService) Login(ctx context.Context, email, password string) (*User, string, string, error) {
    // Authenticate user
    user, err := s.UserService.Login(ctx, email, password)
    if err != nil {
        return nil, "", "", err
    }

    // Generate short-lived access token (15 minutes)
    accessClaims := auth.NewUserClaims(user, 15*time.Minute)
    _, accessToken, _ := s.TokenAuth.Encode(accessClaims.ToMap())

    // Generate long-lived refresh token (7 days)
    refreshToken := generateRandomToken(32)
    expiresAt := time.Now().Add(7 * 24 * time.Hour)

    // Store refresh token in database
    err = s.Repo.StoreRefreshToken(ctx, &RefreshToken{
        Token:     refreshToken,
        UserID:    user.ID,
        ExpiresAt: expiresAt,
    })

    return user, accessToken, refreshToken, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
    // Verify refresh token exists and is valid
    rt, err := s.Repo.GetRefreshToken(ctx, refreshToken)
    if err != nil {
        return "", errors.New("invalid refresh token")
    }

    // Check expiration
    if time.Now().After(rt.ExpiresAt) {
        s.Repo.DeleteRefreshToken(ctx, refreshToken)
        return "", errors.New("refresh token expired")
    }

    // Get user
    user, err := s.UserService.GetUser(ctx, rt.UserID)
    if err != nil {
        return "", err
    }

    // Generate new access token
    accessClaims := auth.NewUserClaims(user, 15*time.Minute)
    _, accessToken, _ := s.TokenAuth.Encode(accessClaims.ToMap())

    return accessToken, nil
}
Step 3: Add handlers

go
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
    // ... parse request ...

    user, accessToken, refreshToken, err := h.Service.Login(ctx, req.Email, req.Password)

    // Store refresh token in HttpOnly cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        MaxAge:   7 * 24 * 60 * 60, // 7 days
        Path:     "/auth/refresh",
    })

    // Return access token in response
    utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
        "access_token": accessToken,
        "user":         user,
    })
}

func (h *AuthHandlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
    // Get refresh token from cookie
    cookie, err := r.Cookie("refresh_token")
    if err != nil {
        http.Error(w, "No refresh token", http.StatusUnauthorized)
        return
    }

    // Generate new access token
    accessToken, err := h.Service.RefreshAccessToken(r.Context(), cookie.Value)
    if err != nil {
        http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
        return
    }

    utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
        "access_token": accessToken,
    })
}
Step 4: Setup routes

go
// Public routes
r.Post("/login", authHandlers.Login)
r.Post("/auth/refresh", authHandlers.RefreshToken)
r.Post("/logout", authHandlers.Logout)

// Protected routes (access token required)
r.Group(func(r chi.Router) {
    r.Use(jwtauth.Verifier(tokenAuth
thanks, but i need .md file tp put it inmy repo for learning
Claude Sonnet 4.5
Gondolkodási folyamat



JWT Authentication Learning Workbook
Complete Guide to JWT Authentication in Go with Clean Architecture

📚 Table of Contents
1. JWT Fundamentals
1.1 What is JWT?
1.2 JWT Structure
1.3 Standard JWT Claims
1.4 JWT vs Session Authentication
1.5 How JWT Authentication Works
2. Clean Architecture Layers
2.1 Layer Overview
2.2 Responsibility Matrix
2.3 Data Flow
3. Password Hashing with bcrypt
3.1 Why Hash Passwords?
3.2 bcrypt Implementation
3.3 Common Mistakes
4. Implementation Guide
4.1 Domain Layer
4.2 Service Layer
4.3 Delivery Layer - Auth Infrastructure
4.4 Delivery Layer - Middleware
4.5 Delivery Layer - Handlers
4.6 Router Setup
5. Security Best Practices
5.1 Secret Key Management
5.2 Token Expiration
5.3 HTTPS in Production
5.4 Token Storage
5.5 What NOT to Store in JWT
5.6 Token Revocation
6. Common Pitfalls
7. Testing JWT Handlers
8. Interview Questions & Answers
9. Practical Exercises
10. Quick Reference
1. JWT Fundamentals
1.1 What is JWT?
JWT (JSON Web Token) is a compact, URL-safe means of representing claims to be transferred between two parties.

Key Characteristics:

✅ Self-contained - Contains all information about the user
✅ Stateless - Server doesn't need to store session data
✅ Signed - Can verify authenticity and integrity
❌ Not encrypted - Anyone can decode and read it
1.2 JWT Structure
A JWT consists of three Base64Url-encoded parts separated by dots:

xxxxx.yyyyy.zzzzz
Example:

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImVtYWlsIjoidXNlckBleGFtcGxlLmNvbSIsImV4cCI6MTczNTY4OTYwMH0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
Part 1: Header
Contains metadata about the token

i need a readme.md file please that o can download
Claude Sonnet 4.5
Gondolkodási folyamat



I'll create a complete markdown file for you. Here it is:

markdown
# JWT Authentication Learning Workbook & Interview Prep Guide

**Complete Guide to JWT Authentication in Go with Clean Architecture**

> 📚 **Purpose:** Master JWT authentication implementation in Go, understand clean architecture patterns, and prepare for technical interviews.

---

## 📑 Table of Contents

- [1. JWT Fundamentals](#1-jwt-fundamentals)
  - [1.1 What is JWT?](#11-what-is-jwt)
  - [1.2 JWT Structure Deep Dive](#12-jwt-structure-deep-dive)
  - [1.3 Standard JWT Claims](#13-standard-jwt-claims)
  - [1.4 JWT vs Session Authentication](#14-jwt-vs-session-authentication)
  - [1.5 How JWT Authentication Works](#15-how-jwt-authentication-works)
- [2. Clean Architecture Layers](#2-clean-architecture-layers)
  - [2.1 Layer Overview](#21-layer-overview)
  - [2.2 Responsibility Matrix](#22-responsibility-matrix)
  - [2.3 Data Flow - Login](#23-data-flow---login)
  - [2.4 Data Flow - Protected Routes](#24-data-flow---protected-routes)
- [3. Password Hashing with bcrypt](#3-password-hashing-with-bcrypt)
  - [3.1 Why Hash Passwords?](#31-why-hash-passwords)
  - [3.2 bcrypt Overview](#32-bcrypt-overview)
  - [3.3 Implementation Guide](#33-implementation-guide)
  - [3.4 Common Mistakes](#34-common-mistakes)
  - [3.5 Best Practices](#35-best-practices)
- [4. Complete Implementation Guide](#4-complete-implementation-guide)
  - [4.1 Domain Layer](#41-domain-layer)
  - [4.2 Service Layer](#42-service-layer)
  - [4.3 Delivery Layer - Auth Infrastructure](#43-delivery-layer---auth-infrastructure)
  - [4.4 Delivery Layer - Middleware](#44-delivery-layer---middleware)
  - [4.5 Delivery Layer - Handlers](#45-delivery-layer---handlers)
  - [4.6 Router Setup](#46-router-setup)
- [5. Security Best Practices](#5-security-best-practices)
  - [5.1 Secret Key Management](#51-secret-key-management)
  - [5.2 Token Expiration](#52-token-expiration)
  - [5.3 HTTPS in Production](#53-https-in-production)
  - [5.4 Token Storage](#54-token-storage)
  - [5.5 What NOT to Store in JWT](#55-what-not-to-store-in-jwt)
  - [5.6 Token Revocation Strategies](#56-token-revocation-strategies)
- [6. Common Pitfalls & Solutions](#6-common-pitfalls--solutions)
- [7. Testing JWT Authentication](#7-testing-jwt-authentication)
  - [7.1 Unit Testing Handlers](#71-unit-testing-handlers)
  - [7.2 Testing Middleware](#72-testing-middleware)
  - [7.3 Integration Tests](#73-integration-tests)
- [8. Interview Questions & Answers](#8-interview-questions--answers)
  - [8.1 Basic Questions](#81-basic-questions)
  - [8.2 Intermediate Questions](#82-intermediate-questions)
  - [8.3 Advanced Questions](#83-advanced-questions)
- [9. Practical Exercises](#9-practical-exercises)
- [10. Quick Reference Cheat Sheet](#10-quick-reference-cheat-sheet)

---

## 1. JWT Fundamentals

### 1.1 What is JWT?

**JWT (JSON Web Token)** is a compact, URL-safe means of representing claims to be transferred between two parties as a JSON object.

**Key Characteristics:**
- ✅ **Self-contained** - Contains all user information needed
- ✅ **Stateless** - Server doesn't store session data
- ✅ **Signed** - Can verify authenticity using cryptographic signature
- ⚠️ **Not encrypted** - Anyone can decode and read the payload (Base64)
- ✅ **Portable** - Works across different domains and services

**When to use JWT:**
- Stateless REST APIs
- Microservices architecture
- Mobile applications
- Single Page Applications (SPAs)
- Cross-domain authentication

---

### 1.2 JWT Structure Deep Dive

A JWT consists of three Base64Url-encoded parts separated by dots (`.`):

HEADER.PAYLOAD.SIGNATURE

markdown

**Example JWT:**
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImVtYWlsIjoidXNlckBleGFtcGxlLmNvbSIsImV4cCI6MTczNTY4OTYwMH0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

yaml

---

#### **Part 1: Header (Red)**

Contains metadata about the token:

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```
alg - Signing algorithm (HS256, RS256, etc.)
typ - Token type (always "JWT")
Base64Url encoded:

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
Part 2: Payload (Purple)
Contains the claims (user data):

```json
{
  "user_id": 123,
  "email": "user@example.com",
  "name": "John Doe",
  "exp": 1735689600,
  "iat": 1735603200
}
```
Base64Url encoded:

eyJ1c2VyX2lkIjoxMjMsImVtYWlsIjoidXNlckBleGFtcGxlLmNvbSIsImV4cCI6MTczNTY4OTYwMH0
Part 3: Signature (Blue)
Ensures the token hasn't been tampered with:

scss
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret
)
```
Result:

SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
🔐 Important: Only the server with the secret key can create valid signatures!

1.3 Standard JWT Claims
Registered Claims (Standardized)
Táblázat
Claim	Full Name	Type	Description	Example
iss	Issuer	String	Who created the token	"auth.example.com"
sub	Subject	String	Subject identifier	"user123"
aud	Audience	String	Intended recipient	"api.example.com"
exp	Expiration Time	Number	Unix timestamp	1735689600
nbf	Not Before	Number	Token valid after this	1735603200
iat	Issued At	Number	Creation time	1735603200
jti	JWT ID	String	Unique identifier	"abc123"
Custom Claims (Application-specific)
```json
{
  "user_id": 123,
  "email": "user@example.com",
  "name": "John Doe",
  "role": "admin",
  "permissions": ["read", "write", "delete"]
}
⚠️ Important: Custom claims should NOT contain sensitive data (passwords, credit cards, etc.)
```
1.4 JWT vs Session Authentication
Táblázat
Feature	JWT	Session
Storage Location	Client (localStorage, cookie)	Server (database, Redis)
Stateless	✅ Yes	❌ No
Scalability	✅ Easy (no shared state)	⚠️ Harder (shared session store)
Server Memory	✅ None	⚠️ Required
Revocation	❌ Difficult	✅ Easy (delete session)
Size	⚠️ Larger (sent every request)	✅ Small (session ID only)
Cross-domain	✅ Easy	⚠️ Requires configuration
Security	⚠️ XSS risk (localStorage)	✅ Better (HttpOnly cookie)
When to use JWT:
Microservices architecture
Mobile applications
Stateless REST APIs
Need horizontal scaling
Cross-domain authentication
When to use Sessions:
Traditional web applications
Need immediate revocation
Smaller data transfer requirements
Single server or shared session store available
1.5 How JWT Authentication Works
sql
┌─────────┐                                    ┌─────────┐
│ Client  │                                    │ Server  │
└────┬────┘                                    └────┬────┘
     │                                              │
     │  1. POST /login                             │
     │     Body: {                                 │
     │       "email": "user@example.com",          │
     │       "password": "Password123"             │
     │     }                                       │
     │─────────────────────────────────────────────>│
     │                                              │
     │                       2. Verify Credentials │
     │                          - Get user from DB │
     │                          - Compare password │
     │                            hash with bcrypt │
     │                                              │
     │                          3. Generate JWT    │
     │                          - Create claims    │
     │                          - Sign with secret │
     │                                              │
     │  4. Response:                               │
     │     {                                       │
     │       "token": "eyJhbG...",                 │
     │       "user": {...}                         │
     │     }                                       │
     │<─────────────────────────────────────────────│
     │                                              │
     │  5. Store Token                             │
     │     localStorage.setItem('token', ...)      │
     │                                              │
     │  6. GET /protected-resource                 │
     │     Headers:                                │
     │       Authorization: Bearer eyJhbG...       │
     │─────────────────────────────────────────────>│
     │                                              │
     │                          7. Validate Token  │
     │                          - Extract from hdr │
     │                          - Verify signature │
     │                          - Check expiration │
     │                          - Extract user info│
     │                                              │
     │  8. Response: { data: [...] }               │
     │<─────────────────────────────────────────────│
     │                                              │
2. Clean Architecture Layers
2.1 Layer Overview
Clean Architecture separates concerns into distinct layers:

yaml
┌───────────────────────────────────────────────────────────────┐
│                      DELIVERY LAYER                           │
│  (HTTP Handlers, Middleware, Request/Response DTOs)           │
│                                                                │
│  Responsibilities:                                             │
│  • Handle HTTP requests/responses                             │
│  • Validate input (JSON parsing, validation tags)             │
│  • Generate JWT tokens                                        │
│  • Extract JWT from headers                                   │
│  • Return appropriate HTTP status codes                       │
│                                                                │
│  Files: handlers/, middleware/, auth/                         │
└───────────────────────┬───────────────────────────────────────┘
                        │
                        ▼
┌───────────────────────────────────────────────────────────────┐
│                       SERVICE LAYER                           │
│  (Business Logic, Use Cases)                                  │
│                                                                │
│  Responsibilities:                                             │
│  • User registration logic                                    │
│  • Password verification (bcrypt)                             │
│  • Business rules validation                                  │
│  • Orchestrate repository calls                              │
│  • NO JWT knowledge (infrastructure concern)                  │
│  • NO HTTP knowledge (delivery concern)                       │
│                                                                │
│  Files: service/                                              │
└───────────────────────┬───────────────────────────────────────┘
                        │
                        ▼
┌───────────────────────────────────────────────────────────────┐
│                       DOMAIN LAYER                            │
│  (Entities, Interfaces, Business Errors)                      │
│                                                                │
│  Responsibilities:                                             │
│  • Define core entities (User, Todo)                          │
│  • Define service interfaces                                  │
│  • Define domain errors                                       │
│  • Business validation rules                                  │
│  • NO external dependencies                                   │
│                                                                │
│  Files: domain/                                               │
└───────────────────────┬───────────────────────────────────────┘
                        │
                        ▼
┌───────────────────────────────────────────────────────────────┐
│                     REPOSITORY LAYER                          │
│  (Database Access, Data Persistence)                          │
│                                                                │
│  Responsibilities:                                             │
│  • CRUD operations                                            │
│  • SQL queries                                                │
│  • Data mapping (DB ↔ Domain)                                │
│  • Transaction management                                     │
│                                                                │
│  Files: repository/                                           │
└───────────────────────────────────────────────────────────────┘
### 2.2 Responsibility Matrix
Táblázat
Layer	JWT	Password	HTTP	Database	Business Logic
Domain	❌ No	❌ No	❌ No	❌ No	✅ Yes (rules)
Service	❌ No	✅ Yes (hash/verify)	❌ No	❌ No	✅ Yes (logic)
Delivery	✅ Yes (generate/validate)	❌ No	✅ Yes	❌ No	❌ No
Repository	❌ No	❌ No	❌ No	✅ Yes	❌ No
Key Principle:

JWT is a delivery mechanism, not business logic!

Password hashing is business logic, not delivery concern!

2.3 Data Flow - Login
```go
// ────────────────────────────────────────────────────────────────
// 1. CLIENT REQUEST
// ────────────────────────────────────────────────────────────────
POST /login HTTP/1.1
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password123"
}
````
// ────────────────────────────────────────────────────────────────
// 2. DELIVERY LAYER (Handler)
// ────────────────────────────────────────────────────────────────
 ```go
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
    // Parse and validate request
    var req LoginDTO
    json.NewDecoder(r.Body).Decode(&req)
    validate.Struct(req) // Email format, required fields

    // Call service (business logic)
    user, err := h.Service.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        // Handle errors (401, 500, etc.)
        return
    }

    // ✅ GENERATE JWT (delivery concern - here!)
    claims := auth.NewUserClaims(user, 24*time.Hour)
    _, tokenString, _ := h.TokenAuth.Encode(claims.ToMap())

    // Return response
    utils.WriteJSON(w, http.StatusOK, LoginResponseDTO{
        Token: tokenString,
        User:  UserResponseDTO{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
        },
    })
}

// ────────────────────────────────────────────────────────────────
// 3. SERVICE LAYER (Business Logic)
// ────────────────────────────────────────────────────────────────
func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
    // Fetch user from database
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, domain.ErrInvalidCredentials
    }

    // ✅ VERIFY PASSWORD (business logic - here!)
    err = bcrypt.CompareHashAndPassword(
        []byte(user.Password), // Hashed password from DB
        []byte(password),      // Plain password from user
    )
    if err != nil {
        return nil, domain.ErrInvalidCredentials
    }

    // Return authenticated user (NO JWT knowledge!)
    return user, nil
}

// ────────────────────────────────────────────────────────────────
// 4. REPOSITORY LAYER (Database)
// ────────────────────────────────────────────────────────────────
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User

    query := "SELECT id, name, email, password FROM users WHERE email = ?"
    err := r.DB.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password, // Hashed password
    )

    if err != nil {
        return nil, err
    }

    return &user, nil
}
2.4 Data Flow - Protected Routes
go
// ────────────────────────────────────────────────────────────────
// 1. CLIENT REQUEST
// ────────────────────────────────────────────────────────────────
GET /todos HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

// ────────────────────────────────────────────────────────────────
// 2. MIDDLEWARE CHAIN
// ────────────────────────────────────────────────────────────────

// 2a. Verifier (jwtauth library) - Extracts and validates token
func Verifier(tokenAuth *JWTAuth) func(http.Handler) http.Handler {
    // • Extract token from Authorization header
    // • Verify signature with secret key
    // • Validate expiration (exp claim)
    // • Store token in request context
    // • Continue even if invalid (stores error in context)
}

// 2b. Authenticator (custom) - Checks validation result
func Authenticator(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, _, err := jwtauth.FromContext(r.Context())

        if err != nil || token == nil {
            http.Error(w, "Unauthorized", 401)
            return // Stop here!
        }

        // Validate claims structure
        claims, err := auth.FromTokenClaims(token.PrivateClaims())
        if err != nil {
            http.Error(w, "Invalid token claims", 401)
            return // Stop here!
        }

        next.ServeHTTP(w, r) // Continue
    })
}

// 2c. UserContext (custom) - Extracts user info
func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, _, _ := jwtauth.FromContext(r.Context())

        // Extract claims
        claims, _ := auth.FromTokenClaims(token.PrivateClaims())

        // Create UserContext
        userCtx := auth.NewUserContext(claims)

        // Add to request context
        ctx := userCtx.AddToContext(r.Context())

        // Pass updated context to handler
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// ────────────────────────────────────────────────────────────────
// 3. HANDLER (uses authenticated user)
// ────────────────────────────────────────────────────────────────
func (h *TodoHandlers) ListTodos(w http.ResponseWriter, r *http.Request) {
    // Get user from context (set by middleware)
    userCtx, ok := auth.FromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", 401)
        return
    }

    // Call service with user ID
    todos, err := h.Service.ListTodos(r.Context(), userCtx.ID)

    utils.WriteJSON(w, http.StatusOK, todos)
}

// ────────────────────────────────────────────────────────────────
// 4. SERVICE LAYER
// ────────────────────────────────────────────────────────────────
func (s *TodoService) ListTodos(ctx context.Context, userID int64) ([]*domain.Todo, error) {
    // Business rule: user can only see their own todos
    return s.Repo.GetByUserID(ctx, userID)
}
Middleware Chain Summary:

sql
Request → [1] Verifier → [2] Authenticator → [3] UserContext → Handler
          ↓             ↓                    ↓                  ↓
          Extract &     Check result        Extract user       Use user info
          Validate      Validate claims     Add to context     from context
          signature     structure
3. Password Hashing with bcrypt
3.1 Why Hash Passwords?
❌ NEVER store passwords in plain text!

go
// ❌ TERRIBLE - Anyone with DB access sees passwords
user := User{
    Email:    "user@example.com",
    Password: "MyPassword123", // Plain text!
}
db.Insert(user)
Risks of plain text passwords:

🔓 Database breach exposes all passwords
🔓 Admins/DBAs can see passwords
🔓 Users reuse passwords across sites
🔓 Password hints/reset becomes dangerous
✅ Always hash passwords before storing:

go
// ✅ GOOD - Password is hashed
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("MyPassword123"), bcrypt.DefaultCost)
user := User{
    Email:    "user@example.com",
    Password: string(hashedPassword), // "$2a$10$N9qo8uLOick..."
}
db.Insert(user)
3.2 bcrypt Overview
bcrypt is a password-hashing function designed to be slow and resistant to brute-force attacks.

Key Features:

✅ Salted - Random salt added to each password (prevents rainbow tables)
✅ Adaptive - Cost factor can increase as computers get faster
✅ Slow by design - Makes brute force attacks impractical
✅ One-way - Cannot decrypt back to original password
How bcrypt works:

java
Password: "MyPassword123"
    ↓
bcrypt.GenerateFromPassword(password, cost)
    ↓
1. Generate random salt (16 bytes)
2. Combine salt + password
3. Apply bcrypt algorithm (2^cost iterations)
    ↓
Hashed: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
         └─┬─┘└┬┘└────────┬────────┘└──────────────┬──────────────┘
          alg cost   salt (22 chars)    hash (31 chars)
Parts of a bcrypt hash:

$2a$ - bcrypt algorithm version
10 - Cost factor (2^10 = 1024 iterations)
N9qo8uLOickgx2ZMRZoMye - Salt (22 characters)
IjZAgcfl7p92ldGxad68LJZdL17lhWy - Hash (31 characters)
3.3 Implementation Guide
Registration - Hash Password
go
package service

import (
    "context"
    "github.com/macesz/todo-go/domain"
    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    Repo domain.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
    // Check if email already exists
    existing, _ := s.Repo.GetByEmail(ctx, email)
    if existing != nil {
        return nil, domain.ErrDuplicate
    }

    // ✅ Hash password with bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(password),      // Plain text password
        bcrypt.DefaultCost,    // Cost factor (10)
    )
    if err != nil {
        return nil, err
    }

    // Create user with hashed password
    user := &domain.User{
        Name:     name,
        Email:    email,
        Password: string(hashedPassword), // Store hashed password
    }

    // Save to database
    if err := s.Repo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}
Cost Factors:

Táblázat
Cost	Iterations	Time (approx)	Use Case
4	16	~0.001s	Testing only
8	256	~0.02s	Very weak
10	1,024	~0.1s	Default (recommended)
12	4,096	~0.4s	High security
14	16,384	~1.6s	Very high security
Recommendation: Use bcrypt.DefaultCost (10) for most applications.

Login - Verify Password
go
func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
    // Fetch user from database
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        // Don't reveal if email exists or not
        return nil, domain.ErrInvalidCredentials
    }

    // ✅ Compare password with bcrypt
    err = bcrypt.CompareHashAndPassword(
        []byte(user.Password), // Hashed password from database
        []byte(password),      // Plain text password from login request
    )
    if err != nil {
        // Password doesn't match
        return nil, domain.ErrInvalidCredentials
    }

    // Password is correct!
    return user, nil
}
How CompareHashAndPassword works:

vbnet
Input:
  - hashedPassword: "$2a$10$N9qo8uLO..."
  - plainPassword: "MyPassword123"

Process:
  1. Extract salt from hashedPassword
  2. Hash plainPassword with extracted salt and same cost
  3. Compare resulting hash with stored hash

Result:
  - If hashes match → return nil (success)
  - If hashes don't match → return error
Important Notes:

✅ CompareHashAndPassword is constant-time (prevents timing attacks)
✅ Always returns error for invalid password (don't check error message)
✅ Same password + same salt = same hash (deterministic)
✅ Different passwords = different hashes
3.4 Common Mistakes
❌ Mistake 1: Comparing hashes directly
go
// ❌ WRONG - Hashes are different every time!
if user.Password == password {
    // This will NEVER work!
}

// ❌ WRONG - Even hashing again doesn't work
hashedInput, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if user.Password == string(hashedInput) {
    // Still wrong! Salt is different each time
}

// ✅ CORRECT - Use CompareHashAndPassword
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
if err != nil {
    // Password is wrong
}
❌ Mistake 2: Hashing in handler layer
go
// ❌ WRONG - Business logic in delivery layer
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserDTO
    json.NewDecoder(r.Body).Decode(&req)

    // Hashing here is wrong layer!
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

    user, _ := h.Service.CreateUser(r.Context(), req.Name, req.Email, string(hashedPassword))
    // ...
}

// ✅ CORRECT - Hashing in service layer
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*User, error) {
    // Hash here - it's business logic!
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    // ...
}
❌ Mistake 3: Revealing if email exists
go
// ❌ WRONG - Tells attacker if email is registered
func (s *UserService) Login(ctx context.Context, email, password string) (*User, error) {
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, errors.New("Email not found") // Reveals email doesn't exist
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, errors.New("Invalid password") // Reveals email exists
    }

    return user, nil
}

// ✅ CORRECT - Generic error message
func (s *UserService) Login(ctx context.Context, email, password string) (*User, error) {
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, domain.ErrInvalidCredentials // Generic message
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, domain.ErrInvalidCredentials // Same generic message
    }

    return user, nil
}
❌ Mistake 4: Not handling bcrypt errors
go
// ❌ WRONG - Ignoring errors
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// What if this fails?

// ✅ CORRECT - Handle errors
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
    return nil, fmt.Errorf("failed to hash password: %w", err)
}
3.5 Best Practices
✅ 1. Use bcrypt.DefaultCost
go
// ✅ GOOD - Use default cost (10)
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    bcrypt.DefaultCost, // Cost = 10
)

// ⚠️ Only increase for high-security scenarios
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    12, // Higher cost = slower but more secure
)
✅ 2. Return generic error messages
go
// ✅ GOOD - Don't reveal details
var ErrInvalidCredentials = errors.New("invalid email or password")

func (s *UserService) Login(ctx context.Context, email, password string) (*User, error) {
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, ErrInvalidCredentials // Generic
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, ErrInvalidCredentials // Same message
    }

    return user, nil
}
✅ 3. Validate password strength before hashing
go
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*User, error) {
    // Validate password strength BEFORE hashing
    if len(password) < 8 {
        return nil, domain.ErrWeakPassword
    }

    // Could also check for:
    // - Uppercase letters
    // - Numbers
    // - Special characters

    // Now hash
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    // ...
}
✅ 4. Test with known passwords
go
func TestLogin_Success(t *testing.T) {
    // Create test user with known password
    plainPassword := "Password123"
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

    testUser := &domain.User{
        ID:       1,
        Email:    "test@example.com",
        Password: string(hashedPassword),
    }

    // Mock repository
    mockRepo := new(mocks.MockUserRepository)
    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(testUser, nil)

    service := NewUserService(mockRepo)

    // Test login with correct password
    user, err := service.Login(context.Background(), "test@example.com", plainPassword)

    assert.NoError(t, err)
    assert.Equal(t, testUser.ID, user.ID)
}
✅ 5. Never log passwords
go
// ❌ WRONG - Logging password
log.Printf("Creating user with email %s and password %s", email, password)

// ❌ WRONG - Logging hashed password
log.Printf("Password hash: %s", hashedPassword)

// ✅ CORRECT - Don't log sensitive data
log.Printf("Creating user with email %s", email)
4. Complete Implementation Guide
4.1 Domain Layer
Purpose: Define core business entities and rules

domain/user.go

go
package domain

import (
    "context"
    "errors"
)

// User - Core business entity
type User struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"` // Hashed password, never in JSON response
}

// UserService - Business logic interface
type UserService interface {
    CreateUser(ctx context.Context, name, email, password string) (*User, error)
    Login(ctx context.Context, email, password string) (*User, error)
    GetUser(ctx context.Context, id int64) (*User, error)
}

// UserRepository - Data access interface
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int64) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
}
domain/errors.go

go
package domain

import "errors"

// Authentication errors
var (
    ErrInvalidCredentials = errors.New("invalid email or password")
    ErrUserNotFound       = errors.New("user not found")
    ErrDuplicate          = errors.New("email already exists")
)

// Validation errors
var (
    ErrInvalidEmail   = errors.New("invalid email address")
    ErrWeakPassword   = errors.New("password must be at least 8 characters")
    ErrInvalidPassword = errors.New("password must contain uppercase, lowercase, and digit")
)
domain/dto.go

go
package domain

// CreateUserDTO - Request for creating user
type CreateUserDTO struct {
    Name     string `json:"name" validate:"required,min=5,max=255"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6,max=255,containsany=0123456789,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}

// LoginDTO - Request for login
type LoginDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

// UserResponseDTO - User information without password
type UserResponseDTO struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// LoginResponseDTO - Response after successful login
type LoginResponseDTO struct {
    Token string          `json:"token"`
    User  UserResponseDTO `json:"user"`
}

// ErrorResponse - Generic error response
type ErrorResponse struct {
    Error string `json:"error"`
}
4.2 Service Layer
Purpose: Implement business logic

service/user_service.go

go
package service

import (
    "context"
    "fmt"

    "github.com/macesz/todo-go/domain"
    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    Repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
    return &UserService{Repo: repo}
}

// CreateUser - Register new user
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
    // Business rule: check if email already exists
    existing, err := s.Repo.GetByEmail(ctx, email)
    if err == nil && existing != nil {
        return nil, domain.ErrDuplicate
    }

    // Business rule: hash password before storing
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(password),
        bcrypt.DefaultCost, // Cost factor: 10
    )
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    // Create user entity
    user := &domain.User{
        Name:     name,
        Email:    email,
        Password: string(hashedPassword),
    }

    // Persist to database
    if err := s.Repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}

// Login - Authenticate user
func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
    // Fetch user from database
    user, err := s.Repo.GetByEmail(ctx, email)
    if err != nil {
        // Don't reveal if email exists
        return nil, domain.ErrInvalidCredentials
    }

    // Business rule: verify password with bcrypt
    err = bcrypt.CompareHashAndPassword(
        []byte(user.Password), // Hashed password from DB
        []byte(password),      // Plain text password from request
    )
    if err != nil {
        // Password doesn't match
        return nil, domain.ErrInvalidCredentials
    }

    // Successfully authenticated
    return user, nil
}

// GetUser - Fetch user by ID
func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
    user, err := s.Repo.GetByID(ctx, id)
    if err != nil {
        return nil, domain.ErrUserNotFound
    }
    return user, nil
}
Key Points:

✅ Password hashing in service layer (business logic)
✅ Generic error messages (don't reveal email existence)
✅ No JWT knowledge
✅ No HTTP knowledge
✅ Returns domain entities
4.3 Delivery Layer - Auth Infrastructure
Purpose: JWT claims and context management

delivery/web/auth/claims.go

go
package auth

import (
    "errors"
    "time"

    "github.com/macesz/todo-go/domain"
)

// UserClaims - JWT payload structure
type UserClaims struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
    Exp    int64  `json:"exp"`
}

// NewUserClaims - Convert domain.User to JWT claims
func NewUserClaims(user *domain.User, expiresIn time.Duration) UserClaims {
    return UserClaims{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
        Exp:    time.Now().Add(expiresIn).Unix(),
    }
}

// ToMap - Convert to map for jwtauth library
func (c UserClaims) ToMap() map[string]interface{} {
    return map[string]interface{}{
        "user_id": c.UserID,
        "email":   c.Email,
        "name":    c.Name,
        "exp":     c.Exp,
    }
}

// FromTokenClaims - Extract and validate claims from JWT
// IMPORTANT: JWT stores numbers as float64, not int64!

func FromTokenClaims(claims map[string]interface{}) (*UserClaims, error) {
    // Extract user_id (JWT stores as float64)
    userID, ok := claims["user_id"].(float64)
    if !ok {
        return nil, errors.New("missing or invalid user_id claim")
    }

    // Extract email
    email, ok := claims["email"].(string)
    if !ok {
        return nil, errors.New("missing or invalid email claim")
    }

    // Extract name
    name, ok := claims["name"].(string)
    if !ok {
        return nil, errors.New("missing or invalid name claim")
    }

    return &UserClaims{
        UserID: int64(userID), // Convert float64 to int64
        Email:  email,
        Name:   name,
    }, nil
}
delivery/web/auth/context.go

go
package auth

import "context"

type contextKey string

const userContextKey contextKey = "userContext"

// UserContext - User information in request context
type UserContext struct {
    ID    int64
    Email string
    Name  string
}

// NewUserContext - Create from JWT claims
func NewUserContext(claims *UserClaims) *UserContext {
    return &UserContext{
        ID:    claims.UserID,
        Email: claims.Email,
        Name:  claims.Name,
    }
}

// AddToContext - Store in request context
func (u *UserContext) AddToContext(ctx context.Context) context.Context {
    return context.WithValue(ctx, userContextKey, u)
}

// FromContext - Retrieve from request context
func FromContext(ctx context.Context) (*UserContext, bool) {
    userCtx, ok := ctx.Value(userContextKey).(*UserContext)
    return userCtx, ok
}
4.4 Delivery Layer - Middleware
Purpose: Extract, validate, and process JWT tokens

delivery/web/middlewares/auth.go

go
package middlewares

import (
    "net/http"

    "github.com/go-chi/jwtauth/v5"
    "github.com/macesz/todo-go/delivery/web/auth"
    "github.com/macesz/todo-go/domain"
)

// Authenticator - Validates token and claims structure
// Runs AFTER jwtauth.Verifier
func Authenticator(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get token from context (set by Verifier)
        token, _, err := jwtauth.FromContext(r.Context())

        // Check if Verifier found issues
        if err != nil {
            writeError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
            return
        }

        if token == nil {
            writeError(w, http.StatusUnauthorized, "Missing authentication token")
            return
        }

        // Validate claims structure
        privateClaims := token.PrivateClaims()
        _, err = auth.FromTokenClaims(privateClaims)
        if err != nil {
            writeError(w, http.StatusUnauthorized, "Invalid token claims: "+err.Error())
            return
        }

        // Token is valid, continue to next handler
        next.ServeHTTP(w, r)
    })
}

// UserContext - Extracts user info and adds to request context
// Runs AFTER Authenticator
func UserContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get token from context
        token, _, err := jwtauth.FromContext(r.Context())
        if err != nil || token == nil {
            writeError(w, http.StatusUnauthorized, "Unauthorized")
            return
        }

        // Extract claims
        privateClaims := token.PrivateClaims()
        claims, err := auth.FromTokenClaims(privateClaims)
        if err != nil {
            writeError(w, http.StatusUnauthorized, "Invalid token claims")
            return
        }

        // Create UserContext and add to request context
        userContext := auth.NewUserContext(claims)
        ctx := userContext.AddToContext(r.Context())

        // Pass updated context to next handler
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Helper function to write JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write([]byte(`{"error":"` + message + `"}`))
}
Middleware Chain Explanation:

vbnet
Request
   ↓
[1] jwtauth.Verifier (library)
   • Extracts "Bearer token" from Authorization header
   • Verifies signature with secret key
   • Checks expiration (exp claim)
   • Stores token + any errors in context
   • ALWAYS continues (even if invalid)
   ↓
[2] Authenticator (custom)
   • Retrieves token from context
   • Checks if Verifier found errors
   • Validates claims structure (user_id, email, name)
   • Returns 401 if invalid
   • Continues if valid
   ↓
[3] UserContext (custom)
   • Extracts claims from token
   • Creates UserContext struct
   • Adds UserContext to request context
   • Handler can now call auth.FromContext()
   ↓
Handler
   • Gets user info from context
   • Uses user.ID for business logic
4.5 Delivery Layer - Handlers
delivery/web/handlers/user.go

go
package handlers

import (
    "encoding/json"
    "errors"
    "net/http"
    "time"

    "github.com/go-chi/jwtauth/v5"
    "github.com/go-playground/validator/v10"
    "github.com/macesz/todo-go/delivery/web/auth"
    "github.com/macesz/todo-go/delivery/web/utils"
    "github.com/macesz/todo-go/domain"
)

type UserHandlers struct {
    Service   domain.UserService
    TokenAuth *jwtauth.JWTAuth
    Validator *validator.Validate
}

func NewUserHandlers(service domain.UserService, tokenAuth *jwtauth.JWTAuth) *UserHandlers {
    return &UserHandlers{
        Service:   service,
        TokenAuth: tokenAuth,
        Validator: validator.New(),
    }
}

// CreateUser - Register new user
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    // Parse request
    var req domain.CreateUserDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{
            Error: "Invalid request body",
        })
        return
    }

    // Validate input
    if err := h.Validator.Struct(req); err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{
            Error: translateValidationError(err),
        })
        return
    }

    // Create user (business logic)
    user, err := h.Service.CreateUser(r.Context(), req.Name, req.Email, req.Password)
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrDuplicate):
            utils.WriteJSON(w, http.StatusConflict, domain.ErrorResponse{
                Error: "Email already exists",
            })
            return
        default:
            utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
                Error: "Internal server error",
            })
            return
        }
    }

    // Return response (without password)
    respUser := domain.UserResponseDTO{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }

    utils.WriteJSON(w, http.StatusCreated, respUser)
}

// Login - Authenticate and return JWT
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    // Parse request
    var req domain.LoginDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{
            Error: "Invalid request body",
        })
        return
    }

    // Validate input
    if err := h.Validator.Struct(req); err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{
            Error: translateValidationError(err),
        })
        return
    }

    // Authenticate user (business logic)
    user, err := h.Service.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        if errors.Is(err, domain.ErrInvalidCredentials) {
            utils.WriteJSON(w, http.StatusUnauthorized, domain.ErrorResponse{
                Error: "Invalid email or password",
            })
            return
        }
        utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
            Error: "Internal server error",
        })
        return
    }

    // Generate JWT (delivery concern)
    claims := auth.NewUserClaims(user, 24*time.Hour)
    _, tokenString, err := h.TokenAuth.Encode(claims.ToMap())
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
            Error: "Failed to generate token",
        })
        return
    }

    // Return response
    utils.WriteJSON(w, http.StatusOK, domain.LoginResponseDTO{
        Token: tokenString,
        User: domain.UserResponseDTO{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
        },
    })
}

// GetUser - Get current authenticated user (protected route)
func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
    // Get authenticated user from context (set by middleware)
    userCtx, ok := auth.FromContext(r.Context())
    if !ok {
        utils.WriteJSON(w, http.StatusUnauthorized, domain.ErrorResponse{
            Error: "Unauthorized",
        })
        return
    }

    // Fetch user details from service
    user, err := h.Service.GetUser(r.Context(), userCtx.ID)
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
            Error: "Failed to fetch user",
        })
        return
    }

    // Return response (without password)
    utils.WriteJSON(w, http.StatusOK, domain.UserResponseDTO{
        ID:    user.ID,
        Email: user.Email,
        Name:  user.Name,
    })
}

// Helper function to translate validation errors
func translateValidationError(err error) string {
    validationErrs, ok := err.(validator.ValidationErrors)
    if !ok {
        return "Validation error"
    }

    var messages []string
    for _, e := range validationErrs {
        field := e.Field()
        tag := e.Tag()

        var message string
        switch tag {
        case "required":
            message = field + " is required"
        case "email":
            message = field + " must be a valid email address"
        case "min":
            message = field + " must be at least " + e.Param() + " characters"
        case "max":
            message = field + " must be at most " + e.Param() + " characters"
        case "containsany":
            if e.Param() == "0123456789" {
                message = field + " must contain at least one digit"
            } else if e.Param() == "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
                message = field + " must contain at least one uppercase letter"
            } else {
                message = field + " must contain specific characters"
            }
        default:
            message = field + " is invalid"
        }

        messages = append(messages, message)
    }

    return strings.Join(messages, "; ")
}
delivery/web/utils/response.go

go
package utils

import (
    "encoding/json"
    "net/http"
)

// WriteJSON - Helper to write JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
4.6 Router Setup
delivery/web/server.go

go
package web

import (
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/jwtauth/v5"
    "github.com/macesz/todo-go/delivery/web/handlers"
    "github.com/macesz/todo-go/delivery/web/middlewares"
)

func SetupRouter(userHandlers *handlers.UserHandlers, todoHandlers *handlers.TodoHandlers) http.Handler {
    r := chi.NewRouter()

    // JWT configuration
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        panic("JWT_SECRET environment variable is required")
    }
    tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

    // Set TokenAuth in handlers
    userHandlers.TokenAuth = tokenAuth

    // Global middlewares
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)

    // ────────────────────────────────────────────────────────────
    // PUBLIC ROUTES (no authentication required)
    // ────────────────────────────────────────────────────────────
    r.Group(func(r chi.Router) {
        r.Post("/register", userHandlers.CreateUser)
        r.Post("/login", userHandlers.Login)
    })

    // ────────────────────────────────────────────────────────────
    // PROTECTED ROUTES (JWT authentication required)
    // ────────────────────────────────────────────────────────────
    r.Group(func(r chi.Router) {
        // JWT Middleware Chain
        r.Use(jwtauth.Verifier(tokenAuth))     // 1. Extract & verify token
        r.Use(middlewares.Authenticator)        // 2. Check validation result
        r.Use(middlewares.UserContext)          // 3. Add user to context

        // User routes
        r.Get("/user", userHandlers.GetUser)

        // Todo routes
        r.Route("/todos", func(r chi.Router) {
            r.Get("/", todoHandlers.ListTodos)
            r.Post("/", todoHandlers.CreateTodo)
            r.Get("/{id}", todoHandlers.GetTodo)
            r.Put("/{id}", todoHandlers.UpdateTodo)
            r.Delete("/{id}", todoHandlers.DeleteTodo)
        })
    })

    return r
}
cmd/server/main.go

go
package main

import (
    "log"
    "net/http"

    "github.com/macesz/todo-go/delivery/web"
    "github.com/macesz/todo-go/delivery/web/handlers"
    "github.com/macesz/todo-go/repository"
    "github.com/macesz/todo-go/service"
)

func main() {
    // Setup database connection
    db, err := setupDatabase()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Create repositories
    userRepo := repository.NewUserRepository(db)
    todoRepo := repository.NewTodoRepository(db)

    // Create services
    userService := service.NewUserService(userRepo)
    todoService := service.NewTodoService(todoRepo)

    // Create handlers
    userHandlers := handlers.NewUserHandlers(userService, nil)
    todoHandlers := handlers.NewTodoHandlers(todoService)

    // Setup router
    router := web.SetupRouter(userHandlers, todoHandlers)

    // Start server
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}
5. Security Best Practices
5.1 Secret Key Management
❌ NEVER hardcode secrets
go
// ❌ TERRIBLE - Secret in code
tokenAuth := jwtauth.New("HS256", []byte("mysecret123"), nil)

// ❌ BAD - Secret in config file committed to Git
config.yaml:
  jwt_secret: "mysecret123"
✅ Use environment variables
go
// ✅ GOOD - Load from environment
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
.env file (NOT committed to Git):

bash
JWT_SECRET=8xK7vZ2qT9nR4wL6pY3mJ5hF8dS1cX0bV7nQ9rP4kM2e
.gitignore:

bash
.env
*.env
Generate Strong Secrets
Minimum Requirements:

✅ At least 256 bits (32 bytes) for HS256
✅ Cryptographically random
✅ Different for each environment
Generate secret:

bash
# Using openssl
openssl rand -base64 32

# Using Go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

func main() {
    b := make([]byte, 32)
    rand.Read(b)
    fmt.Println(base64.StdEncoding.EncodeToString(b))
}

### 5.2 Token Expiration
Short-lived tokens are more secure

// ❌ Too long - security risk if token is stolen
claims := auth.NewUserClaims(user, 30*24*time.Hour) // 30 days

// ✅ Short-lived access token
claims := auth.NewUserClaims(user, 15*time.Minute) // 15 minutes

// ✅ Reasonable for development
claims := auth.NewUserClaims(
```
