package pgtodolist

import (
	"time"

	"github.com/macesz/todo-go/domain"
)

type rowDTO struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Title     string    `db:"title"`
	Color     string    `db:"color"`
	Labels    []string  `db:"labels"`
	CreatedAt time.Time `db:"created_at"`
}

func (r rowDTO) ToDomain() *domain.TodoList {
	return &domain.TodoList{
		ID:        r.ID,
		UserID:    r.UserID,
		Title:     r.Title,
		Color:     r.Color,
		Labels:    r.Labels,
		CreatedAt: r.CreatedAt,
	}
}
