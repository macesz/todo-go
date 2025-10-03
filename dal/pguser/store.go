package pguser

import (
	"context"
	"errors"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/pkg"
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

func (s *Store) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[createUserQuery], templateParams)
	if err != nil {
		return nil, err
	}

	queryParams := map[string]any{
		"name":     name,
		"email":    email,
		"password": password,
	}

	result, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return nil, err
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

	user := &domain.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}

	return user, nil
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
		return nil, errors.New("user not found")
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
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
