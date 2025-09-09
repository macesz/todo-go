package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	store := NewInMemoryStore()      // Our in-memory store (like a database)
	service := NewTodoService(store) // Service with business logic
	handlers := NewHandlers(service) // Handlers with HTTP logic

	// Chi router: like Express app or Java Servlet
	r := chi.NewRouter()

	// Chi middlewares: small, composable functions that wrap handlers.
	r.Use(middleware.RequestID) // Adds a unique request ID in the context
	r.Use(middleware.RealIP)    // Sets RemoteAddr to the real client IP from headers
	r.Use(middleware.Logger)    // Logs the start and end of each request
	r.Use(middleware.Recoverer) // Recovers from panics, returns 500 instead of crashing

	// Routes
	r.Get("/todos", handlers.ListTodos)          // List all todos
	r.Post("/todos", handlers.CreateTodo)        // Create a new todo
	r.Get("/todos/{id}", handlers.GetTodo)       // Get a todo by ID\
	r.Put("/todos/{id}", handlers.UpdateTodo)    // Update a todo by ID
	r.Delete("/todos/{id}", handlers.DeleteTodo) // Delete a todo by ID

	// Start the server
	log.Println("listening on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}

}
