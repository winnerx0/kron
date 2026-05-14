package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/winnerx0/kron/internal/job"
	"github.com/winnerx0/kron/internal/response"
	"github.com/winnerx0/kron/internal/validator"
)

type JobHandler struct {
	service job.Service
}

var validate = validator.Get()

func NewJobHandler(service job.Service) *JobHandler {
	return &JobHandler{service: service}
}

func userIDFromContext(r *http.Request) (string, bool) {
	value := r.Context().Value("userId")
	id, ok := value.(string)
	return id, ok
}

func writeJobError(w http.ResponseWriter, err error) {
	if errors.Is(err, job.ErrForbidden) {
		response.WriteError(w, http.StatusForbidden, "Forbidden")
		return
	}
	response.WriteError(w, http.StatusInternalServerError, err.Error())
}

// @Summary Create a job
// @Description Creates a new job.
// @Accept json
// @Produce json
// @Param job body job.CreateJobRequest true "Job"
// @Success 200 {object} job.CreateJobResponse "Success"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/job/create [post]
func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var createJobRequest job.CreateJobRequest

	if err := json.NewDecoder(r.Body).Decode(&createJobRequest); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(createJobRequest); err != nil {
		response.WriteError(w, http.StatusBadRequest, validator.FirstError(err))
		return
	}

	_, err := cron.ParseStandard(createJobRequest.Schedule)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid cron expression")
		return
	}

	newJob, err := h.service.Create(r.Context(), userID, job.Job{
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
		log.Println("Error creating job", err)
		writeJobError(w, err)
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
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/job/{jobID} [delete]
func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jobID := chi.URLParam(r, "jobID")

	if jobID == "" {
		response.WriteError(w, http.StatusBadRequest, "Missing job ID")
		return
	}

	if err := h.service.Delete(r.Context(), userID, jobID); err != nil {
		writeJobError(w, err)
		return
	}

	w.WriteHeader(204)
}

// @Summary Run a job
// @Description Immediately run a job by its ID
// @Param jobID path string true "Job ID"
// @Success 202
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/job/{jobID}/run [post]
func (h *JobHandler) RunJob(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jobID := chi.URLParam(r, "jobID")

	if jobID == "" {
		response.WriteError(w, http.StatusBadRequest, "Missing job ID")
		return
	}

	if err := h.service.RunJob(r.Context(), userID, jobID); err != nil {
		writeJobError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// @Summary Stop a job
// @Description Stop a running job by its ID
// @Param jobID path string true "Job ID"
// @Success 202
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/job/{jobID}/stop [post]
func (h *JobHandler) StopJob(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jobID := chi.URLParam(r, "jobID")

	if jobID == "" {
		response.WriteError(w, http.StatusBadRequest, "Missing job ID")
		return
	}

	stopped, err := h.service.StopJob(r.Context(), userID, jobID)
	if err != nil {
		writeJobError(w, err)
		return
	}

	if !stopped {
		response.WriteError(w, http.StatusNotFound, "Job is not running")
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
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/job/{jobID} [put]
func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jobID := chi.URLParam(r, "jobID")

	if jobID == "" {
		response.WriteError(w, http.StatusBadRequest, "Missing job ID")
		return
	}

	var updateJobRequest job.UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&updateJobRequest); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedJob, err := h.service.Update(r.Context(), userID, job.Job{
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
		writeJobError(w, err)
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
// @Failure 500 {object} response.ErrorResponse
// @Router /api/job/all [get]
func (h *JobHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jobs, err := h.service.FindAll(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(jobs)
}

// @Summary Update job status
// @Description Enable or disable a job by its ID
// @Param jobID path string true "Job ID"
// @Success 202
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/job/{jobID}/status [patch]
func (h *JobHandler) UpdateJobStatus(w http.ResponseWriter, r *http.Request) {

	userID, ok := userIDFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jobID := chi.URLParam(r, "jobID")

	if jobID == "" {
		response.WriteError(w, http.StatusBadRequest, "Missing job ID")
		return
	}

	if err := h.service.UpdateJobStatus(r.Context(), userID, jobID); err != nil {
		writeJobError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
