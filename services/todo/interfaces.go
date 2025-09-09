package todo

import "github.com/macesz/todo-go/domain"

// TodoStore defines the interface for a todo storage backend.
// Like a Java interface
type TodoStore interface {
	List() []domain.Todo
	Create(title string) (domain.Todo, error)
	Get(id int) (domain.Todo, bool)
	Update(id int, title string, done bool) (domain.Todo, error)
	Delete(id int) bool
}
