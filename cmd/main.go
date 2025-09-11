package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/macesz/todo-go/dal/infiletodo"
	"github.com/macesz/todo-go/delivery/web"
	"github.com/macesz/todo-go/services/todo"
)

func main() {
	dir := os.Getenv("DATA_DIR")
	if dir == "" {
		dir = "data"
	}
	// Our in-memory store (like a database)
	// store := inmemorytodo.NewInMemoryStore()

	// Using file-based store
	filePath := filepath.Join(dir, "todos.csv")

	store, err := infiletodo.NewInFileStore(filePath)
	if err != nil {
		log.Fatal(err)
	}

	service := todo.NewTodoService(store) // Service with business logic

	web.StartServer(service) // Start the web server
}
