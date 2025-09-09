package todo

// TodoService contains business logic for managing todos.
// Like a service class in Java or JS
type TodoService struct {
	Store TodoStore // Dependency injection of the store (like a private field in Java)
}

// Constructor for TodoService (Go doesn't have classes, so we use functions)
// like new TodoService(store) in JS or Java
// Here we inject the store dependency (like constructor injection in Java)
func NewTodoService(store TodoStore) *TodoService {
	return &TodoService{
		Store: store, // Assign the store to the service
	}
}
