package tests

import (
	"clothing-store-brands/internal/data"
	s "clothing-store-brands/internal/server"
	"clothing-store-brands/pkg/config"
	"clothing-store-brands/pkg/jsonlog"
	pb "clothing-store-brands/pkg/pb"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)

var server s.Server

func init() {
	config.ReadConfig()
	db, _ := openDB()

	server = s.Server{
		Logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		Models: data.NewModels(db)}
}

func TestCreateBrand(t *testing.T) {
	res, err := server.CreateBrand(context.Background(), &pb.Brand{
		Name:        "Test",
		Country:     "Test",
		Description: "Test",
		ImageUrl:    "Test",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res.Name != "Test" {
		t.Fatalf("Unexpected name, should be Test got: %v", res.Name)
	}
	id = res.Id
}

func TestGetBrand(t *testing.T) {
	res, err := server.ShowBrand(context.Background(), &pb.ShowBrandRequest{Id: id})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if res.Name != "Test" {
		t.Fatalf("Expected to have a name Test, but got: %v", res.Name)
	}
}

func TestUpdateBrand(t *testing.T) {
	res, err := server.UpdateBrand(context.Background(), &pb.UpdateBrandRequest{
		Id:    id,
		Brand: &pb.Brand{Name: "UpdatedTest"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if res.Name != "UpdatedTest" {
		t.Fatalf("Expected to have a name UpdatedTest, but got: %v", res.Name)
	}
}

func TestGetBrandAfterUpdate(t *testing.T) {
	res, err := server.ShowBrand(context.Background(), &pb.ShowBrandRequest{Id: id})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if res.Name != "UpdatedTest" {
		t.Fatalf("Expected to have a name UpdatedTest, but got: %v", res.Name)
	}
}

func TestDeleteBrand2(t *testing.T) {
	_, err := server.DeleteBrand(context.Background(), &pb.DeleteBrandRequest{Id: id})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}

}

func TestGetBrandAfterDelete(t *testing.T) {
	_, err := server.ShowBrand(context.Background(), &pb.ShowBrandRequest{Id: id})
	if err == nil {
		t.Fatalf("Expected error but didn't get it")
	}
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.C.Db.Dsn)
	fmt.Println(config.C.Db.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.C.Db.MaxOpenConns)
	db.SetMaxIdleConns(config.C.Db.MaxIdleConns)
	duration, err := time.ParseDuration(config.C.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
