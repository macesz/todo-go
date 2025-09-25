package pgtodo

import "time"

type RowDTO struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Done      bool      `db:"done"`
	CreatedAt time.Time `db:"created_at"`
}
