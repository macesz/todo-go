package pgtodo

import (
	"testing"

	"github.com/macesz/todo-go/pkg"
)

func TestTemplate(t *testing.T) {
	queries, err := pkg.BuildQueries(files, "queries")
	if err != nil {
		t.Error(err)
	}

	query, err := pkg.PrepareQuery(queries["list_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}

func TestTemplateCreate(t *testing.T) {
	queries, err := pkg.BuildQueries(files, "queries")
	if err != nil {
		t.Error(err)
	}

	query, err := pkg.PrepareQuery(queries["create_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}

func TestTemplateGetTodo(t *testing.T) {
	queries, err := pkg.BuildQueries(files, "queries")
	if err != nil {
		t.Error(err)
	}

	query, err := pkg.PrepareQuery(queries["get_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}

func TestTemplateDeleteTodo(t *testing.T) {
	queries, err := pkg.BuildQueries(files, "queries")
	if err != nil {
		t.Error(err)
	}

	query, err := pkg.PrepareQuery(queries["delete_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}

func TestTemplateUpdateTodo(t *testing.T) {
	queries, err := pkg.BuildQueries(files, "queries")
	if err != nil {
		t.Error(err)
	}

	query, err := pkg.PrepareQuery(queries["update_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}
