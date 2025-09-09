package main

import (
	"encoding/json" // For JSON (like JSON.parse/stringify in JS)
	"net/http"      // Standard HTTP library (like fetch in JS or HttpServlet in Java)
	"strconv"

	chi "github.com/go-chi/chi/v5"
	// String conversions (like parseInt in JS)
	// String utils (like .split() in JS)
)

// TodoHandlers groups HTTP handler functions.
// Like a Java controller class or JS route handler object.
type TodoHandlers struct {
	Service *TodoService
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(service *TodoService) *TodoHandlers {
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
	var todo Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.Service.CreateTodo(todo.Title)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, todo)
}

// GetTodo handles GET /todos/{id} requests.
func (h *TodoHandlers) GetTodo(w http.ResponseWriter, r *http.Request) {
	idr := chi.URLParam(r, "id")

	if idr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be an integer"})
		return
	}

	todo, ok := h.Service.GetTodo(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

// UpdateTodo handles PUT /todos/{id} requests.
func (h *TodoHandlers) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idr := chi.URLParam(r, "id")

	if idr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be an integer"})
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.Service.UpdateTodo(id, todo.Title, todo.Done)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// DeleteTodo handles DELETE /todos/{id} requests.
func (h *TodoHandlers) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idr := chi.URLParam(r, "id")

	if idr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idr)
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
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
