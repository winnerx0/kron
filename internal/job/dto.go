package job

type CreateJobRequest struct {
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
