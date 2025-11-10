package todo

// TodoHandlers groups HTTP handler functions.
// Like a Java controller class or JS route handler object.
type TodoHandlers struct {
	todoService TodoService
	userService UserService
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(todoService TodoService, userService UserService) *TodoHandlers {
	return &TodoHandlers{
		todoService: todoService,
		userService: userService,
	}
}
