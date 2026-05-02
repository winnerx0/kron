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
	"github.com/winnerx0/kron/internal/secret"
)

type Service struct {
	repo          Repository
	executionRepo execution.Repository
	client        http.Client
	secrets       secret.Manager
}

func NewJobService(repo Repository, executionRepo execution.Repository, managers ...secret.Manager) Service {
	var manager secret.Manager = secret.NoopManager{}
	if len(managers) > 0 && managers[0] != nil {
		manager = managers[0]
	}
	return Service{repo: repo, executionRepo: executionRepo, client: http.Client{}, secrets: manager}
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

	headers, err := s.decryptHeaders(job.Headers)
	if err != nil {
		log.Println("Error decrypting job headers", err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
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

	job, err = s.encryptJobSecrets(job, nil)
	if err != nil {
		log.Println("Error encrypting job headers", err)
		return
	}

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
