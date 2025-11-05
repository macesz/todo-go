package pguser

import (
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/pkg"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	queryTemplates map[string]*template.Template

	db *sqlx.DB
}

func CreateStore(db *sqlx.DB) *Store {
	queryTemplates, err := pkg.BuildQueries(files, "queries")
	if err != nil {
		panic(err)
	}
	return &Store{
		queryTemplates: queryTemplates,
		db:             db,
	}
}

func (s *Store) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[createUserQuery], templateParams)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	queryParams := map[string]any{
		"name":     user.Name,
		"email":    user.Email,
		"password": string(hashedPassword),
	}

	result, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // "23505" = unique_violation
			return nil, domain.ErrDuplicate
		}
		return nil, fmt.Errorf("db create user : %w", err)
	}

	defer result.Close()

	var id int64

	if result.Next() {
		err = result.Scan(&id)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("failed to retrieve inserted user ID")
	}

	createdUser := &domain.User{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}
	return createdUser, nil
}

func (s *Store) GetUser(ctx context.Context, id int64) (*domain.User, error) {

	querystr, err := pkg.PrepareQuery(s.queryTemplates[getUserQuery], nil)
	if err != nil {
		return nil, err
	}

	queryParams := map[string]any{
		"id": id,
	}

	result, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	var row rowDTO

	if result.Next() {
		err = result.StructScan(&row)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, domain.ErrUserNotFound
	}

	return row.ToDomain(), nil
}

// get user by email for duplicate check
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	querystr, err := pkg.PrepareQuery(s.queryTemplates[getUserByEmailQuery], nil)
	if err != nil {
		return nil, err
	}

	queryParams := map[string]any{
		"email": email,
	}
	result, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var row rowDTO

	if result.Next() {
		err = result.StructScan(&row)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil // No user found with this email
	}

	return row.ToDomain(), nil

}

// Login user
func (s *Store) Login(ctx context.Context, email, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[loginUserQuery], nil)
	if err != nil {
		return nil, err
	}

	queryParams := map[string]any{
		"email":    email,
		"password": string(hashedPassword),
	}

	result, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	var row rowDTO

	if result.Next() {
		err = result.StructScan(&row)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, domain.ErrUserNotFound
	}

	return row.ToDomain(), nil
}

// deleteUserQuery
func (s *Store) DeleteUser(ctx context.Context, id int64) error {

	querystr, err := pkg.PrepareQuery(s.queryTemplates[deleteUserQuery], nil)
	if err != nil {
		return err
	}

	queryParams := map[string]any{
		"id": id,
	}

	result, err := s.db.NamedExecContext(ctx, querystr, queryParams)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db delete user: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
