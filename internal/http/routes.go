package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/winnerx0/kron/internal/auth"
	"github.com/winnerx0/kron/internal/config"
	"github.com/winnerx0/kron/internal/dashboard"
	"github.com/winnerx0/kron/internal/database"
	"github.com/winnerx0/kron/internal/domain"
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/job"
	"github.com/winnerx0/kron/internal/middleware"
	"github.com/winnerx0/kron/internal/oauth"
	refreshtoken "github.com/winnerx0/kron/internal/refresh_token"
	"github.com/winnerx0/kron/internal/secret"
	"github.com/winnerx0/kron/internal/user"
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

	userRepo := user.NewRepository(db)

	userService := user.NewService(userRepo)

	userHandler := NewUserHandler(userService)

	refreshTokenRepo := refreshtoken.NewRepository(db)

	jobRepo := job.NewRepository(db)

	secretManager, err := secret.NewAESGCMManager(a.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to initialize secret encryption: %w", err)
	}

	jobService := job.NewJobService(jobRepo, executionRepo, secretManager)

	dashboardRepo := dashboard.NewPostgresRepository(db)

	dashboardService := dashboard.NewService(dashboardRepo)

	dashboardHandler := NewDashboardHandler(dashboardService)

	jobsCh := make(chan domain.Job, 10)

	workerCount := 5

	func(jobs <-chan domain.Job) {
		for range workerCount {
			go func() {
				log.Println("Worker started and waiting for jobs")
				for j := range jobs {
					jobService.ExecuteJob(context.Background(), j, true)
				}

			}()
		}
	}(jobsCh)

	go jobService.RunJobs(context.Background(), jobsCh)

	jobHandler := NewJobHandler(jobService)

	authService := auth.NewService(&a.config, userRepo, refreshTokenRepo)

	authHandler := NewAuthHandler(authService)

	oauthService := oauth.NewService(&a.config, userRepo, authService)

	oauthHandler := NewOauthHandler(a.config.FrontendURL, oauthService)

	authMiddleware := middleware.NewAuthMiddleware(&a.config, userRepo)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://kron-pi.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {

		r.Route("/auth", func(r chi.Router) {
			r.Post("/refresh", authHandler.Refresh)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Use)

			r.Get("/dashboard", dashboardHandler.Summary)

			r.Route("/job", func(r chi.Router) {

				r.Post("/create", jobHandler.Create)

				r.Put("/{jobID}", jobHandler.UpdateJob)

				r.Delete("/{jobID}", jobHandler.DeleteJob)

				r.Get("/all", jobHandler.FindAll)

				r.Post("/{jobID}/run", jobHandler.RunJob)

				r.Post("/{jobID}/stop", jobHandler.StopJob)

				r.Patch("/{jobID}/status", jobHandler.UpdateJobStatus)
			})

			r.Route("/execution", func(r chi.Router) {
				r.Get("/all", executionHandler.FindAll)
				r.Get("/{id}", executionHandler.FindByID)
			})

			r.Route("/user", func(r chi.Router) {
				r.Get("/me", userHandler.Me)
			})
		})

	})

	r.Route("/oauth", func(r chi.Router) {
		r.Route("/google", func(r chi.Router) {
			r.Get("/login", oauthHandler.GoogleLogin)
			r.Get("/callback", oauthHandler.GoogleCallback)
		})
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Listening to server on port 5000")
	return http.ListenAndServe(":5000", r)
}
