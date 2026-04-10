package job

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/winnerx0/kron/internal/execution"
)

type Service struct {
	repo          Repository
	executionRepo execution.Repository
}

func NewJobService(repo Repository, executionRepo execution.Repository) Service {
	return Service{repo: repo, executionRepo: executionRepo}

}

var httpClient = http.Client{}

var wg sync.WaitGroup

func (s *Service) RunJobs(ctx context.Context) {

	fmt.Println("Starting")

	ticker := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			jobs, err := s.repo.FindAll(ctx)
			if err != nil {
				log.Fatal("Error getting jobs", err)
			}
			
			for _, job := range jobs {
				go func() {

					fmt.Println(job.Name)

					req, err := http.NewRequest(job.Method, job.Endpoint, bytes.NewReader([]byte(job.Body)))

					if err != nil {
						log.Println("Error creating request", err)
						return
					}

					resp, err := httpClient.Do(req)
					if err != nil {
						log.Println("Error sending request", err)
						return
					}
					resp.Body.Close()

					if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
						log.Printf("Error: received status code %d for job %s", resp.StatusCode, job.Name)
						return
					}

					sched, _ := cron.ParseStandard(job.Schedule)

					job.NextRunAt = sched.Next(time.Now())
					s.repo.Update(ctx, job)

				}()

			}
		}
	}
}

func (s *Service) Create(ctx context.Context, job Job) (Job, error) {
	return s.repo.Create(ctx, job)
}
