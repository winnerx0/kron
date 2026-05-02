package job

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/secret"
)

func TestUserService_Create_Success(t *testing.T) {

	ctx := context.Background()

	job := Job{
		ID:     uuid.NewString(),
		Name:   "Job 1",
		Method: "GET",
		Headers: map[string]any{
			"Content-Type": "application/json",
		},
		Endpoint: "https://example.com",
		Body:     "",
		Schedule: "*/5 * * * *",
	}

	mockRepo := new(MockRepository)

	mockRepo.On("Create", ctx, mock.MatchedBy(func(j Job) bool {
		return j.ID == job.ID && j.NextRunAt.After(time.Now())
	})).Return(job, nil)

	executionRepo := new(execution.MockRepository)

	service := Service{
		repo:          mockRepo,
		executionRepo: executionRepo,
	}

	result, err := service.Create(ctx, job)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != job.ID {
		t.Errorf("got ID %s, want %s", result.ID, job.ID)
	}
	if result.Name != job.Name {
		t.Errorf("got Name %s, want %s", result.Name, job.Name)
	}

	mockRepo.AssertExpectations(t)

}

func TestUserService_ExecuteJob_Success(t *testing.T) {

	ctx := context.Background()

	job := Job{
		ID:     uuid.NewString(),
		Name:   "Job 1",
		Method: "GET",
		Headers: map[string]any{
			"Content-Type": "application/json",
		},
		Endpoint: "https://example.com",
		Body:     "",
		Schedule: "*/5 * * * *",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("got Content-Type %s, want application/json", got)
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(CreateJobResponse{
			ID:       job.ID,
			Name:     job.Name,
			Method:   job.Method,
			Endpoint: job.Endpoint,
			Body:     job.Body,
			Schedule: job.Schedule,
		})
	}))

	defer mockServer.Close()
	job.Endpoint = mockServer.URL

	mockRepo := new(MockRepository)

	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(j Job) bool {
		return j.ID == job.ID && j.NextRunAt.After(time.Now())
	})).Return(job, nil)

	executionRepo := new(execution.MockRepository)

	executionRepo.On("Save", mock.Anything, mock.MatchedBy(func(e execution.Execution) bool {
		return e.JobID == job.ID && e.Status == execution.RUNNING
	})).Return(nil)

	executionRepo.On("Update", mock.Anything, mock.MatchedBy(func(e execution.Execution) bool {
		return e.JobID == job.ID && e.Status == execution.SUCCESS
	})).Return(nil)

	service := Service{
		repo:          mockRepo,
		executionRepo: executionRepo,
		client:        *mockServer.Client(),
	}

	service.ExecuteJob(ctx, job, true)

	mockRepo.AssertExpectations(t)
	executionRepo.AssertExpectations(t)
}

func TestJobService_Create_EncryptsAndReturnsSensitiveHeaders(t *testing.T) {
	ctx := context.Background()
	manager, err := secret.NewAESGCMManager("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("unexpected manager error: %v", err)
	}

	job := Job{
		ID:       uuid.NewString(),
		Name:     "Job 1",
		Method:   "GET",
		Endpoint: "https://example.com",
		Schedule: "*/5 * * * *",
		Headers: map[string]any{
			"Authorization": "Bearer raw-token",
			"Content-Type":  "application/json",
		},
	}
	storedAuthorization, err := manager.Encrypt("Bearer raw-token")
	if err != nil {
		t.Fatalf("unexpected encryption error: %v", err)
	}

	mockRepo := new(MockRepository)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(j Job) bool {
		authorization, ok := j.Headers["Authorization"].(string)
		if !ok || !strings.HasPrefix(authorization, "kron:v1:") {
			return false
		}

		decrypted, err := manager.Decrypt(authorization)
		return err == nil &&
			decrypted == "Bearer raw-token" &&
			j.Headers["Content-Type"] == "application/json" &&
			j.NextRunAt.After(time.Now())
	})).Return(Job{
		ID:       job.ID,
		Name:     job.Name,
		Method:   job.Method,
		Endpoint: job.Endpoint,
		Schedule: job.Schedule,
		Headers: map[string]any{
			"Authorization": storedAuthorization,
			"Content-Type":  "application/json",
		},
	}, nil)

	service := NewJobService(mockRepo, new(execution.MockRepository), manager)

	result, err := service.Create(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Headers["Authorization"] != "Bearer raw-token" {
		t.Errorf("got Authorization %v, want raw value", result.Headers["Authorization"])
	}
	if result.Headers["Content-Type"] != "application/json" {
		t.Errorf("got Content-Type %v, want application/json", result.Headers["Content-Type"])
	}

	mockRepo.AssertExpectations(t)
}

func TestJobService_Update_PreservesMaskedSensitiveHeaders(t *testing.T) {
	ctx := context.Background()
	manager, err := secret.NewAESGCMManager("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("unexpected manager error: %v", err)
	}

	encrypted, err := manager.Encrypt("Bearer existing-token")
	if err != nil {
		t.Fatalf("unexpected encryption error: %v", err)
	}

	jobID := uuid.NewString()
	existingJob := Job{
		ID: jobID,
		Headers: map[string]any{
			"Authorization": encrypted,
		},
	}
	updateJob := Job{
		ID:       jobID,
		Name:     "Updated",
		Method:   "GET",
		Endpoint: "https://example.com",
		Schedule: "*/5 * * * *",
		Headers: map[string]any{
			"Authorization": secret.MaskedValue,
		},
	}

	mockRepo := new(MockRepository)
	mockRepo.On("FindByID", ctx, jobID).Return(existingJob, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(j Job) bool {
		return j.Headers["Authorization"] == encrypted
	})).Return(Job{
		ID:       updateJob.ID,
		Name:     updateJob.Name,
		Method:   updateJob.Method,
		Endpoint: updateJob.Endpoint,
		Schedule: updateJob.Schedule,
		Headers: map[string]any{
			"Authorization": encrypted,
		},
	}, nil)

	service := NewJobService(mockRepo, new(execution.MockRepository), manager)

	result, err := service.Update(ctx, updateJob)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Headers["Authorization"] != "Bearer existing-token" {
		t.Errorf("got Authorization %v, want raw value", result.Headers["Authorization"])
	}

	mockRepo.AssertExpectations(t)
}
