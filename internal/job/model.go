package job

import (
	"time"

	"github.com/winnerx0/kron/internal/execution"
	"gorm.io/datatypes"
)

type Job struct {
	ID string `gorm:"type:uuid;primaryKey,default:uuid_generate_v7()"`

	Name string `gorm:"type:varchar(50);not null"`

	Description string `gorm:"type:varchar(255)"`

	Schedule string `gorm:"type:varchar(20);not null"`

	Endpoint string `gorm:"not null"`

	Method string `gorm:"type:varchar(8);not null"`

	Headers datatypes.JSONMap `gorm:"not null;type:jsonb"`

	Body string `gorm:""`

	NextRunAt time.Time `gorm:"not null"`

	Status bool `gorm:"not null;default:true"`

	Executions []execution.Execution `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`

	UserID string `gorm:"type:uuid;user_id"`
}
