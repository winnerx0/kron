package job

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/mock"
	"github.com/winnerx0/kron/internal/execution"
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

	mockRepo.On("Create", ctx, job).Return(job, nil)

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

	mockRepo := new(MockRepository)

	sched, _ := cron.ParseStandard(job.Schedule)

	mockRepo.On("Update", ctx, mock.MatchedBy(func(j Job) bool {

		return j.ID == job.ID && time.Time.Equal(j.NextRunAt, sched.Next(time.Now()))
	})).Return(nil)

	executionRepo := new(execution.MockRepository)

	service := Service{
		repo:          mockRepo,
		executionRepo: executionRepo,
		client:        *mockServer.Client(),
	}

	service.ExecuteJob(ctx, job)
}
