package todolist

type TodoListHandlers struct {
	todoListService TodoListService
	userService     UserService
}

func NewHandlers(todoListService TodoListService, userService UserService) *TodoListHandlers {
	return &TodoListHandlers{
		todoListService: todoListService,
		userService:     userService,
	}
}
