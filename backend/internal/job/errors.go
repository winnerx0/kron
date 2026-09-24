package job

import "fmt"

type ForbiddenError struct{}

func (ForbiddenError) Error() string {
	return "forbidden"
}

type JobDisabledError struct{}

func (JobDisabledError) Error() string {
	return "job is disabled"
}

type InvalidScheduleError struct {
	Schedule string
	Err      error
}

func (e InvalidScheduleError) Error() string {
	return fmt.Sprintf("invalid cron expression: %s", e.Schedule)
}

func (e InvalidScheduleError) Unwrap() error {
	return e.Err
}

var (
	ErrForbidden   = ForbiddenError{}
	ErrJobDisabled = JobDisabledError{}
)
