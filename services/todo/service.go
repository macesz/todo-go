package todo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/macesz/todo-go/domain"
)

// ListTodos returns all todos
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, filtering, sorting, etc.

func (s *TodoService) ListTodos(ctx context.Context, userID int64, listID int64) ([]*domain.Todo, error) {
	todos, err := s.Store.List(ctx, userID, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}
	return todos, nil
}

// CreateTodo creates a new todo with the given title
// Returns the created Todo or an error
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, checking for duplicates, logging, etc.
func (s *TodoService) CreateTodo(ctx context.Context, userID int64, listID int64, title string, priority int64) (*domain.Todo, error) {
	// Validate title
	if title == "" {
		return nil, domain.ErrInvalidTitle
	}

	// Validate priority
	if priority < 1 || priority > 5 {
		return nil, fmt.Errorf("priority must be between 1 and 5: %w", domain.ErrInvalidInput)
	}
	createdAt := time.Now()

	todo := &domain.Todo{
		UserID:    userID,
		ListID:    listID,
		Title:     title,
		Done:      false,
		Priority:  priority,
		CreatedAt: createdAt,
	}

	err := s.Store.Create(ctx, listID, todo) // Delegate to the store
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil

}

// GetTodo retrieves a todo by ID
// Returns the Todo and a boolean indicating if it was found
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, logging, access control, etc.
func (s *TodoService) GetTodo(ctx context.Context, userID int64, id int64) (*domain.Todo, error) {
	todo, err := s.Store.Get(ctx, id) // Delegate to the store
	if err != nil {
		// Convert sql.ErrNoRows to domain.ErrNotFound
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	if todo.UserID != userID {
		return nil, domain.ErrNotFound
	}

	return todo, nil
}

// UpdateTodo updates an existing todo by ID

func (s *TodoService) UpdateTodo(ctx context.Context, userID int64, id int64, title string, done bool, priority int64) (*domain.Todo, error) {

	if priority < 1 || priority > 5 {
		return nil, fmt.Errorf("priority must be between 1 and 5: %w", domain.ErrInvalidInput)
	}

	_, err := s.GetTodo(ctx, userID, id)
	if err != nil {
		// GetTodo already returns domain.ErrNotFound if not found or not owned
		return nil, err
	}

	updated, err := s.Store.Update(ctx, id, title, done, priority)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return updated, nil
}

// DeleteTodo deletes a todo by ID

func (s *TodoService) DeleteTodo(ctx context.Context, userID int64, id int64) error {
	if _, err := s.GetTodo(ctx, userID, id); err != nil {
		return err
	}

	err := s.Store.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil

}
