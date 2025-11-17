package pgtodolist

import (
	"embed"
)

//go:embed queries/*.sql.tpl
var files embed.FS

const (
	listTodoListQuery   = "list_todo_list"
	createTodoListQuery = "create_todo_list"
	getTodoListQuery    = "get_todo_list"
	updateTodoListQuery = "update_todo_list"
	deleteTodoListQuery = "delete_todo_list"
)
