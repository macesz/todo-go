package pgtodolist

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/pkg"
)

// Here is the Store struct where we store the queries and the database connection.
type Store struct {
	queryTemplates map[string]*template.Template
	db             *sqlx.DB
}

// CreateStore creates a new Store instance.
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

func (s *Store) List(ctx context.Context, userID int64) ([]*domain.TodoList, error) {
	todoLists := make([]*domain.TodoList, 0)

	// Template parameters are not safe to use directly in the query, because they can be used to inject SQL code.
	// I can use anything that is not a user input, like Table Name, Column Name, etc.
	templateParams := map[string]any{}

	// Prepare the query string, by using the template.
	querystr, err := pkg.PrepareQuery(s.queryTemplates[listTodoListQuery], templateParams)
	if err != nil {
		return nil, err
	}

	// Prepare the query parameters.
	// This is safe to use directly in the query, because it uses named parameters.
	queryParams := map[string]any{
		"user_id": userID,
	}

	// Execute the query. You can add parameters to the query if needed instead of using nil.
	//NamedQueryContext ✅ - Multiple rows (ListTodos, Search, etc.)
	rows, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return nil, err
	}

	defer rows.Close() // Important: Always close rows!

	var row rowDTO

	for rows.Next() {
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		todoLists = append(todoLists, row.ToDomain())
	}

	return todoLists, nil
}

func (s *Store) GetListByID(ctx context.Context, id int64) (*domain.TodoList, error) {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[getTodoListQuery], templateParams)
	if err != nil {
		return nil, err
	}

	queryParams := map[string]any{
		"id": id,
	}

	var row rowDTO
	rows, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
	} else {
		// Return sql.ErrNoRows so the service layer can handle it properly
		return nil, sql.ErrNoRows
	}

	return row.ToDomain(), nil
}

func (s *Store) Create(ctx context.Context, todoList *domain.TodoList) error {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[createTodoListQuery], templateParams)
	if err != nil {
		return err
	}

	queryParams := map[string]any{
		"user_id":    todoList.UserID,
		"title":      todoList.Title,
		"color":      todoList.Color,
		"labels":     strings.Join(todoList.Labels, ","),
		"created_at": todoList.CreatedAt,
	}

	result, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return err
	}
	defer result.Close()

	var (
		id int64
	)

	if result.Next() {
		err = result.Scan(&id)
		if err != nil {
			return err
		}
	} else {
		return errors.New("failed to retrieve inserted todo list ID")
	}

	// Create a new Todo instance with the retrieved ID and other fields
	todoList.ID = id

	return nil
}

func (s *Store) Update(ctx context.Context, id int64, title string, color string, labels []string, deleted bool) (*domain.TodoList, error) {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[updateTodoListQuery], templateParams)
	if err != nil {
		return nil, err
	}

	queryParams := map[string]any{
		"id":      id,
		"title":   title,
		"color":   color,
		"labels":  strings.Join(labels, ","),
		"deleted": deleted,
	}

	result, err := s.db.NamedExecContext(ctx, querystr, queryParams)
	if err != nil {
		return nil, err
	}

	// Optional: Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		// Return sql.ErrNoRows so the service layer can handle it properly
		return nil, sql.ErrNoRows
	}

	return s.GetListByID(ctx, id)
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[deleteTodoListQuery], templateParams)
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

	// Optional: Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Return sql.ErrNoRows so the service layer can handle it properly
		return sql.ErrNoRows
	}

	return nil
}
