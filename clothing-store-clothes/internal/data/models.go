package data

import (
	"database/sql"
)

type Models struct {
	Clothes ClotheModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Clothes: ClotheModel{DB: db},
	}
}
