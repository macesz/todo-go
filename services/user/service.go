package user

import (
	"context"
	"fmt"

	"github.com/macesz/todo-go/domain"
	// "golang.org/x/crypto/bcrypt"
)

// create user
func (u *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, fmt.Errorf("missing required fields: %w", domain.ErrInvalidInput)
	}

	// check user email is already exists
	existingUser, err := u.UserStore.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, fmt.Errorf("email already in use: %w", domain.ErrDuplicate)
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	// Call the UserStore to save the user
	createduser, err := u.UserStore.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in store: %w", err) // Wrap unexpected errors
	}

	return createduser, nil
}

// get user by id
func (u *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return u.UserStore.GetUser(ctx, id)
}

// user login
func (u *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	return u.UserStore.Login(ctx, email, password)
}

// delete user by id
func (u *UserService) DeleteUser(ctx context.Context, id int64) error {
	return u.UserStore.DeleteUser(ctx, id)
}
