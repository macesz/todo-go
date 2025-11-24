package todolist

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/delivery/web/utils"
	"github.com/macesz/todo-go/domain"
)

func (h *TodoListHandlers) List(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	todoLists, err := h.todoListService.List(r.Context(), user.ID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, todoLists)

}

func (h *TodoListHandlers) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	userctx, ok := auth.UserFromContext(ctx)
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	user, err := h.userService.GetUser(ctx, userctx.ID)
	if err != nil || user == nil {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	var reqTodoList domain.CreateTodoListRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&reqTodoList); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}
	todoList, err := h.todoListService.Create(ctx, user.ID, reqTodoList.Title, *reqTodoList.Color, reqTodoList.Labels)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidTitle) {
			utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"})
		return
	}

	respTodoList := domain.TodoListDTO{
		ID:        todoList.ID,
		UserID:    todoList.UserID,
		Title:     todoList.Title,
		Color:     &todoList.Color,
		Labels:    todoList.Labels,
		CreatedAt: todoList.CreatedAt.Format(time.RFC3339),
	}

	utils.WriteJSON(w, http.StatusCreated, respTodoList)

}

func (h *TodoListHandlers) GetListByID(w http.ResponseWriter, r *http.Request) {
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

	todoList, err := h.todoListService.GetListByID(r.Context(), user.ID, id)
	if err != nil {
		if errors.Is(err, domain.ErrListNotFound) { // Check custom error
			utils.WriteJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: err.Error()}) // e.g., {"error": "todo not found"}
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"}) // Generic for security
		return
	}

	itemDTOs := make([]domain.TodoDTO, len(todoList.Items))
	for i, item := range todoList.Items {
		itemDTOs[i] = domain.TodoDTO{
			ID:         item.ID,
			UserID:     item.UserID,
			TodoListID: item.TodoListID,
			Title:      item.Title,
			Done:       item.Done,
			Priority:   item.Priority,
			CreatedAt:  item.CreatedAt.Format(time.RFC3339),
		}
	}

	// Map to response DTO
	respTodoList := domain.TodoListDTO{
		ID:        todoList.ID,
		UserID:    todoList.UserID,
		Title:     todoList.Title,
		Color:     &todoList.Color,
		Labels:    todoList.Labels,
		CreatedAt: todoList.CreatedAt.Format(time.RFC3339),
		Items:     itemDTOs,
	}
	utils.WriteJSON(w, http.StatusOK, respTodoList) // Return the todo as JSON

}

func (h *TodoListHandlers) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	user, ok := auth.UserFromContext(ctx)
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	idr := chi.URLParam(r, "id")
	if idr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	id, err := strconv.ParseInt(idr, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	var todoListDtO domain.UpdateTodoListRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&todoListDtO); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()}) // Using struct for consistency
		return
	}

	updated, err := h.todoListService.Update(ctx, user.ID, id, todoListDtO.Title, *todoListDtO.Color, todoListDtO.Labels)
	if err != nil {
		if errors.Is(err, domain.ErrListNotFound) { // Check custom error )
			utils.WriteJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: err.Error()}) // e.g., {"error": "todo not found"}
			return
		} else if errors.Is(err, domain.ErrInvalidTitle) { // Optional: If service returns this
			utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"}) // Generic for security
		return
	}

	respTodoList := domain.TodoListDTO{
		ID:     updated.ID,
		UserID: user.ID,
		Title:  updated.Title,
		Color:  &updated.Color,
		Labels: updated.Labels,
	}

	utils.WriteJSON(w, http.StatusOK, respTodoList)
}

func (h *TodoListHandlers) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, ok := auth.UserFromContext(ctx)
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, domain.ErrorResponse{Error: "missing user"})
		return
	}

	idr := chi.URLParam(r, "id")
	if idr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	id, err := strconv.ParseInt(idr, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	if err := h.todoListService.Delete(ctx, user.ID, id); err != nil {
		if errors.Is(err, domain.ErrListNotFound) {
			utils.WriteJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"}) // Generic for security
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
