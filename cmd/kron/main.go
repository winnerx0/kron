package main

import (
	"log"

	_ "github.com/winnerx0/kron/docs"
	"github.com/winnerx0/kron/internal/config"
	"github.com/winnerx0/kron/internal/http"
)

// @title Kron API
// @version 1.0
// @description This is a simple API for scheduling and executing jobs.
// @host localhost:5000
// @BasePath /
func main() {

	cfg := config.Load()
	app := http.NewApp(cfg)

	if err := app.Start(); err != nil {
		log.Fatal("Failed to start server ", err)
	}
}
