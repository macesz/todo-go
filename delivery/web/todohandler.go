package web

import (
	"encoding/json" // For JSON (like JSON.parse/stringify in JS)
	"io"
	"net/http" // Standard HTTP library (like fetch in JS or HttpServlet in Java)
	"strconv"
	"time"

	chi "github.com/go-chi/chi/v5"
	validate "github.com/go-playground/validator/v10" // For struct validation (like Joi in JS or Hibernate Validator in Java)
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

// NewHealthHandler returns a handler that reports basic health info.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

// ListTodos handles GET /todos requests.
func (h *TodoHandlers) ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.Service.ListTodos(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, todos)
}

// CreateTodo handles POST /todos requests.
func (h *TodoHandlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var reqTodo domain.CreateTodoDTO // Empty Todo struct to decode into

	// Decode the JSON body into the todo struct
	// json.NewDecoder is like JSON.parse in JS
	// r.Body is the request body (like req.body in Express)
	// &todo is the address of the todo variable (like passing by reference in Java)
	// If decoding fails, return 400 Bad Request
	if err := json.NewDecoder(r.Body).Decode(&reqTodo); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if validate.New().Struct(reqTodo) != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required and must be between 1 and 255 characters"})
		return
	}

	// Create the todo using the service
	// If creation fails, return 400 Bad Request
	todo, err := h.Service.CreateTodo(r.Context(), reqTodo.Title)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	respTodo := domain.TodoDTO{
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

	todo, err := h.Service.GetTodo(r.Context(), id) // Get the todo from the service
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	respTodo := domain.TodoDTO{
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

	var todoDTO domain.UpdateTodoDTO // Empty Todo struct to decode into

	// Decode the JSON body into the todo struct
	// If decoding fails, return 400 Bad Request
	if err := json.NewDecoder(r.Body).Decode(&todoDTO); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the UpdateTodoDTO struct
	if validate.New().Struct(todoDTO) != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required and must be between 1 and 255 characters; done is required"})
		return

	}

	// Update the todo using the service
	// If update fails, return 400 Bad Request
	updated, err := h.Service.UpdateTodo(r.Context(), id, todoDTO.Title, todoDTO.Done)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	respTodo := domain.TodoDTO{
		ID:        updated.ID,
		Title:     updated.Title,
		Done:      updated.Done,
		CreatedAt: updated.CreatedAt.Format(time.RFC3339), // Format time as ISO string
	}

	writeJSON(w, http.StatusOK, respTodo) // Return the updated todo as JSON
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

	if err := h.Service.DeleteTodo(r.Context(), id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// writeJSON is a helper to write JSON responses.
// type any = interface{} any is an alias for interface{} and is equivalent to interface{} in all ways.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json") // Set content type header

	w.WriteHeader(status)           // Set the status code
	json.NewEncoder(w).Encode(data) // Encode and write the JSON response
}
