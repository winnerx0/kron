package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	_ "github.com/winnerx0/kron/docs"
	"github.com/winnerx0/kron/internal/config"
	"github.com/winnerx0/kron/internal/database"
	"github.com/winnerx0/kron/internal/http"
	"github.com/winnerx0/kron/internal/job"
	rabbitmq "github.com/winnerx0/kron/internal/queue"
	"github.com/winnerx0/kron/internal/scheduler"
)

// @title Kron API
// @version 1.0
// @description This is a simple API for scheduling and executing jobs.
// @host localhost:5000
// @BasePath /
func main() {

	cfg := config.Load()
	app := http.NewApp(cfg)

	conn := rabbitmq.NewRabbitMQClient(cfg.RabbitMQURL)

	err := rabbitmq.Setup(conn.Ch)

	if err != nil {
		log.Fatal("Failed to setup RabbitMQ Queues ", err)
	}

	database := database.NewDatabase(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort, cfg.DBName)

	db := database.Start()

	jobRepo := job.NewRepository(db)

	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379", DB: 0, Protocol: 2})

	ctx := context.Background()
	
	publisher, err := scheduler.NewPublisher(conn.Ch, "jobs_queue")
	
	if err != nil {
		log.Fatal("Failed to create publisher", err)
	}
	
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal("Failed to get hostname", err)
	}

	locker := scheduler.NewLocker(rdb, time.Minute * 5, hostname)

	poller := scheduler.NewPoller(jobRepo, publisher, locker, time.Second * 30)

	go poller.Run(ctx)
	
	if err := app.Start(); err != nil {
		log.Fatal("Failed to start server ", err)
	}
}
