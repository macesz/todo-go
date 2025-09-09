package web

import (
	"log"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartServer(todoService TodoService) {
	// Chi router: like Express app or Java Servlet
	r := chi.NewRouter()

	// Chi middlewares: small, composable functions that wrap handlers.
	r.Use(middleware.RequestID) // Adds a unique request ID in the context
	r.Use(middleware.RealIP)    // Sets RemoteAddr to the real client IP from headers
	r.Use(middleware.Logger)    // Logs the start and end of each request
	r.Use(middleware.Recoverer) // Recovers from panics, returns 500 instead of crashing

	todoHandler := NewHandlers(todoService) // Create handlers with the service

	// Routes
	r.Get("/todos", todoHandler.ListTodos)          // List all todos
	r.Post("/todos", todoHandler.CreateTodo)        // Create a new todo
	r.Get("/todos/{id}", todoHandler.GetTodo)       // Get a todo by ID\
	r.Put("/todos/{id}", todoHandler.UpdateTodo)    // Update a todo by ID
	r.Delete("/todos/{id}", todoHandler.DeleteTodo) // Delete a todo by ID

	// Start the server
	log.Println("listening on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
