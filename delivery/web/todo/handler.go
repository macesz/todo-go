package todo

import (
	"encoding/json" // For JSON (like JSON.parse/stringify in JS)
	"errors"
	"fmt"
	"net/http" // Standard HTTP library (like fetch in JS or HttpServlet in Java)
	"strconv"
	"strings"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	validate "github.com/go-playground/validator/v10" // For struct validation (like Joi in JS or Hibernate Validator in Java)
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/delivery/web/utils"
	"github.com/macesz/todo-go/domain"
	// String conversions (like parseInt in JS)
	// String utils (like .split() in JS)
)

// ListTodos handles GET /todos requests.
func (h *TodoHandlers) ListTodos(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	idr := chi.URLParam(r, "listID") // Get the "id" URL parameter

	// Check if id parameter exists
	if idr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	// Convert id string to int64
	listID, err := strconv.ParseInt(idr, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	todos, err := h.todoService.ListTodos(r.Context(), user.ID, listID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"})
		return
	}

	respTodos := make([]domain.TodoDTO, len(todos))
	for _, todo := range todos {
		respTodo := domain.TodoDTO{
			ID:         todo.ID,
			UserID:     todo.UserID,
			TodoListID: todo.TodoListID,
			Title:      todo.Title,
			Done:       todo.Done,
			CreatedAt:  todo.CreatedAt.Format(time.RFC3339),
		}
		respTodos = append(respTodos, respTodo)
	}
	utils.WriteJSON(w, http.StatusOK, respTodos)
}

// CreateTodo handles POST /todos requests.
func (h *TodoHandlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	userCtx, ok := auth.UserFromContext(ctx)
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	user, err := h.userService.GetUser(ctx, userCtx.ID)
	if err != nil || user == nil {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	idr := chi.URLParam(r, "listID") // Get the "id" URL parameter

	// Check if id parameter exists
	if idr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	// Convert id string to int64
	listID, err := strconv.ParseInt(idr, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	var reqTodo domain.CreateTodoDTO // Empty Todo struct to decode into

	// Decode the JSON body into the todo struct
	// json.NewDecoder is like JSON.parse in JS
	// r.Body is the request body (like req.body in Express)
	// &todo is the address of the todo variable (like passing by reference in Java)
	if err := json.NewDecoder(r.Body).Decode(&reqTodo); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	if err := validate.New().Struct(reqTodo); err != nil {
		useErr := translateValidationError(err)
		// Dynamic message, e.g., "Title is required"
		// Similar to Joi validation errors in JS or Bean Validation in Java
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: useErr})
		return
	}

	// Create the todo using the service
	// If creation fails, return 400 Bad Request
	todo, err := h.todoService.CreateTodo(r.Context(), user.ID, listID, reqTodo.Title)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidTitle) {
			utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"})
		return
	}

	respTodo := domain.TodoDTO{
		ID:         todo.ID,
		UserID:     todo.UserID,
		TodoListID: todo.TodoListID,
		Title:      todo.Title,
		Done:       todo.Done,
		CreatedAt:  todo.CreatedAt.Format(time.RFC3339), // Format time as ISO string
	}

	utils.WriteJSON(w, http.StatusCreated, respTodo)
}

// GetTodo handles GET /lists/{listID}/todos/{id} requests.
func (h *TodoHandlers) GetTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	// get listId parameter
	idrl := chi.URLParam(r, "listID") // Get the "id" URL parameter
	if idrl == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	todolistID, err := strconv.ParseInt(idrl, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	idr := chi.URLParam(r, "id") // Get the "id" URL parameter
	if idr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	// Convert id string to int64
	id, err := strconv.ParseInt(idr, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	// Get the todo from the service
	todo, err := h.todoService.GetTodo(r.Context(), user.ID, id)
	if err != nil {

		if errors.Is(err, domain.ErrNotFound) {
			utils.WriteJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: err.Error()}) // e.g., {"error": "todo not found"}
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"}) // Generic for security
		return
	}

	// Map to response DTO
	respTodo := domain.TodoDTO{
		ID:         todo.ID,
		UserID:     todo.UserID,
		TodoListID: todolistID,
		Title:      todo.Title,
		Done:       todo.Done,
		CreatedAt:  todo.CreatedAt.Format(time.RFC3339), // Format time as ISO string
	}

	utils.WriteJSON(w, http.StatusOK, respTodo) // Return the todo as JSON
}

// UpdateTodo handles PUT /lists/{listID}/todos/{id} requests.
func (h *TodoHandlers) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	//get {listID} parameter
	idrl := chi.URLParam(r, "listID") // Get the "id" URL parameter

	if idrl == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	todolistID, err := strconv.ParseInt(idrl, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	idr := chi.URLParam(r, "id") // Get the "id" URL parameter

	if idr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	id, err := strconv.ParseInt(idr, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	var todoDTO domain.UpdateTodoDTO // Empty Todo struct to decode into

	// Decode the JSON body into the todo struct
	// If decoding fails, return 400 Bad Request
	if err := json.NewDecoder(r.Body).Decode(&todoDTO); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()}) // Using struct for consistency
		return
	}

	defer r.Body.Close() // Clean up - like closing a file; prevents leaks

	// Validate using tags in UpdateTodoDTO (like Joi.validate in JS)
	if err := validate.New().Struct(todoDTO); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()}) // Dynamic message, e.g., "Title is required"
		return
	}

	// Call service to update (passes context for timeouts/cancellation)
	updated, err := h.todoService.UpdateTodo(r.Context(), user.ID, id, todoDTO.Title, todoDTO.Done)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) { // Check custom error )
			utils.WriteJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: err.Error()}) // e.g., {"error": "todo not found"}
			return
		} else if errors.Is(err, domain.ErrInvalidTitle) { // Optional: If service returns this
			utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
			return
		}
		// TODO: Add logging here, e.g., log.Printf("Internal error updating todo %d: %v", id, err)
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"}) // Generic for security
		return
	}

	respTodo := domain.TodoDTO{
		ID:         updated.ID,
		UserID:     user.ID,
		TodoListID: todolistID,
		Title:      updated.Title,
		Done:       updated.Done,
	}

	utils.WriteJSON(w, http.StatusOK, respTodo) // Return the updated todo as JSON
}

// DeleteTodo handles DELETE /todos/{id} requests.
func (h *TodoHandlers) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	idr := chi.URLParam(r, "id") // Get the "id" URL parameter
	if idr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	id, err := strconv.ParseInt(idr, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	if err := h.todoService.DeleteTodo(r.Context(), user.ID, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.WriteJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: err.Error()}) // e.g., {"error": "todo not found"}
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"}) // Generic for security
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// translateValidationError converts validator errors to user-friendly strings
func translateValidationError(err error) string {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return "validation failed"
	}

	messages := []string{}
	for _, fieldErr := range validationErrs {
		switch fieldErr.Field() {
		case "Title":
			switch fieldErr.Tag() {
			case "required":
				messages = append(messages, "title is required")
			case "max":
				messages = append(messages, "title must be at most 255 characters")
			default:
				messages = append(messages, "title is invalid")
			}
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", strings.ToLower(fieldErr.Field())))
		}

	}

	if len(messages) == 0 {
		return "validation failed"
	}

	return strings.Join(messages, "; ") // Combine if multiple errors
}
