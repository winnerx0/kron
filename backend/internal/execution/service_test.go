package execution

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/winnerx0/kron/internal/domain"
)

func TestExecutionService_Create_Success(t *testing.T){

	ctx := context.Background()

	execution := domain.Execution{
		ID:     uuid.NewString(),
		JobID:  uuid.NewString(),
		Status: domain.RUNNING,
		Started: time.Now(),
		Finished: time.Now().Add(1 * time.Minute),
	}

	mockRepo := new(MockRepository)

	mockRepo.On("Save", ctx, execution).Return(nil)

	service := NewExecutionService(mockRepo)

	err := service.Create(ctx, execution)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mockRepo.AssertExpectations(t)
}