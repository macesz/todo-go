package main

import (
	"github.com/macesz/todo-go/dal/inmemorytodo"
	"github.com/macesz/todo-go/delivery/web"
	"github.com/macesz/todo-go/services/todo"
)

func main() {
	store := inmemorytodo.NewInMemoryStore() // Our in-memory store (like a database)
	service := todo.NewTodoService(store)    // Service with business logic

	web.StartServer(service) // Start the web server
}
