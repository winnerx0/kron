package execution_test

import (
	"context"
	"testing"

	"github.com/winnerx0/kron/internal/mocks"
	"time"

	"github.com/google/uuid"
	"github.com/winnerx0/kron/internal/domain"
	"github.com/winnerx0/kron/internal/execution"
)

func TestExecutionService_Create_Success(t *testing.T) {

	ctx := context.Background()

	executionRecord := domain.Execution{
		ID:       uuid.NewString(),
		JobID:    uuid.NewString(),
		Status:   domain.RUNNING,
		Started:  time.Now(),
		Finished: time.Now().Add(1 * time.Minute),
	}

	mockRepo := new(mocks.MockExecutionRepository)

	mockRepo.On("Save", ctx, executionRecord).Return(nil)

	service := execution.NewExecutionService(mockRepo)

	err := service.Create(ctx, executionRecord)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mockRepo.AssertExpectations(t)
}
