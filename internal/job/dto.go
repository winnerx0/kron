package job

type CreateJobRequest struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description"`
	Schedule    string         `json:"schedule" validate:"required"`
	Endpoint    string         `json:"endpoint" validate:"required,url"`
	Method      string         `json:"method" validate:"required,oneof=GET POST PUT DELETE PATCH"`
	Headers     map[string]any `json:"headers" validate:"required"`
	Body        string         `json:"body"`
}

type UpdateJobRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schedule    string         `json:"schedule"`
	Endpoint    string         `json:"endpoint"`
	Method      string         `json:"method"`
	Headers     map[string]any `json:"headers"`
	Body        string         `json:"body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateJobResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schedule    string         `json:"schedule"`
	Endpoint    string         `json:"endpoint"`
	Method      string         `json:"method"`
	Headers     map[string]any `json:"headers"`
	Body        string         `json:"body"`
}

type UpdateJobResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schedule    string         `json:"schedule"`
	Endpoint    string         `json:"endpoint"`
	Method      string         `json:"method"`
	Headers     map[string]any `json:"headers"`
	Body        string         `json:"body"`
}

type JobResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schedule    string         `json:"schedule"`
	Endpoint    string         `json:"endpoint"`
	Method      string         `json:"method"`
	Headers     map[string]any `json:"headers"`
	Body        string         `json:"body"`
}