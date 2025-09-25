package pgtodo

import "testing"

func TestTemplate(t *testing.T) {
	queries, err := buildQueries("queries")
	if err != nil {
		t.Error(err)
	}

	query, err := prepareQuery(queries["list_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}

func TestTemplateCreate(t *testing.T) {
	queries, err := buildQueries("queries")
	if err != nil {
		t.Error(err)
	}

	query, err := prepareQuery(queries["create_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}

func TestTemplateGetTodo(t *testing.T) {
	queries, err := buildQueries("queries")
	if err != nil {
		t.Error(err)
	}

	query, err := prepareQuery(queries["get_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}

func TestTemplateDeleteTodo(t *testing.T) {
	queries, err := buildQueries("queries")
	if err != nil {
		t.Error(err)
	}

	query, err := prepareQuery(queries["delete_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}

func TestTemplateUpdateTodo(t *testing.T) {
	queries, err := buildQueries("queries")
	if err != nil {
		t.Error(err)
	}

	query, err := prepareQuery(queries["update_todo"], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(query)
}
