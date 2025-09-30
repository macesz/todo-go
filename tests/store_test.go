package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres" // used by migrator
	_ "github.com/golang-migrate/migrate/source/file"       // used by migrator
	"github.com/jmoiron/sqlx"
	"github.com/macesz/todo-go/dal/pgtodo"
	"github.com/macesz/todo-go/domain"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DbName = "test_db"
	DbUser = "test_user"
	DbPass = "test_password"
)

func SetUpTestDB(ctx context.Context) (*sqlx.DB, error) {
	// create postgres container
	var env = map[string]string{
		"POSTGRES_PASSWORD": DbPass,
		"POSTGRES_USER":     DbUser,
		"POSTGRES_DB":       DbName,
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())

	// execute migration
	migrateErr := migrateDb(dbAddr)
	if migrateErr != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", migrateErr)
	}

	// Connect directly with sqlx - MUCH SIMPLER!
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDb(dbAddr string) error {
	// get location of test
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get path")
	}

	pathToMigrationFiles := filepath.Join(filepath.Dir(path), "..", "infra", "postgres", "migrations")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)

	m, err := migrate.New(fmt.Sprintf("file:%s", pathToMigrationFiles), databaseURL)
	if err != nil {
		return err
	}

	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migration done")

	return nil
}

func TestPgTodoStore(t *testing.T) {
	// ctx := t.Context()

	// db, err := SetUpTestDB(ctx)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// store := pgtodo.CreateStore(db)

	tests := []struct {
		name string
		exec func(*testing.T)
	}{
		{
			name: "create 3 todos",
			exec: func(t *testing.T) {
				store, db := setupTestStore(t)

				defer db.Close()

				ctx := t.Context()

				for i := range int64(3) {
					ct, err := store.Create(ctx, fmt.Sprintf("test%d", i+1))
					if err != nil {
						t.Error(err)
					}

					if ct.ID != i+1 {
						t.Errorf("expected id to be %d, got %d", i+1, ct.ID)
					}

					if ct.Title != fmt.Sprintf("test%d", i+1) {
						t.Errorf("expected title to be 'test%d', got %s", i+1, ct.Title)
					}

					if ct.Done != false {
						t.Errorf("expected completed to be false, got %t", ct.Done)
					}
				}

			},
		},
		{
			name: "ListTodo",
			exec: func(t *testing.T) {
				store, db := setupTestStore(t)
				defer db.Close()

				// Create sample data for this specific test
				createSampleTodos(t, store, 3)

				ctx := t.Context()

				todos, err := store.List(ctx)
				if err != nil {
					t.Error(err)
				}

				if len(todos) != 3 {
					t.Errorf("expected 3 todos, got %d", len(todos))
				}
			},
		},
		{
			name: "GetTodo",
			exec: func(t *testing.T) {
				store, db := setupTestStore(t)
				defer db.Close()

				// Create sample data for this specific test
				createSampleTodos(t, store, 3)

				ctx := t.Context()

				todo, err := store.Get(ctx, 1)
				if err != nil {
					t.Error(err)
				}

				if todo.ID != 1 {
					t.Errorf("expected id to be 1, got %d", todo.ID)
				}

				if todo.Title != "test1" {
					t.Errorf("expected title to be 'test1', got %s", todo.Title)
				}
			},
		},
		{
			name: "UpdateTodo",
			exec: func(t *testing.T) {
				store, db := setupTestStore(t)
				defer db.Close()

				// Create sample data for this specific test
				createSampleTodos(t, store, 3)

				ctx := t.Context()

				var todoId int64
				todoId = 1
				todo, err := store.Get(ctx, todoId)
				if err != nil {
					t.Error(err)
				}

				todo.Title = "test1 updated"
				todo.Done = true

				todo, err = store.Update(ctx, todoId, todo.Title, todo.Done)
				if err != nil {
					t.Error(err)
				}

				todo, err = store.Get(ctx, 1)
				if err != nil {
					t.Error(err)
				}

				if todo.Title != "test1 updated" {
					t.Errorf("expected title to be 'test1 updated', got %s", todo.Title)
				}

				if todo.Done != true {
					t.Errorf("expected completed to be true, got %t", todo.Done)
				}
			},
		},
		{
			name: "DeleteTodo",
			exec: func(t *testing.T) {
				store, db := setupTestStore(t)
				defer db.Close()

				// Create sample data for this specific test
				createSampleTodos(t, store, 3)

				ctx := t.Context()

				err := store.Delete(ctx, 1)
				if err != nil {
					t.Error(err)
				}

				todo, err := store.Get(ctx, 1)
				if err == nil {
					t.Error(fmt.Errorf("expected error, got nil"))
				}

				if todo != nil {
					t.Errorf("expected todo to be nil, got %v", todo)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.exec(t)
		})
	}
}

func setupTestStore(t *testing.T) (*pgtodo.Store, *sqlx.DB) {
	ctx := t.Context()

	db, err := SetUpTestDB(ctx)
	if err != nil {
		t.Fatal(err)
	}

	store := pgtodo.CreateStore(db)
	return store, db
}

func createSampleTodos(t *testing.T, store *pgtodo.Store, count int) []*domain.Todo {

	ctx := t.Context()
	var todos []*domain.Todo

	for i := range int64(count) {
		todo, err := store.Create(ctx, fmt.Sprintf("test%d", i+1))
		if err != nil {
			t.Error(err)
		}
		todos = append(todos, todo)
	}

	return todos
}
