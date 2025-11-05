package user

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
}
