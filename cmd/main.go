package main

import (
	"fmt"
	"os"

	// "github.com/macesz/todo-go/dal/infiletodo"
	"github.com/jmoiron/sqlx"
	"github.com/macesz/todo-go/dal/pgtodo"
	"github.com/macesz/todo-go/delivery/web"
	"github.com/macesz/todo-go/services/todo"
)

func main() {
	// dir := os.Getenv("DATA_DIR")
	// if dir == "" {
	// 	dir = "data"
	// }
	// // Our in-memory store (like a database)
	// store := inmemorytodo.NewInMemoryStore()

	// Using file-based store
	// filePath := filepath.Join(dir, "todos.csv")

	// store, err := infiletodo.NewInFileStore(filePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// PG CONNECTION
	DbUser := os.Getenv("DB_USER")
	DbPass := os.Getenv("DB_PASS")
	dbAddr := os.Getenv("DB_ADDR")
	DbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}

	store := pgtodo.CreateStore(db)

	service := todo.NewTodoService(store) // Service with business logic

	web.StartServer(service, nil) // Start the web server
}
