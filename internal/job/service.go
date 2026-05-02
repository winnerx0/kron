package job

import (
	"bytes"
	"context"
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
				go s.ExecuteJob(ctx, job, true)
			}
		}
	}
}

func (s *Service) Create(ctx context.Context, job Job) (Job, error) {
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

func (s *Service) Update(ctx context.Context, job Job) (Job, error) {
	existingJob, err := s.repo.FindByID(ctx, job.ID)
	if err != nil {
		return Job{}, err
	}

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
			Headers:     s.decryptHeadersForResponse(job.Headers),
			Body:        job.Body,
		})
	}

	return jobResponses, nil
}

func (s *Service) RunJob(ctx context.Context, id string) error {
	job, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	go s.ExecuteJob(ctx, job, false)
	return nil
}

func (s *Service) StopJob(id string) bool {
	s.ensureActiveExecutionTracking()

	s.activeMu.Lock()
	run, ok := s.activeRuns[id]
	s.activeMu.Unlock()

	if !ok {
		return false
	}

	run.cancel()
	return true
}

func (s *Service) ExecuteJob(ctx context.Context, job Job, advanceSchedule bool) {
	s.ensureActiveExecutionTracking()

	executionCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if advanceSchedule {
		s.advanceNextRun(context.Background(), job)
	}

	newExecution := execution.Execution{
		ID:      uuid.NewString(),
		JobID:   job.ID,
		Status:  execution.RUNNING,
		Started: time.Now(),
	}

	s.activeMu.Lock()
	if existingRun, ok := s.activeRuns[job.ID]; ok {
		existingRun.cancel()
	}
	s.activeRuns[job.ID] = activeRun{executionID: newExecution.ID, cancel: cancel}
	s.activeMu.Unlock()

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
	}

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
		if executionCtx.Err() != nil {
			finish(execution.STOPPED)
			return
		}
		finish(execution.FAILED)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Error: received status code %d for job %s", resp.StatusCode, job.Name)
		finish(execution.FAILED)
		return
	}

	finish(execution.SUCCESS)
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
		return err
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
