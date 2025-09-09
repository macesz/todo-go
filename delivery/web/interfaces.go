package web

import "github.com/macesz/todo-go/domain"

type TodoService interface {
	ListTodos() []domain.Todo
	CreateTodo(title string) (domain.Todo, error)
	GetTodo(id int) (domain.Todo, bool)
	UpdateTodo(id int, title string, done bool) (domain.Todo, error)
	DeleteTodo(id int) bool
}
