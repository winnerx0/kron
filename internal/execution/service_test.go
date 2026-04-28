package execution

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestExecutionService_Create_Success(t *testing.T){

	ctx := context.Background()

	execution := Execution{
		ID:     uuid.NewString(),
		JobID:  uuid.NewString(),
		Status: PENDING,
		Started: time.Now(),
		Finished: time.Now().Add(1 * time.Minute),
	}

	mockRepo := new(MockRepository)

	mockRepo.On("Save", ctx, execution).Return(nil)

	service := Service{
		repo: mockRepo,
	}

	err := service.Create(ctx, execution)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mockRepo.AssertExpectations(t)
}