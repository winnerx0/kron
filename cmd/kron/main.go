package main

import (
	"context"
	"log"

	"github.com/winnerx0/kron/internal/config"
	"github.com/winnerx0/kron/internal/database"
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/http"
	"github.com/winnerx0/kron/internal/job"
	_"github.com/winnerx0/kron/docs"
)

// @title Kron API
// @version 1.0
// @description This is a simple API for scheduling and executing jobs.
// @host localhost:5000
// @BasePath /
func main() {
	
	cfg := config.Load()
	app := http.NewApp(cfg)
	
	database := database.NewDatabase(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort, cfg.DBName)

	db := database.Start()
	
	jobRepo := job.NewRepository(db)
	
	executionRepo := execution.NewPostgresRepository(db)
	
	jobService := job.NewJobService(jobRepo, executionRepo)

	go jobService.RunJobs(context.Background())
	
	err := app.Start()
	if err != nil {
		log.Fatal("Failed to start server ", err)
	}
}
