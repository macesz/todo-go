package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/middlewares"
	"github.com/macesz/todo-go/domain"
)

// StartServer initializes the router, sets up routes, and starts the HTTP server.
// Accepts services parameter (dependency injection)
func StartServer(ctx context.Context, conf domain.Config, services *ServerServices, handlers *Handlers) {
	// Chi router: like Express app or Java Servlet
	r := chi.NewRouter()

	// Chi middlewares: small, composable functions that wrap handlers.
	r.Use(middleware.RequestID) // Adds a unique request ID in the context
	r.Use(middleware.RealIP)    // Sets RemoteAddr to the real client IP from headers
	r.Use(middleware.Logger)    // Logs the start and end of each request
	r.Use(middleware.Recoverer) // Recovers from panics, returns 500 instead of crashing

	// ============================================
	// PUBLIC ROUTES (No authentication required)
	// ============================================
	// r.Group(func(r chi.Router) {
	// r.Get("/", indexPage)
	// r.Get("/{AssetUrl}", GetAsset)
	r.Post("/user", handlers.User.CreateUser) // Create a new user
	r.Post("/login", handlers.User.Login)     // Login a user
	// })

	// ============================================
	// PROTECTED ROUTES (JWT authentication required)
	// ============================================
	r.Group(func(r chi.Router) {
		// r.Use(AuthMiddleware)

		// Seek, verify and validate JWT tokens
		// Using the injected TokenAuth from services
		r.Use(jwtauth.Verifier(services.TokenAuth))
		r.Use(middlewares.Authenticator)
		r.Use(middlewares.UserContext)

		r.Use(middleware.AllowContentType("application/json", "text/xml"))

		r.Route("/lists", func(r chi.Router) {
			r.Get("/", handlers.TodoList.List)
			r.Get("/{id}", handlers.TodoList.Get)
			r.Post("/", handlers.TodoList.Create)
			r.Put("/{id}", handlers.TodoList.Update)
			r.Delete("/{id}", handlers.TodoList.Delete)
		})

		r.Route("/todos", func(r chi.Router) {
			r.Get("/", handlers.Todo.ListTodos)         // List all todos
			r.Get("/{id}", handlers.Todo.GetTodo)       // Get specific todo by ID
			r.Post("/", handlers.Todo.CreateTodo)       // Create a new todo
			r.Put("/{id}", handlers.Todo.UpdateTodo)    // Update a todo by ID
			r.Delete("/{id}", handlers.Todo.DeleteTodo) // Delete a todo by ID
		})

		// changed to /users from /user to follow REST conventions, as we need separation for private and protected routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", handlers.User.GetUser)
			r.Delete("/{id}", handlers.User.DeleteUser) // Delete a user by ID
		})
	})

	// Start the server
	log.Printf("listening on :%s", conf.ServerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", conf.ServerPort), r); err != nil {
		log.Fatal(err)
	}
}
