package web

import (
	"encoding/json" // For JSON (like JSON.parse/stringify in JS)
	"net/http"      // Standard HTTP library (like fetch in JS or HttpServlet in Java)
	"strconv"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/macesz/todo-go/domain"
	// String conversions (like parseInt in JS)
	// String utils (like .split() in JS)
)

// TodoHandlers groups HTTP handler functions.
// Like a Java controller class or JS route handler object.
type TodoHandlers struct {
	Service TodoService
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(service TodoService) *TodoHandlers {
	return &TodoHandlers{
		Service: service,
	}
}

// ListTodos handles GET /todos requests.
func (h *TodoHandlers) ListTodos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Service.ListTodos())
}

// CreateTodo handles POST /todos requests.
func (h *TodoHandlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo domain.Todo // Empty Todo struct to decode into

	// Decode the JSON body into the todo struct
	// json.NewDecoder is like JSON.parse in JS
	// r.Body is the request body (like req.body in Express)
	// &todo is the address of the todo variable (like passing by reference in Java)
	// If decoding fails, return 400 Bad Request
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create the todo using the service
	// If creation fails, return 400 Bad Request
	todo, err := h.Service.CreateTodo(todo.Title)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	respTodo := Todo{
		ID:        todo.ID,
		Title:     todo.Title,
		Done:      todo.Done,
		CreatedAt: todo.CreatedAt.Format(time.RFC3339), // Format time as ISO string
	}

	writeJSON(w, http.StatusCreated, respTodo)
}

// GetTodo handles GET /todos/{id} requests.
func (h *TodoHandlers) GetTodo(w http.ResponseWriter, r *http.Request) {
	idr := chi.URLParam(r, "id") // Get the "id" URL parameter

	if idr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idr) // Convert id string to int
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be an integer"})
		return
	}

	todo, ok := h.Service.GetTodo(id) // Get the todo from the service
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	respTodo := Todo{
		ID:        todo.ID,
		Title:     todo.Title,
		Done:      todo.Done,
		CreatedAt: todo.CreatedAt.Format(time.RFC3339), // Format time as ISO string
	}

	writeJSON(w, http.StatusOK, respTodo) // Return the todo as JSON
}

// UpdateTodo handles PUT /todos/{id} requests.
func (h *TodoHandlers) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idr := chi.URLParam(r, "id") // Get the "id" URL parameter

	if idr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idr) // Convert id string to int
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be an integer"})
		return
	}

	var todo Todo // Empty Todo struct to decode into

	// Decode the JSON body into the todo struct
	// If decoding fails, return 400 Bad Request
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update the todo using the service
	// If update fails, return 400 Bad Request
	updated, err := h.Service.UpdateTodo(id, todo.Title, todo.Done)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updated) // Return the updated todo as JSON
}

// DeleteTodo handles DELETE /todos/{id} requests.
func (h *TodoHandlers) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idr := chi.URLParam(r, "id") // Get the "id" URL parameter

	if idr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idr) // Convert id string to int
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be an integer"})
		return
	}

	if ok := h.Service.DeleteTodo(id); !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// writeJSON is a helper to write JSON responses.
// type any = interface{} any is an alias for interface{} and is equivalent to interface{} in all ways.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json") // Set content type header
	// Set the status code and encode the data as JSON
	// json.NewEncoder is like JSON.stringify in JS
	// w is the response writer (like res in Express)
	// data is the data to encode (can be struct, map, slice, etc.)
	// If encoding fails, we can't do much here, so we ignore the error
	// In a real app, consider logging the error
	w.WriteHeader(status)           // Set the status code
	json.NewEncoder(w).Encode(data) // Encode and write the JSON response
}
