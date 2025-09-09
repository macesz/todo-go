package main

type TodoStore interface {
	List() []Todo
	Create(title string) (Todo, error)
	Get(id int) (Todo, bool)
	Update(id int, title string, done bool) (Todo, error)
	Delete(id int) bool
}

type TodoService struct {
	Store TodoStore
}

func NewTodoService(store TodoStore) *TodoService {
	return &TodoService{
		Store: store,
	}
}

func (s *TodoService) ListTodos() []Todo {
	return s.Store.List()
}

func (s *TodoService) CreateTodo(title string) (Todo, error) {
	return s.Store.Create(title)
}

func (s *TodoService) GetTodo(id int) (Todo, bool) {
	return s.Store.Get(id)
}

func (s *TodoService) UpdateTodo(id int, title string, done bool) (Todo, error) {
	return s.Store.Update(id, title, done)
}

func (s *TodoService) DeleteTodo(id int) bool {
	return s.Store.Delete(id)
}
