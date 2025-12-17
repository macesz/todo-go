package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/jmoiron/sqlx"

	"github.com/macesz/todo-go/cmd/composition"
	"github.com/macesz/todo-go/delivery/web"
	"github.com/macesz/todo-go/domain"
	infraPG "github.com/macesz/todo-go/infra/postgres"
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

	services := composition.ComposeServices(cfg, db)

	// Create WEB HANDLERS
	handlers, err := web.CreateHandlers(ctx, services)
	if err != nil {
		panic(err)
	}

	// Pass services to both handlers AND router
	router, err := web.CreateRouter(ctx, cfg, services, handlers)
	if err != nil {
		log.Fatal(err)
	}

	// Start the server
	log.Printf("listening on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServerPort), router); err != nil {
		log.Fatal(err)
	}
}

// This follows Dependency Inversion Principle - high-level modules (server) depend on abstractions (services struct)
// rather than creating dependencies internally.
