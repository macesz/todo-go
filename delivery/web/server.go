package web

import (
	"io"
	"log"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/macesz/todo-go/delivery/web/todo"
	"github.com/macesz/todo-go/delivery/web/user"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

// StartServer initializes the router, sets up routes, and starts the HTTP server.
// It takes a TodoService to handle business logic.
// Like setting up an Express app or Java Servlet.

func StartServer(todoService todo.TodoService, userService user.UserService) {
	// Chi router: like Express app or Java Servlet
	r := chi.NewRouter()

	// Chi middlewares: small, composable functions that wrap handlers.
	r.Use(middleware.RequestID) // Adds a unique request ID in the context
	r.Use(middleware.RealIP)    // Sets RemoteAddr to the real client IP from headers
	r.Use(middleware.Logger)    // Logs the start and end of each request
	r.Use(middleware.Recoverer) // Recovers from panics, returns 500 instead of crashing

	todoHandler := todo.NewHandlers(todoService) // Create handlers with the service
	userHandler := user.NewHandlers(userService) // Create handlers with the service

	// Routes

	// Public Routes
	r.Group(func(r chi.Router) {
		// r.Get("/", indexPage)
		// r.Get("/{AssetUrl}", GetAsset)
		r.Post("/user", userHandler.CreateUser) // Create a new user
		r.Post("/login", userHandler.LoginUser) // Login a user
	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		// r.Use(AuthMiddleware)
		r.Use(middleware.AllowContentType("application/json", "text/xml"))

		r.Route("/todos", func(r chi.Router) {
			r.Get("/", todoHandler.ListTodos)         // List all todos
			r.Get("/{id}", todoHandler.GetTodo)       // Get specific todo by ID
			r.Post("/", todoHandler.CreateTodo)       // Create a new todo
			r.Put("/{id}", todoHandler.UpdateTodo)    // Update a todo by ID
			r.Delete("/{id}", todoHandler.DeleteTodo) // Delete a todo by ID
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/", userHandler.GetUser)           // Get specific user by ID
			r.Delete("/{id}", userHandler.DeleteUser) // Delete a user by ID
		})
	})

	// Start the server
	log.Println("listening on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
