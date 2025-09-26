package pgtodo

import (
	"bytes"
	"embed"
	"path/filepath"
	"strings"
	"text/template"
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

func buildQueries(dir string) (map[string]*template.Template, error) {
	queries := make(map[string]*template.Template)

	tmpfiles, err := files.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, tmpf := range tmpfiles {
		if tmpf.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, filepath.Join(dir, tmpf.Name()))
		if err != nil {
			continue
		}

		queries[strings.Split(tmpf.Name(), ".")[0]] = pt
	}

	return queries, nil
}

func prepareQuery(queryTpl *template.Template, params any) (string, error) {
	var buff bytes.Buffer

	err := queryTpl.Execute(&buff, params)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
