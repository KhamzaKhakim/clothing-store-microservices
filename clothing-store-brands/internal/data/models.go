package data

import (
	"database/sql"
)

type Models struct {
	Brands BrandModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Brands: BrandModel{DB: db},
	}
}
