package web

import (
	"context"
	"log"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/middlewares"
	"github.com/macesz/todo-go/domain"
)

// StartServer initializes the router, sets up routes, and starts the HTTP server.
// It takes a TodoService to handle business logic.
// Like setting up an Express app or Java Servlet.
func StartServer(ctx context.Context, conf domain.Config, handlers *Handlers) {
	// Chi router: like Express app or Java Servlet
	r := chi.NewRouter()

	// JWT Auth setup with HS256 and secret from config
	tokenAuth := jwtauth.New("HS256", []byte(conf.JWTSecret), nil)

	// Chi middlewares: small, composable functions that wrap handlers.
	r.Use(middleware.RequestID) // Adds a unique request ID in the context
	r.Use(middleware.RealIP)    // Sets RemoteAddr to the real client IP from headers
	r.Use(middleware.Logger)    // Logs the start and end of each request
	r.Use(middleware.Recoverer) // Recovers from panics, returns 500 instead of crashing

	// ============================================
	// PUBLIC ROUTES (No authentication required)
	// ============================================
	r.Group(func(r chi.Router) {
		// r.Get("/", indexPage)
		// r.Get("/{AssetUrl}", GetAsset)
		r.Post("/user", handlers.User.CreateUser) // Create a new user
		r.Post("/login", handlers.User.Login)     // Login a user
	})

	// ============================================
	// PROTECTED ROUTES (JWT authentication required)
	// ============================================
	r.Group(func(r chi.Router) {
		// r.Use(AuthMiddleware)

		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middlewares.Authenticator)
		r.Use(middlewares.UserContext)

		r.Use(middleware.AllowContentType("application/json", "text/xml"))

		r.Route("/todos", func(r chi.Router) {
			r.Get("/", handlers.Todo.ListTodos)         // List all todos
			r.Get("/{id}", handlers.Todo.GetTodo)       // Get specific todo by ID
			r.Post("/", handlers.Todo.CreateTodo)       // Create a new todo
			r.Put("/{id}", handlers.Todo.UpdateTodo)    // Update a todo by ID
			r.Delete("/{id}", handlers.Todo.DeleteTodo) // Delete a todo by ID
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/", handlers.User.GetUser)
			r.Delete("/{id}", handlers.User.DeleteUser) // Delete a user by ID
		})
	})

	// Start the server
	log.Println("listening on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
