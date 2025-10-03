package user

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

type UserService interface {
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	CreateUser(ctx context.Context, name, email, password string) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
}
