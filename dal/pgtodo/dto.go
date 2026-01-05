package pgtodo

import (
	"time"

	"github.com/macesz/todo-go/domain"
)

type rowDTO struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	TodlistID int64     `db:"todolist_id"`
	Title     string    `db:"title"`
	Done      bool      `db:"done"`
	CreatedAt time.Time `db:"created_at"`
}

func (r rowDTO) ToDomain() *domain.Todo {
	return &domain.Todo{
		ID:         r.ID,
		UserID:     r.UserID,
		TodoListID: r.TodlistID,
		Title:      r.Title,
		Done:       r.Done,
		CreatedAt:  r.CreatedAt,
	}
}
