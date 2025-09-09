package todo

import "github.com/macesz/todo-go/domain"


// ListTodos returns all todos
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, filtering, sorting, etc.

func (s *TodoService) ListTodos() []domain.Todo {
	return s.Store.List() // Delegate to the store
}

// CreateTodo creates a new todo with the given title
// Returns the created Todo or an error
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, checking for duplicates, logging, etc.
func (s *TodoService) CreateTodo(title string) (domain.Todo, error) {
	return s.Store.Create(title) // Delegate to the store
}

// GetTodo retrieves a todo by ID
// Returns the Todo and a boolean indicating if it was found
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, logging, access control, etc.
func (s *TodoService) GetTodo(id int) (domain.Todo, bool) {
	return s.Store.Get(id) // Delegate to the store
}

// UpdateTodo updates an existing todo by ID
// Returns the updated Todo or an error
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, validation, logging, etc.
func (s *TodoService) UpdateTodo(id int, title string, done bool) (domain.Todo, error) {
	return s.Store.Update(id, title, done) // Delegate to the store
}

// DeleteTodo deletes a todo by ID
// Returns a boolean indicating if the deletion was successful
// Like a service method in Java or JS
// Here we could add more business logic if needed
// For example, logging, cascading deletes, etc.
func (s *TodoService) DeleteTodo(id int) bool {
	return s.Store.Delete(id) // Delegate to the store
}
