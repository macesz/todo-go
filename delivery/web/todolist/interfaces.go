package todolist

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

type TodoListService interface {
	ListTodos(ctx context.Context, userID int64) ([]*domain.TodoList, error)
	CreateTodo(ctx context.Context, userID int64, title string, color string, labels []string) (*domain.TodoList, error)
	GetTodo(ctx context.Context, userID int64, id int64) (*domain.TodoList, error)
	UpdateTodo(ctx context.Context, userID int64, id int64, title string, color string, labes []string) (*domain.TodoList, error)
	DeleteTodo(ctx context.Context, userID int64, id int64) error
}

type UserService interface {
	GetUser(ctx context.Context, id int64) (*domain.User, error)
}
