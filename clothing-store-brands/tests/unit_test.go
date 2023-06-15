package tests

import (
	"clothing-store-brands/internal/data"
	s "clothing-store-brands/internal/server"
	"clothing-store-brands/internal/validator"
	"clothing-store-brands/pkg/config"
	"clothing-store-brands/pkg/jsonlog"
	"fmt"
	"os"
	"testing"
	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)

var brand data.Brand
var brandsList []*data.Brand
var id int64

func init() {
	config.ReadConfig()
	db, err := openDB()
	fmt.Println("------------")
	fmt.Println(err)
	fmt.Println("------------")

	server = s.Server{
		Logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		Models: data.NewModels(db)}
}

func TestValidateBrand(t *testing.T) {
	brand := data.Brand{
		ID:          1,
		Name:        "",
		Country:     "Jojo",
		Description: "Jojo",
		ImageURL:    "JOjo",
	}
	v := validator.New()
	if data.ValidateBrand(v, &brand); v.Valid() {
		t.Fatalf(`Expected to be not valid but got: %v`, v.Valid())
	}

	v1 := validator.New()

	brand.Name = "Jojo"
	if data.ValidateBrand(v1, &brand); !v1.Valid() {
		t.Fatalf(`Expected to be valid but got: %v`, v1.Valid())
	}

}

func TestInsertBrand(t *testing.T) {
	newBrand := &data.Brand{
		Name:        "Test",
		Country:     "Test",
		Description: "Test",
		ImageURL:    "Test",
	}
	err := server.Models.Brands.Insert(newBrand)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf("Expected to be okay, but got error: %v", err.Error())
	}
	id = newBrand.ID
}

func TestShowBrand(t *testing.T) {
	res, err := server.Models.Brands.Get(id)
	if err != nil {
		t.Fatalf("Expected not to get error but got: %v", err.Error())
	}
	if res.Name != "Test" {
		t.Fatalf("Unexpected response name. Expected Test but got: %v", res.Name)
	}
}

func TestListBrandBeforeDelete(t *testing.T) {
	brands, err := server.Models.Brands.GetAll()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	brandsList = brands
}

func TestDeleteBrand(t *testing.T) {
	err := server.Models.Brands.Delete(id)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestListBrandAfterDelete(t *testing.T) {
	brands, err := server.Models.Brands.GetAll()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if len(brands) != len(brandsList)-1 {
		t.Fatalf("Unexpected deletion or listing error")
	}
}
