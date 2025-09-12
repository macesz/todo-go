package todo

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

// ListTodos returns all todos
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, filtering, sorting, etc.

func (s *TodoService) ListTodos(ctx context.Context) ([]domain.Todo, error) {
	return s.Store.List(ctx) // Delegate to the store
}

// CreateTodo creates a new todo with the given title
// Returns the created Todo or an error
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, checking for duplicates, logging, etc.
func (s *TodoService) CreateTodo(ctx context.Context, title string) (domain.Todo, error) {
	return s.Store.Create(ctx, title) // Delegate to the store
}

// GetTodo retrieves a todo by ID
// Returns the Todo and a boolean indicating if it was found
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, logging, access control, etc.
func (s *TodoService) GetTodo(ctx context.Context, id int) (domain.Todo, error) {
	return s.Store.Get(ctx, id) // Delegate to the store
}

// UpdateTodo updates an existing todo by ID
// Returns the updated Todo or an error
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, validation, logging, etc.
func (s *TodoService) UpdateTodo(ctx context.Context, id int, title string, done bool) (domain.Todo, error) {
	return s.Store.Update(ctx, id, title, done) // Delegate to the store
}

// DeleteTodo deletes a todo by ID
// Returns a boolean indicating if the deletion was successful
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, logging, cascading deletes, etc.
func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	return s.Store.Delete(ctx, id) // Delegate to the store
}
