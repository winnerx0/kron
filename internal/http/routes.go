package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/swaggo/http-swagger/v2"
	"github.com/winnerx0/kron/internal/config"
	"github.com/winnerx0/kron/internal/database"
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/job"
)

type App struct {
	config config.Config
}

func NewApp(config config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Start() error {

	database := database.NewDatabase(a.config.DBHost, a.config.DBUser, a.config.DBPassword, a.config.DBPort, a.config.DBName)

	db := database.Start()

	executionRepo := execution.NewPostgresRepository(db)

	// executionService := execution.NewExecutionService(executionRepo)

	// executionHandler := NewExecutionHandler(executionService)

	jobRepo := job.NewRepository(db)

	jobService := job.NewJobService(jobRepo, executionRepo)

	jobHandler := NewJobHandler(jobService)
	
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	r.Post("/create", jobHandler.Create)

	return http.ListenAndServe(":5000", r)
}
