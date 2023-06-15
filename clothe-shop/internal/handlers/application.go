package handlers

import (
	"clothing-store/internal/data"
	"clothing-store/internal/jsonlog"
	"clothing-store/pkg/config"
	"sync"
)

type Application struct {
	Config config.Config
	Logger *jsonlog.Logger
	Models data.Models
	Wg     sync.WaitGroup
}
