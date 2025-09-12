package web

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

type TodoService interface {
	ListTodos(ctx context.Context) ([]domain.Todo, error)
	CreateTodo(ctx context.Context, title string) (domain.Todo, error)
	GetTodo(ctx context.Context, id int) (domain.Todo, error)
	UpdateTodo(ctx context.Context, id int, title string, done bool) (domain.Todo, error)
	DeleteTodo(ctx context.Context, id int) error
}
