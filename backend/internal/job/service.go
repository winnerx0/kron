package job

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"github.com/robfig/cron/v3"
	"github.com/winnerx0/kron/internal/domain"
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/secret"
)

type Service interface {
	Create(ctx context.Context, userID string, job domain.Job) (domain.Job, error)
	Update(ctx context.Context, userID string, job domain.Job) (domain.Job, error)
	Delete(ctx context.Context, userID string, id string) error
	FindAll(ctx context.Context, userID string) ([]JobResponse, error)
	RunJob(ctx context.Context, userID string, id string) error
	StopJob(ctx context.Context, userID string, id string) (bool, error)
	UpdateJobStatus(ctx context.Context, userID string, jobID string) error
	RunJobs(ctx context.Context, jobsCh chan<- domain.Job)
}

type JobService struct {
	repo          Repository
	executionRepo execution.Repository
	client        http.Client
	secrets       secret.Manager
	activeMu      *sync.Mutex
	activeRuns    map[string]activeRun
}

type activeRun struct {
	executionID string
	cancel      context.CancelFunc
}

func NewJobService(repo Repository, executionRepo execution.Repository, managers ...secret.Manager) *JobService {
	var manager secret.Manager = secret.NoopManager{}
	if len(managers) > 0 && managers[0] != nil {
		manager = managers[0]
	}
	return &JobService{
		repo:          repo,
		executionRepo: executionRepo,
		client:        http.Client{},
		secrets:       manager,
		activeMu:      &sync.Mutex{},
		activeRuns:    map[string]activeRun{},
	}
}

func (s *JobService) RunJobs(ctx context.Context, jobsCh chan<- domain.Job) {

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
				jobsCh <- job
			}
		}
	}
}

func (s *JobService) Create(ctx context.Context, userID string, job domain.Job) (domain.Job, error) {
	job.UserID = userID

	if err := setNextRun(&job); err != nil {
		return domain.Job{}, err
	}

	job, err := s.encryptJobSecrets(job, nil)
	if err != nil {
		return domain.Job{}, err
	}

	createdJob, err := s.repo.Create(ctx, job)
	if err != nil {
		return domain.Job{}, err
	}

	return s.decryptJobSecrets(createdJob)
}

func (s *JobService) Update(ctx context.Context, userID string, job domain.Job) (domain.Job, error) {
	existingJob, err := s.repo.FindByID(ctx, job.ID)
	if err != nil {
		return domain.Job{}, err
	}

	if existingJob.UserID != userID {
		return domain.Job{}, ErrForbidden
	}

	job.UserID = existingJob.UserID

	if err := setNextRun(&job); err != nil {
		return domain.Job{}, err
	}

	job, err = s.encryptJobSecrets(job, existingJob.Headers)
	if err != nil {
		return domain.Job{}, err
	}

	updatedJob, err := s.repo.Update(ctx, job)
	if err != nil {
		return domain.Job{}, err
	}

	return s.decryptJobSecrets(updatedJob)
}

func (s *JobService) Delete(ctx context.Context, userID string, id string) error {
	existingJob, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if existingJob.UserID != userID {
		return ErrForbidden
	}

	activeRun, ok := s.activeRuns[existingJob.ID]

	if ok {
		activeRun.cancel()
		delete(s.activeRuns, existingJob.ID)
	}

	return s.repo.Delete(ctx, id)
}

func (s *JobService) FindAll(ctx context.Context, userID string) ([]JobResponse, error) {
	jobs, err := s.repo.FindAll(ctx, userID)
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
			Headers:     s.decryptHeadersForResponse(job.Headers),
			Body:        job.Body,
			Enabled:     job.Status,
		})
	}

	return jobResponses, nil
}

func (s *JobService) RunJob(ctx context.Context, userID string, id string) error {
	job, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if job.UserID != userID {
		return ErrForbidden
	}

	if job.Status == false {
		return ErrJobDisabled
	}
	return nil
}

func (s *JobService) StopJob(ctx context.Context, userID string, id string) (bool, error) {
	job, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return false, err
	}

	if job.UserID != userID {
		return false, ErrForbidden
	}

	s.ensureActiveExecutionTracking()

	s.activeMu.Lock()
	run, ok := s.activeRuns[id]
	s.activeMu.Unlock()

	if !ok {
		return false, nil
	}

	run.cancel()
	return true, nil
}

// readResponseBody reads up to 64KB of the response body so that execution
// details can show what the endpoint returned without storing unbounded data.
func readResponseBody(r io.Reader) string {
	const maxBodyBytes = 64 * 1024
	data, err := io.ReadAll(io.LimitReader(r, maxBodyBytes))
	if err != nil {
		return ""
	}
	return string(data)
}

func (s *JobService) advanceNextRun(ctx context.Context, job domain.Job) {
	if err := setNextRun(&job); err != nil {
		log.Println("Error parsing cron expression", err)
		return
	}

	job, err := s.encryptJobSecrets(job, nil)
	if err != nil {
		log.Println("Error encrypting job headers", err)
		return
	}

	if _, err := s.repo.Update(ctx, job); err != nil {
		log.Println("Error updating job", err)
	}
}

func setNextRun(job *domain.Job) error {
	sched, err := cron.ParseStandard(job.Schedule)
	if err != nil {
		return InvalidScheduleError{Schedule: job.Schedule, Err: err}
	}
	job.NextRunAt = sched.Next(time.Now())
	return nil
}

func (s *JobService) ensureActiveExecutionTracking() {
	if s.activeMu == nil {
		s.activeMu = &sync.Mutex{}
	}
	if s.activeRuns == nil {
		s.activeRuns = map[string]activeRun{}
	}
}

func (s *JobService) UpdateJobStatus(ctx context.Context, userID string, jobID string) error {
	job, err := s.repo.FindByID(ctx, jobID)
	if err != nil {
		return err
	}

	if job.UserID != userID {
		return ErrForbidden
	}

	job.Status = !job.Status

	if _, err := s.repo.Update(ctx, job); err != nil {
		return err
	}

	s.ensureActiveExecutionTracking()

	s.activeMu.Lock()
	activeRun, ok := s.activeRuns[job.ID]
	if ok {
		delete(s.activeRuns, jobID)
	}
	s.activeMu.Unlock()

	if ok {
		activeRun.cancel()
	}

	return nil
}

func (s *JobService) secretManager() secret.Manager {
	if s.secrets == nil {
		return secret.NoopManager{}
	}
	return s.secrets
}

func (s *JobService) encryptJobSecrets(job domain.Job, existingHeaders map[string]any) (domain.Job, error) {
	headers := make(map[string]any, len(job.Headers))
	for key, value := range job.Headers {
		valueString := fmt.Sprintf("%v", value)
		if !secret.IsSensitiveHeader(key) {
			headers[key] = value
			continue
		}

		if secret.IsMasked(valueString) && existingHeaders != nil {
			if existingValue, ok := existingHeaders[key]; ok {
				headers[key] = existingValue
				continue
			}
		}

		encryptedValue, err := s.secretManager().Encrypt(valueString)
		if err != nil {
			return domain.Job{}, err
		}
		headers[key] = encryptedValue
	}

	job.Headers = headers
	return job, nil
}

func (s *JobService) decryptHeaders(headers map[string]any) (map[string]string, error) {
	decrypted := make(map[string]string, len(headers))
	for key, value := range headers {
		valueString := fmt.Sprintf("%v", value)
		if !secret.IsSensitiveHeader(key) {
			decrypted[key] = valueString
			continue
		}

		decryptedValue, err := s.secretManager().Decrypt(valueString)
		if err != nil {
			return nil, err
		}
		decrypted[key] = decryptedValue
	}
	return decrypted, nil
}

func (s *JobService) decryptJobSecrets(job domain.Job) (domain.Job, error) {
	headers, err := s.decryptHeaders(job.Headers)
	if err != nil {
		return domain.Job{}, err
	}

	decryptedHeaders := make(map[string]any, len(headers))
	for key, value := range headers {
		decryptedHeaders[key] = value
	}

	job.Headers = decryptedHeaders
	return job, nil
}

func (s *JobService) decryptHeadersForResponse(headers map[string]any) map[string]any {
	decrypted, err := s.decryptHeaders(headers)
	if err != nil {
		return headers
	}

	responseHeaders := make(map[string]any, len(decrypted))
	for key, value := range decrypted {
		responseHeaders[key] = value
	}
	return responseHeaders
}

func exponentialBackoff(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	delay := baseDelay * time.Duration(1<<attempt)

	if delay > maxDelay {
		return maxDelay
	}

	return delay
}
