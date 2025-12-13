package main

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/jmoiron/sqlx"

	"github.com/macesz/todo-go/dal/pgtodo"
	"github.com/macesz/todo-go/dal/pgtodolist"
	"github.com/macesz/todo-go/dal/pguser"
	"github.com/macesz/todo-go/delivery/web"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/domain"
	infraPG "github.com/macesz/todo-go/infra/postgres"
	"github.com/macesz/todo-go/services/todo"
	"github.com/macesz/todo-go/services/todolist"
	"github.com/macesz/todo-go/services/user"
)

func main() {
	ctx := context.Background()

	// Load CONFIG from ENV variables
	cfg := domain.Config{
		DBAddr:     os.Getenv("DB_ADDR"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASS"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}

	// Connect to POSTGRESQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBAddr,
		cfg.DBName)

	// check arg contains --migrate
	if slices.Contains(os.Args, "migrate") {
		if err := infraPG.MigrateDb(cfg.DBUser, cfg.DBPassword, cfg.DBAddr, cfg.DBName); err != nil {
			panic(err)
		}
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Create DATA STORES
	todoStore := pgtodo.CreateStore(db)
	todolistStore := pgtodolist.CreateStore(db)
	userStore := pguser.CreateStore(db)

	// Create SERVICES
	// NEW: Create auth at application startup
	tokenAuth := auth.CreateTokenAuth(cfg.JWTSecret)
	todoService := todo.NewTodoService(todoStore) // Service with business logic
	todoListService := todolist.NewTodoListService(todolistStore)
	userService := user.NewUserService(userStore) // Service with business logic

	services := &web.ServerServices{
		TodoList:  todoListService,
		Todo:      todoService,
		User:      userService,
		TokenAuth: tokenAuth, // ← Injected dependency
	}

	// Create WEB HANDLERS
	handlers, err := web.CreateHandlers(ctx, services)
	if err != nil {
		panic(err)
	}

	// Pass services to both handlers AND server
	web.StartServer(ctx, cfg, services, handlers) // Start the web server
}

// This follows Dependency Inversion Principle - high-level modules (server) depend on abstractions (services struct)
// rather than creating dependencies internally.
