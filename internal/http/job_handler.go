package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/winnerx0/kron/internal/job"
	"github.com/winnerx0/kron/internal/validator"
)

type JobHandler struct {
	service job.Service
}

var validate = validator.Get()

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
// @Router /job/create [post]
func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createJobRequest job.CreateJobRequest

	if err := json.NewDecoder(r.Body).Decode(&createJobRequest); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: err.Error()})
		return
	}

	if err := validate.Struct(createJobRequest); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: validator.FirstError(err)})
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

// @Summary Delete a job
// @Description Delete a job by its ID
// @Param jobID path string true "Job ID"
// @Success 204
// @Failure 400 {object} job.ErrorResponse
// @Failure 500 {object} job.ErrorResponse
// @Router /jobs/{jobID} [delete]
func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobID")

	if jobID == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: "Missing job ID"})
		return
	}

	if err := h.service.Delete(r.Context(), jobID); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(204)
}

// @Summary Run a job
// @Description Immediately run a job by its ID
// @Param jobID path string true "Job ID"
// @Success 202
// @Failure 400 {object} job.ErrorResponse
// @Failure 500 {object} job.ErrorResponse
// @Router /jobs/{jobID}/run [post]
func (h *JobHandler) RunJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobID")

	if jobID == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: "Missing job ID"})
		return
	}

	if err := h.service.RunJob(r.Context(), jobID); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// @Summary Stop a job
// @Description Stop a running job by its ID
// @Param jobID path string true "Job ID"
// @Success 202
// @Failure 400 {object} job.ErrorResponse
// @Failure 404 {object} job.ErrorResponse
// @Router /jobs/{jobID}/stop [post]
func (h *JobHandler) StopJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobID")

	if jobID == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: "Missing job ID"})
		return
	}

	if !h.service.StopJob(jobID) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: "Job is not running"})
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// @Summary Update a job
// @Description Update a job by its ID
// @Param jobID path string true "Job ID"
// @Produce json
// @Param job body job.UpdateJobRequest true "Job"
// @Success 200 {object} job.UpdateJobResponse
// @Failure 400 {object} job.ErrorResponse
// @Failure 500 {object} job.ErrorResponse
// @Router /jobs/{jobID} [put]
func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobID")

	if jobID == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: "Missing job ID"})
		return
	}

	var updateJobRequest job.UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&updateJobRequest); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: err.Error()})
		return
	}

	updatedJob, err := h.service.Update(r.Context(), job.Job{
		ID:          jobID,
		Name:        updateJobRequest.Name,
		Description: updateJobRequest.Description,
		Schedule:    updateJobRequest.Schedule,
		Endpoint:    updateJobRequest.Endpoint,
		Method:      updateJobRequest.Method,
		Headers:     updateJobRequest.Headers,
		Body:        updateJobRequest.Body,
	})

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(job.UpdateJobResponse{
		ID:          updatedJob.ID,
		Name:        updatedJob.Name,
		Description: updatedJob.Description,
		Schedule:    updatedJob.Schedule,
		Endpoint:    updatedJob.Endpoint,
		Method:      updatedJob.Method,
		Headers:     updatedJob.Headers,
		Body:        updatedJob.Body,
	})
}

// @Summary Find all jobs
// @Description Find all created jobs
// @Produce json
// @Success 200 {array} job.JobResponse
// @Failure 500 {object} job.ErrorResponse
// @Router /job/all [get]
func (h *JobHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.FindAll(r.Context())
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(job.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(jobs)
}
