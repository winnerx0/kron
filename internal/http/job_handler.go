package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/winnerx0/kron/internal/job"
)

type JobHandler struct {
	service job.Service
}

func NewJobHandler(service job.Service) *JobHandler {
	return &JobHandler{service: service}
}

// @Summary Create a job
// @Description Creates a new job.
// @Accept json
// @Produce json
// @Param job body job.CreateJobRequest true "Job"
// @Success 200 {object} job.CreateJobResponse "Success"
// @Failure 400 {object} job.ErrorResponse "Invalid request"
// @Failure 500 {object} job.ErrorResponse "Internal server error"
// @Router /create [post]
func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createJobRequest job.CreateJobRequest

	if err := json.NewDecoder(r.Body).Decode(&createJobRequest); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: err.Error()})
		return
	}

	_, err := cron.ParseStandard(createJobRequest.Schedule)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: "Invalid cron expression"})
		return
	}

	newJob, err := h.service.Create(r.Context(), job.Job{
		ID:          uuid.NewString(),
		Name:        createJobRequest.Name,
		Description: createJobRequest.Description,
		Schedule:    createJobRequest.Schedule,
		Endpoint:    createJobRequest.Endpoint,
		Method:      createJobRequest.Method,
		Headers:     createJobRequest.Headers,
		Body:        createJobRequest.Body,
	})

	if err != nil {
		log.Println("Error parsing cron expression", err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(200)

	json.NewEncoder(w).Encode(job.CreateJobResponse{
		ID:          newJob.ID,
		Name:        newJob.Name,
		Description: newJob.Description,
		Schedule:    newJob.Schedule,
		Endpoint:    newJob.Endpoint,
		Method:      newJob.Method,
		Headers:     newJob.Headers,
		Body:        newJob.Body,
	})
}
