package pguser

import (
	"embed"
)

//go:embed queries/*.sql.tpl
var files embed.FS

const (
	createUserQuery     = "create_user"
	getUserQuery        = "get_user"
	getUserByEmailQuery = "get_user_by_email"
	deleteUserQuery     = "delete_user"
	loginUserQuery      = "login_user"
)
