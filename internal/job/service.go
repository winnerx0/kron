package job

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/secret"
)

type Service struct {
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

func NewJobService(repo Repository, executionRepo execution.Repository, managers ...secret.Manager) Service {
	var manager secret.Manager = secret.NoopManager{}
	if len(managers) > 0 && managers[0] != nil {
		manager = managers[0]
	}
	return Service{
		repo:          repo,
		executionRepo: executionRepo,
		client:        http.Client{},
		secrets:       manager,
		activeMu:      &sync.Mutex{},
		activeRuns:    map[string]activeRun{},
	}
}

func (s *Service) RunJobs(ctx context.Context, jobsCh chan<- Job) {

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
				log.Printf("%+v\n", job)
				jobsCh <- job
			}
		}
	}
}

func (s *Service) Create(ctx context.Context, userID string, job Job) (Job, error) {
	job.UserID = userID

	if err := setNextRun(&job); err != nil {
		return Job{}, err
	}

	job, err := s.encryptJobSecrets(job, nil)
	if err != nil {
		return Job{}, err
	}

	createdJob, err := s.repo.Create(ctx, job)
	if err != nil {
		return Job{}, err
	}

	return s.decryptJobSecrets(createdJob)
}

func (s *Service) Update(ctx context.Context, userID string, job Job) (Job, error) {
	existingJob, err := s.repo.FindByID(ctx, job.ID)
	if err != nil {
		return Job{}, err
	}

	if existingJob.UserID != userID {
		return Job{}, ErrForbidden
	}

	job.UserID = existingJob.UserID

	if err := setNextRun(&job); err != nil {
		return Job{}, err
	}

	job, err = s.encryptJobSecrets(job, existingJob.Headers)
	if err != nil {
		return Job{}, err
	}

	updatedJob, err := s.repo.Update(ctx, job)
	if err != nil {
		return Job{}, err
	}

	return s.decryptJobSecrets(updatedJob)
}

func (s *Service) Delete(ctx context.Context, userID string, id string) error {
	existingJob, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if existingJob.UserID != userID {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) FindAll(ctx context.Context, userID string) ([]JobResponse, error) {
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

func (s *Service) RunJob(ctx context.Context, userID string, id string) error {
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

	go s.ExecuteJob(context.Background(), job, false)
	return nil
}

func (s *Service) StopJob(ctx context.Context, userID string, id string) (bool, error) {
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

func (s *Service) ExecuteJob(ctx context.Context, job Job, advanceSchedule bool) {

	s.ensureActiveExecutionTracking()

	executionCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	newExecution := execution.Execution{
		ID:      uuid.NewString(),
		JobID:   job.ID,
		Status:  execution.RUNNING,
		Started: time.Now(),
	}

	s.activeMu.Lock()
	if existingRun, ok := s.activeRuns[job.ID]; ok {
		log.Println("Cancelling active run")
		existingRun.cancel()
	}
	s.activeRuns[job.ID] = activeRun{executionID: newExecution.ID, cancel: cancel}
	s.activeMu.Unlock()

	// remove run from active runs when the current run is active run
	defer func() {
		s.activeMu.Lock()
		if run, ok := s.activeRuns[job.ID]; ok && run.executionID == newExecution.ID {
			delete(s.activeRuns, job.ID)
		}
		s.activeMu.Unlock()
	}()

	err := s.executionRepo.Save(executionCtx, newExecution)
	if err != nil {
		log.Println("Error saving execution", err)
		return
	}

	finish := func(status execution.ExecutionStatus) {
		newExecution.Finished = time.Now()
		newExecution.Status = status
		if err := s.executionRepo.Update(context.Background(), newExecution); err != nil {
			log.Println("Error updating execution", err)
		}

		if advanceSchedule {
			s.advanceNextRun(context.Background(), job)
		}
	}

	select {
	case <-executionCtx.Done():
		finish(execution.STOPPED)
		return
	default:
	}

	for attempt := range 5 {
		req, err := http.NewRequestWithContext(executionCtx, job.Method, job.Endpoint, bytes.NewReader([]byte(job.Body)))
		if err != nil {
			log.Println("Error creating request", err)
			finish(execution.FAILED)
			return
		}

		headers, err := s.decryptHeaders(job.Headers)
		if err != nil {
			log.Println("Error decrypting job headers", err)
			finish(execution.FAILED)
			return
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := s.client.Do(req)
		if err != nil {
			log.Println("Error sending request", err)
			if errors.Is(executionCtx.Err(), context.Canceled) || errors.Is(executionCtx.Err(), context.DeadlineExceeded) {
				finish(execution.STOPPED)
				return
			}
			if attempt < 4 {
				select {
				case <-executionCtx.Done():
					finish(execution.STOPPED)
					return
				case <-time.After(exponentialBackoff(attempt, time.Second, 30*time.Second)):
				}
				continue
			}
			finish(execution.FAILED)
			return
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			log.Printf("Error: received status code %d for job %s", resp.StatusCode, job.Name)
			resp.Body.Close()
			if resp.StatusCode >= 500 && attempt < 4 {
				select {
				case <-executionCtx.Done():
					finish(execution.STOPPED)
					return
				case <-time.After(exponentialBackoff(attempt, time.Second, 30*time.Second)):
				}
				continue
			}
			finish(execution.FAILED)
			return
		}

		resp.Body.Close()
		finish(execution.SUCCESS)
		return
	}
}

func (s *Service) advanceNextRun(ctx context.Context, job Job) {
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

func setNextRun(job *Job) error {
	sched, err := cron.ParseStandard(job.Schedule)
	if err != nil {
		return InvalidScheduleError{Schedule: job.Schedule, Err: err}
	}
	job.NextRunAt = sched.Next(time.Now())
	return nil
}

func (s *Service) ensureActiveExecutionTracking() {
	if s.activeMu == nil {
		s.activeMu = &sync.Mutex{}
	}
	if s.activeRuns == nil {
		s.activeRuns = map[string]activeRun{}
	}
}

func (s *Service) UpdateJobStatus(ctx context.Context, userID string, jobID string) error {
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

func (s *Service) secretManager() secret.Manager {
	if s.secrets == nil {
		return secret.NoopManager{}
	}
	return s.secrets
}

func (s *Service) encryptJobSecrets(job Job, existingHeaders map[string]any) (Job, error) {
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
			return Job{}, err
		}
		headers[key] = encryptedValue
	}

	job.Headers = headers
	return job, nil
}

func (s *Service) decryptHeaders(headers map[string]any) (map[string]string, error) {
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

func (s *Service) decryptJobSecrets(job Job) (Job, error) {
	headers, err := s.decryptHeaders(job.Headers)
	if err != nil {
		return Job{}, err
	}

	decryptedHeaders := make(map[string]any, len(headers))
	for key, value := range headers {
		decryptedHeaders[key] = value
	}

	job.Headers = decryptedHeaders
	return job, nil
}

func (s *Service) decryptHeadersForResponse(headers map[string]any) map[string]any {
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
