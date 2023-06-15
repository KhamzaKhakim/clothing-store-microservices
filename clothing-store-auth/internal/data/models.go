package data

import (
	"database/sql"
)

type Models struct {
	Users  UserModel
	Roles  RolesModel
	Tokens TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:  UserModel{DB: db},
		Roles:  RolesModel{DB: db},
		Tokens: TokenModel{DB: db},
	}
}
