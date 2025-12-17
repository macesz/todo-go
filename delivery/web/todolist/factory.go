package todolist

type TodoListHandlers struct {
	todoListService TodoListService
	todoService     TodoService
	userService     UserService
}

func NewHandlers(todoListService TodoListService, todoService TodoService, userService UserService) *TodoListHandlers {
	return &TodoListHandlers{
		todoListService: todoListService,
		todoService:     todoService,
		userService:     userService,
	}
}
