package job

import (
	"time"

	"github.com/winnerx0/kron/internal/execution"
	"gorm.io/datatypes"
)

type Job struct {
	ID string `gorm:"primaryKey,default:uuid_generate_v4()"`

	Name string `gorm:"not null"`

	Description string `gorm:""`

	Schedule string `gorm:"not null"`

	Endpoint string `gorm:"not null"`

	Method string `gorm:"not null"`

	Headers datatypes.JSONMap `gorm:"not null,type:jsonb"`

	Body string `gorm:""`

	NextRunAt time.Time `gorm:"not null"`

	Executions []execution.Execution `gorm:"foreignKey:JobID"`
}
