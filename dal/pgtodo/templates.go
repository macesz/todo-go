package pgtodo

import (
	"embed"
)

//go:embed queries/*.sql.tpl
var files embed.FS

const (
	listTodoQuery   = "list_todo"
	createTodoQuery = "create_todo"
	getTodoQuery    = "get_todo"
	updateTodoQuery = "update_todo"
	deleteTodoQuery = "delete_todo"
)
