package domain

import "time"

type TodoList struct {
	ID     int64
	UserID int64

	Title     string
	Color     string
	Labels    []string
	CreatedAt time.Time
}
