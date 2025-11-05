package pkg

import (
	"bytes"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

func BuildQueries(files fs.ReadDirFS, dir string) (map[string]*template.Template, error) {
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

func PrepareQuery(queryTpl *template.Template, params any) (string, error) {
	var buff bytes.Buffer

	err := queryTpl.Execute(&buff, params)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
