package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/winnerx0/kron/internal/config"
	"github.com/winnerx0/kron/internal/database"
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/job"
	"github.com/winnerx0/kron/internal/secret"
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

	executionService := execution.NewExecutionService(executionRepo)

	executionHandler := NewExecutionHandler(executionService)

	jobRepo := job.NewRepository(db)

	secretManager, err := secret.NewAESGCMManager(a.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to initialize secret encryption: %w", err)
	}

	jobService := job.NewJobService(jobRepo, executionRepo, secretManager)

	jobHandler := NewJobHandler(jobService)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {

		r.Route("/job", func(r chi.Router) {

			r.Post("/create", jobHandler.Create)

			r.Put("/{jobID}", jobHandler.UpdateJob)

			r.Delete("/{jobID}", jobHandler.DeleteJob)

			r.Get("/all", jobHandler.FindAll)
		})

		r.Route("/execution", func(r chi.Router) {

			r.Get("/all", executionHandler.FindAll)
		})

	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Listening to server on port 5000")
	return http.ListenAndServe(":5000", r)
}
