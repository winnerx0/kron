package execution

import (
	"time"
)

type Execution struct {
	ID       int             `gorm:"primaryKey"`
	JobID    string          `gorm:"job_id"`
	Status   ExecutionStatus `gorm:"type:varchar(20);not null"`
	Started  time.Time       `gorm:"not null"`
	Finished time.Time       `gorm:"not null"`
}

type ExecutionStatus int

const (
	PENDING = iota
	SUCCESS 
	FAILED
)

var StatusName = map[ExecutionStatus]string{
	SUCCESS: "Success",
	PENDING: "Pending",
	FAILED:  "Failed",
}

func (s ExecutionStatus) String() string {
	return StatusName[s]
}
