package pguser

import (
	"embed"
)

//go:embed queries/*.sql.tpl
var files embed.FS

const (
	createUserQuery = "create_user"
	getUserQuery    = "get_user"
	deleteUserQuery = "delete_user"
)
