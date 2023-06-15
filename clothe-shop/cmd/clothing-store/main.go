package main

import (
	"clothing-store/internal/data"
	"clothing-store/internal/handlers"
	"clothing-store/internal/jsonlog"
	"clothing-store/pkg/config"
	"context"      // New import
	"database/sql" // New import
	"os"
	"time"
	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)

const version = "1.0.0"

var app1 handlers.Application

func main() {
	config.ReadConfig()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(config.C)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	logger.PrintInfo("database connection pool established", nil)

	app := &handlers.Application{
		Config: config.C,
		Logger: logger,
		Models: data.NewModels(db),
	}

	err = app.Serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
