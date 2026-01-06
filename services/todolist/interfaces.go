package todolist

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

type TodoListStore interface {
	List(ctx context.Context, userId int64) ([]*domain.TodoList, error)
	GetListByID(ctx context.Context, id int64) (*domain.TodoList, error)
	Create(ctx context.Context, todoList *domain.TodoList) error
	Update(ctx context.Context, id int64, title string, color string, labels []string, deleted bool) (*domain.TodoList, error)
	Delete(ctx context.Context, id int64) error
}
