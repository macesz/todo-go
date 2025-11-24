package todolist

type TodoListService struct {
	Store TodoListStore
}

func NewTodoListService(store TodoListStore) *TodoListService {
	return &TodoListService{
		Store: store, // Assign the store to the service
	}
}
