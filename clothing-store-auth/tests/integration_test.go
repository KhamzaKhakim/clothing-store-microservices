package tests

import (
	"clothing-store-auth/internal/data"
	s "clothing-store-auth/internal/server"
	"clothing-store-auth/pkg/config"
	"clothing-store-auth/pkg/jsonlog"
	pb "clothing-store-auth/pkg/pb"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
)

var server s.Server
var adminToken string
var newUser pb.RegisterResponse

func init() {
	config.ReadConfig()
	db, _ := openDB()

	server = s.Server{
		Logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		Models: data.NewModels(db)}
}

func TestLogin(t *testing.T) {
	_, err := server.Login(context.Background(), &pb.LoginRequest{
		Email:    "admin@gmail.com",
		Password: "123456",
	})
	if err == nil {
		t.Fatalf("Expected to get error but everuthing was okay")
	}

	res, err := server.Login(context.Background(), &pb.LoginRequest{
		Email:    "admin@gmail.com",
		Password: "12345678",
	})
	if err != nil {
		t.Fatalf("Got unexpected error")
	}
	adminToken = res.Token
}

func TestRegister(t *testing.T) {
	res, err := server.Register(context.Background(),
		&pb.RegisterRequest{
			Name:     "Test",
			Email:    "test789467@gmail.com",
			Password: "12345678",
		})
	if err != nil {
		t.Fatalf("Expected to be ok, but got error: %v", err)
	}
	newUser = *res
}

func TestAuthenticate(t *testing.T) {
	_, err := server.Authenticate(context.Background(), &pb.AuthenticateRequest{Token: adminToken})
	if err == nil {
		t.Fatalf("Expected to be okay but got error")
	}
}

func TestAuthorize(t *testing.T) {
	res, err := server.Authorize(context.Background(), &pb.AuthorizeRequest{Id: 1, Role: "ADMIN"})
	if err != nil {
		t.Fatalf("Expected to not get error, since user with id 1 is admin")
	}
	if res.Msg == "" {
		t.Fatalf("Expected to get message, not empty string")
	}
	_, err = server.Authorize(context.Background(), &pb.AuthorizeRequest{Id: 1, Role: "USER"})
	if err == nil {
		t.Fatalf("Expected to get error, since user with id 1 is admin")
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
