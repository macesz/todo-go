package todo

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

type TodoService interface {
	ListTodos(ctx context.Context, userID int64) ([]*domain.Todo, error)
	CreateTodo(ctx context.Context, userID int64, title string, priority int64) (*domain.Todo, error)
	GetTodo(ctx context.Context, userID int64, id int64) (*domain.Todo, error)
	UpdateTodo(ctx context.Context, userID int64, id int64, title string, done bool, priority int64) (*domain.Todo, error)
	DeleteTodo(ctx context.Context, userID int64, id int64) error
}
