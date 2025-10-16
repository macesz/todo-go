package user

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

// create user
func (u *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, domain.ErrInvalidInput
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	return u.UserStore.CreateUser(ctx, user)
}

// get user by id
func (u *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return u.UserStore.GetUser(ctx, id)
}

// delete user by id
func (u *UserService) DeleteUser(ctx context.Context, id int64) error {
	return u.UserStore.DeleteUser(ctx, id)
}
