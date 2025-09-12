package todo

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

// TodoStore defines the interface for a todo storage backend.
// Like a Java interface
type TodoStore interface {
	List(ctx context.Context) ([]domain.Todo, error)
	Create(ctx context.Context, title string) (domain.Todo, error)
	Get(ctx context.Context, id int) (domain.Todo, error)
	Update(ctx context.Context, id int, title string, done bool) (domain.Todo, error)
	Delete(ctx context.Context, id int) error
}
