package todolist

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

type TodoListService interface {
	List(ctx context.Context, userID int64) ([]*domain.TodoList, error)
	Create(ctx context.Context, userID int64, title string, color string, labels []string) (*domain.TodoList, error)
	Get(ctx context.Context, userID int64, id int64) (*domain.TodoList, error)
	Update(ctx context.Context, userID int64, id int64, title string, color string, labes []string) (*domain.TodoList, error)
	Delete(ctx context.Context, userID int64, id int64) error
}

type UserService interface {
	GetUser(ctx context.Context, id int64) (*domain.User, error)
}
