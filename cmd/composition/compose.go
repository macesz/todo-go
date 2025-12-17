package composition

import (
	"github.com/jmoiron/sqlx"
	"github.com/macesz/todo-go/dal/pgtodo"
	"github.com/macesz/todo-go/dal/pgtodolist"
	"github.com/macesz/todo-go/dal/pguser"
	"github.com/macesz/todo-go/delivery/web"
	"github.com/macesz/todo-go/delivery/web/auth"
	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/services/todo"
	"github.com/macesz/todo-go/services/todolist"
	"github.com/macesz/todo-go/services/user"
)

func ComposeServices(cfg domain.Config, db *sqlx.DB) *web.ServerServices {
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

	return services
}
