package main

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/jmoiron/sqlx"

	"github.com/macesz/todo-go/dal/pgtodo"
	"github.com/macesz/todo-go/dal/pguser"
	"github.com/macesz/todo-go/delivery/web"
	"github.com/macesz/todo-go/domain"
	infraPG "github.com/macesz/todo-go/infra/postgres"
	"github.com/macesz/todo-go/services/todo"
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

	// Create DATA STORES
	todoStore := pgtodo.CreateStore(db)
	userStore := pguser.CreateStore(db)

	// Create SERVICES
	todoService := todo.NewTodoService(todoStore) // Service with business logic
	userService := user.NewUserService(userStore) // Service with business logic

	// Create WEB HANDLERS
	handlers, err := web.CreateHandlers(ctx, &web.ServerServices{
		Todo: todoService,
		User: userService,
	})
	if err != nil {
		panic(err)
	}

	web.StartServer(ctx, cfg, handlers) // Start the web server
}
