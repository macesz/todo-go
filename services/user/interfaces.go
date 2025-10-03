package user

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
}
