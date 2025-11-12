package todolist

type TodoListService struct {
	Store TodoListStore
}

func NewTodoService(store TodoListStore) *TodoListService {
	return &TodoListService{
		Store: store, // Assign the store to the service
	}
}
