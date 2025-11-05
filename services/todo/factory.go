package todo

// TodoService contains business logic for managing todos.
// Like a service class in Java or JS
type TodoService struct {
	Store TodoStore // Dependency injection of the store (like a private field in Java)
}

// Factory function - Go's equivalent to a constructor in Java
// Java: new TodoService(store)
// Go:   NewTodoService(store)
// Since Go has no classes/constructors, we use factory functions to create and initialize structs
// The "factory" name emphasizes that we're manufacturing instances rather than just initializing them.

// Here we inject the store dependency (like constructor injection in Java)
func NewTodoService(store TodoStore) *TodoService {
	return &TodoService{
		Store: store, // Assign the store to the service
	}
}
