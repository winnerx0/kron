package job

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/winnerx0/kron/internal/execution"
)

type Service struct {
	repo          Repository
	executionRepo execution.Repository
	client        http.Client
}

func NewJobService(repo Repository, executionRepo execution.Repository) Service {
	return Service{repo: repo, executionRepo: executionRepo, client: http.Client{}}
}

func (s *Service) RunJobs(ctx context.Context) {

	fmt.Println("Starting")

	ticker := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			jobs, err := s.repo.FindAllNextRun(ctx)
			if err != nil {
				log.Fatal("Error getting jobs", err)
			}

			for _, job := range jobs {
				go s.ExecuteJob(ctx, job)

			}
		}
	}
}

func (s *Service) Create(ctx context.Context, job Job) (Job, error) {
	return s.repo.Create(ctx, job)
}

func (s *Service) Update(ctx context.Context, job Job) (Job, error) {
	return s.repo.Update(ctx, job)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]JobResponse, error) {
	jobs, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var jobResponses []JobResponse
	for _, job := range jobs {
		jobResponses = append(jobResponses, JobResponse{
			ID:          job.ID,
			Name:        job.Name,
			Description: job.Description,
			Schedule:    job.Schedule,
			Endpoint:    job.Endpoint,
			Method:      job.Method,
			Headers:     job.Headers,
			Body:        job.Body,
		})
	}

	return jobResponses, nil
}

func (s *Service) ExecuteJob(ctx context.Context, job Job) {

	newExecution := execution.Execution{
		ID:      uuid.NewString(),
		JobID:   job.ID,
		Status:  execution.PENDING,
		Started: time.Now(),
	}

	err := s.executionRepo.Save(ctx, newExecution)

	if err != nil {
		log.Println("Error saving execution", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, job.Method, job.Endpoint, bytes.NewReader([]byte(job.Body)))

	if err != nil {
		log.Println("Error creating request", err)
		return
	}

	for key, value := range job.Headers {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}

	resp, err := s.client.Do(req)

	if err != nil {
		log.Println("Error sending request", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Error: received status code %d for job %s", resp.StatusCode, job.Name)
		return
	}

	sched, _ := cron.ParseStandard(job.Schedule)

	job.NextRunAt = sched.Next(time.Now())

	_, err = s.repo.Update(ctx, job)

	if err != nil {
		log.Println("Error updating job", err)
		return
	}

	newExecution.Finished = time.Now()
	newExecution.Status = execution.SUCCESS

	err = s.executionRepo.Update(ctx, newExecution)

	if err != nil {
		log.Println("Error saving execution", err)
		return
	}

}
