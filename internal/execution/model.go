package execution

import (
	"time"
)

type Execution struct {
	ID       string          `gorm:"primaryKey,default:uuid_generate_v4()" json:"id"`
	JobID    string          `gorm:"job_id" json:"jobID"`
	Status   ExecutionStatus `gorm:"type:varchar(20);not null" json:"status"`
	Started  time.Time       `gorm:"not null" json:"startedAt"`
	Finished time.Time       `gorm:"not null" json:"finishedAt"`
}

type ExecutionStatus string

const (
	PENDING ExecutionStatus = "pending"
	SUCCESS ExecutionStatus = "success"
	FAILED  ExecutionStatus = "failed"
)