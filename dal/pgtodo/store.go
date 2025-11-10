package pgtodo

import (
	"context"
	"errors"
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

// List retrieves a list of todos from the database.
func (s *Store) List(ctx context.Context, userID int64) ([]*domain.Todo, error) {
	todos := make([]*domain.Todo, 0)

	// Template parameters are not safe to use directly in the query, because they can be used to inject SQL code.
	// I can use anything that is not a user input, like Table Name, Column Name, etc.
	templateParams := map[string]any{}

	// Prepare the query string, by using the template.
	querystr, err := pkg.PrepareQuery(s.queryTemplates[listTodoQuery], templateParams)
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
		err := rows.StructScan(&row) // Fixed: Added & (pointer) and error handling
		if err != nil {
			return nil, err
		}

		todos = append(todos, row.ToDomain())
	}

	return todos, nil
}

func (s *Store) Create(ctx context.Context, todo *domain.Todo) error {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[createTodoQuery], templateParams)
	if err != nil {
		return err
	}

	queryParams := map[string]any{
		"user_id":    todo.UserID,
		"title":      todo.Title,
		"priority":   todo.Priority,
		"created_at": todo.CreatedAt,
	}

	// NamedQueryContext ✅ - Single row with RETURNING clause
	result, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return err
	}
	defer result.Close()

	var (
		id int64
	)

	// Scan the result into the variables
	if result.Next() {
		err = result.Scan(&id)
		if err != nil {
			return err
		}
	} else {
		return errors.New("failed to retrieve inserted todo ID")
	}

	// Create a new Todo instance with the retrieved ID and other fields
	todo.ID = id

	return nil
}

func (s *Store) Get(ctx context.Context, id int64) (*domain.Todo, error) {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[getTodoQuery], templateParams)
	if err != nil {
		return nil, err
	}

	queryParams := map[string]any{
		"id": id,
	}

	var row rowDTO
	//NamedQueryContext ✅ - Single row with named parameters (GetTodo, GetUser, etc.)
	rows, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return nil, err
	}

	// don't forget to close the rows
	defer rows.Close()

	// Scan the row into the todo struct, first call `Next()` and then `StructScan()` to get the data from the result
	if rows.Next() {
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("todo not found")
	}

	return row.ToDomain(), nil
}

func (s *Store) Update(ctx context.Context, id int64, title string, done bool, priority int64) (*domain.Todo, error) {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[updateTodoQuery], templateParams)
	if err != nil {
		return nil, err
	}

	queryParams := map[string]any{
		"id":       id,
		"title":    title,
		"done":     done,
		"priority": priority,
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
		return nil, errors.New("todo not found")
	}

	return s.Get(ctx, id)
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	templateParams := map[string]any{}

	querystr, err := pkg.PrepareQuery(s.queryTemplates[deleteTodoQuery], templateParams)
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
		return errors.New("todo not found")
	}

	return nil
}
