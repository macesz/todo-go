package pguser

import "github.com/macesz/todo-go/domain"

type rowDTO struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

func (r rowDTO) ToDomain() *domain.User {
	return &domain.User{
		ID:       r.ID,
		Email:    r.Email,
		Name:     r.Name,
		Password: r.Password,
	}
}
