package execution

import (
	"time"
)

type Execution struct {
	ID       string          `gorm:"type:uuid;primaryKey,default:uuid_generate_v7()" json:"id"`
	JobID    string          `gorm:"type:uuid;not null;foreignKey:JobID;references:ID" json:"jobID"`
	Status   ExecutionStatus `gorm:"type:varchar(20);not null" json:"status"`
	Started  time.Time       `gorm:"not null" json:"startedAt"`
	Finished time.Time       `gorm:"not null" json:"finishedAt"`
}

type ExecutionStatus string

const (
	RUNNING ExecutionStatus = "running"
	SUCCESS ExecutionStatus = "success"
	FAILED  ExecutionStatus = "failed"
	STOPPED ExecutionStatus = "stopped"
)
