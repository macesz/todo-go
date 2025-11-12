package testutils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/services/todo/mocks"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DbName = "test_db"
	DbUser = "test_user"
	DbPass = "test_password"
)

// TestContainer holds the container and connection info
type TestContainer struct {
	Container testcontainers.Container
	DB        *sqlx.DB
	DSN       string
}

// SetupTestDB creates a PostgreSQL container and runs migrations
func SetupTestDB(t *testing.T) *TestContainer {
	t.Helper()

	ctx := context.Background()

	// Create PostgreSQL container, same as  docker run -e POSTGRES_PASSWORD="paas" -e POSTGRES_USER="user" -p 5432:54320 postgres:14-alpine
	// it will connect to my docker
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": DbPass,
			"POSTGRES_USER":     DbUser,
			"POSTGRES_DB":       DbName,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithStartupTimeout(60 * time.Second),
	}

	//here we execute the request
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "failed to start container")

	// Cleanup container when test completes
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	})

	// Get mapped port
	mappedPort, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err, "failed to get container port")

	log.Printf("PostgreSQL container ready on port: %s", mappedPort.Port())

	// Small delay to ensure DB is fully ready
	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", mappedPort.Port())

	// Run migrations
	err = runMigrations(dbAddr)
	require.NoError(t, err, "failed to run migrations")

	// Connect to database
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)
	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err, "failed to connect to database")

	// Cleanup DB connection when test completes
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	})

	return &TestContainer{
		Container: container,
		DB:        db,
		DSN:       dsn,
	}
}

// runMigrations runs database migrations
func runMigrations(dbAddr string) error {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get caller path")
	}

	// Adjust this path to match your project structure
	projectRoot := filepath.Join(filepath.Dir(path), "..", "..")
	migrationsPath := filepath.Join(projectRoot, "infra", "postgres", "migrations")

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	log.Printf("Looking for migrations at: %s", absPath)

	// Check if directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", absPath)
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)

	// Use file:// protocol for the source URL
	sourceURL := fmt.Sprintf("file://%s", absPath)
	log.Printf("Migration source URL: %s", sourceURL)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// setupMockStore creates and configures a mock store with cleanup
func SetupMockStore(t *testing.T) *mocks.TodoStore {
	t.Helper()
	store := mocks.NewTodoStore(t)
	t.Cleanup(func() { store.AssertExpectations(t) })
	return store
}

// Helper function to add user context to request (simulating authenticated user)
func WithUserContext(req *http.Request, userID int64) *http.Request {
	userCtx := &auth.UserContext{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	ctx := userCtx.AddToContext(req.Context())
	return req.WithContext(ctx)
}

// CleanupDB cleans all tables for a fresh test state
func CleanupDB(t *testing.T, db *sqlx.DB) {
	t.Helper()

	tables := []string{"todos", "users"} // Add your table names
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			t.Logf("Warning: failed to truncate table %s: %v", table, err)
		}
	}
}
