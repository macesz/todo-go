package pgtodo

import (
	"context"
	"errors"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/macesz/todo-go/domain"
)

// Here is the Store struct where we store the queries and the database connection.
type Store struct {
	queryTemplates map[string]*template.Template
	db             *sqlx.DB
}

// CreateStore creates a new Store instance.
func CreateStore(db *sqlx.DB) *Store {
	queryTemplates, err := buildQueries("queries")
	if err != nil {
		panic(err)
	}

	return &Store{
		queryTemplates: queryTemplates,
		db:             db,
	}
}

// ListTodo retrieves a list of todos from the database.
func (s *Store) ListTodo(ctx context.Context) ([]domain.Todo, error) {
	todos := make([]domain.Todo, 0)

	// Template parameters are not safe to use directly in the query, because they can be used to inject SQL code.
	// I can use anything that is not a user input, like Table Name, Column Name, etc.
	templateParams := map[string]any{}

	// Prepare the query string, by using the template.
	querystr, err := prepareQuery(s.queryTemplates[listTodoQuery], templateParams)
	if err != nil {
		return nil, err
	}

	// Prepare the query parameters.
	// This is safe to use directly in the query, because it uses named parameters.
	queryParams := map[string]any{}

	// Execute the query. You can add parameters to the query if needed instead of using nil.
	//NamedQueryContext ✅ - Multiple rows (ListTodos, Search, etc.)
	rows, err := s.db.NamedQueryContext(ctx, querystr, queryParams)
	if err != nil {
		return nil, err
	}

	defer rows.Close() // Important: Always close rows!

	var row RowDTO

	for rows.Next() {
		err := rows.StructScan(&row) // Fixed: Added & (pointer) and error handling
		if err != nil {
			return nil, err
		}

		todos = append(todos, domain.Todo{
			ID:        row.ID,
			Title:     row.Title,
			Done:      row.Done,
			CreatedAt: row.CreatedAt,
		})
	}

	return todos, nil
}

func (s *Store) CreateTodo(ctx context.Context, todo *domain.Todo) (int, error) {
	templateParams := map[string]any{}

	querystr, err := prepareQuery(s.queryTemplates[createTodoQuery], templateParams)
	if err != nil {
		return 0, err
	}

	queryParams := map[string]any{
		"title": todo.Title,
	}

	var id int
	// GetContext ✅ - Single row to add a new todo
	err = s.db.GetContext(ctx, &id, querystr, queryParams)
	if err != nil {
		return 0, err
	}

	todo.ID = id

	return id, nil
}

func (s *Store) GetTodo(ctx context.Context, id int) (*domain.Todo, error) {
	templateParams := map[string]any{}

	querystr, err := prepareQuery(s.queryTemplates[getTodoQuery], templateParams)
	if err != nil {
		return nil, err
	}

	queryParams := map[string]any{
		"id": id,
	}

	var todo domain.Todo
	//QueryRowxContext ✅ - Single row (GetTodo, GetUser, etc.)
	err = s.db.QueryRowxContext(ctx, querystr, queryParams).StructScan(&todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *Store) UpdateTodo(ctx context.Context, todo domain.Todo) error {
	templateParams := map[string]any{}

	querystr, err := prepareQuery(s.queryTemplates[updateTodoQuery], templateParams)
	if err != nil {
		return err
	}

	queryParams := map[string]any{
		"id":    todo.ID,
		"title": todo.Title,
		"done":  todo.Done,
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

func (s *Store) DeleteTodo(ctx context.Context, id int) error {
	templateParams := map[string]any{}

	querystr, err := prepareQuery(s.queryTemplates[deleteTodoQuery], templateParams)
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
