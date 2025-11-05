package pgtodo

import (
	"time"

	"github.com/macesz/todo-go/domain"
)

type rowDTO struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"userId"`
	Title     string    `db:"title"`
	Done      bool      `db:"done"`
	Priority  int64     `db:"priority"`
	CreatedAt time.Time `db:"created_at"`
}

func (r rowDTO) ToDomain() *domain.Todo {
	return &domain.Todo{
		ID:        r.ID,
		UserID:    r.UserID,
		Title:     r.Title,
		Done:      r.Done,
		Priority:  r.Priority,
		CreatedAt: r.CreatedAt,
	}
}
