package todo

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

// ListTodos returns all todos
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, filtering, sorting, etc.

func (s *TodoService) ListTodos(ctx context.Context, userID int64) ([]*domain.Todo, error) {
	return s.Store.List(ctx, userID) // Delegate to the store
}

// CreateTodo creates a new todo with the given title
// Returns the created Todo or an error
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, checking for duplicates, logging, etc.
func (s *TodoService) CreateTodo(ctx context.Context, userID int64, title string, priority int64) (*domain.Todo, error) {
	return s.Store.Create(ctx, userID, title, priority) // Delegate to the store
}

// GetTodo retrieves a todo by ID
// Returns the Todo and a boolean indicating if it was found
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, logging, access control, etc.
func (s *TodoService) GetTodo(ctx context.Context, userID int64, id int64) (*domain.Todo, error) {
	todo, err := s.Store.Get(ctx, id) // Delegate to the store
	if err != nil {
		return nil, err
	}

	if todo.UserID != userID {
		return nil, domain.ErrTodoNotFound
	}

	return todo, nil
}

// UpdateTodo updates an existing todo by ID

func (s *TodoService) UpdateTodo(ctx context.Context, userID int64, id int64, title string, done bool, priority int64) (*domain.Todo, error) {
	if _, err := s.GetTodo(ctx, userID, id); err != nil {
		return nil, err
	}

	return s.Store.Update(ctx, id, title, done, priority) // Delegate to the store
}

// DeleteTodo deletes a todo by ID

func (s *TodoService) DeleteTodo(ctx context.Context, userID int64, id int64) error {
	if _, err := s.GetTodo(ctx, userID, id); err != nil {
		return err
	}

	return s.Store.Delete(ctx, id) // Delegate to the store
}
