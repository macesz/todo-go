package pgtodo

import (
	"time"

	"github.com/macesz/todo-go/domain"
)

type RowDTO struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	Done      bool      `db:"done"`
	CreatedAt time.Time `db:"created_at"`
}

func (r RowDTO) ToDomain() *domain.Todo {
	return &domain.Todo{
		ID:        r.ID,
		Title:     r.Title,
		Done:      r.Done,
		CreatedAt: r.CreatedAt,
	}
}
