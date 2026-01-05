package web

import (
	"context"
	"io"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/macesz/todo-go/delivery/web/todo"
	"github.com/macesz/todo-go/delivery/web/todolist"
	"github.com/macesz/todo-go/delivery/web/user"
)

type ServerServices struct {
	TodoList  todolist.TodoListService
	Todo      todo.TodoService
	User      user.UserService
	TokenAuth *jwtauth.JWTAuth
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

type Handlers struct {
	TodoList *todolist.TodoListHandlers
	Todo     *todo.TodoHandlers
	User     *user.UserHandlers
}

func CreateHandlers(ctx context.Context, services *ServerServices) (*Handlers, error) {
	todoListHandler := todolist.NewHandlers(services.TodoList, services.Todo, services.User)
	todoHandler := todo.NewHandlers(services.Todo, services.User)      // Create handlers with the service
	userHandler := user.NewHandlers(services.User, services.TokenAuth) // Create handlers with the service

	handlers := &Handlers{
		TodoList: todoListHandler,
		Todo:     todoHandler,
		User:     userHandler,
	}

	return handlers, nil
}
