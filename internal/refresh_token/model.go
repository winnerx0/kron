package refreshtoken

import (
	"time"
)

type RefreshToken struct {
	ID        string    `gorm:"type:uuid;primaryKey,default:gen_random_uuidv7()"`
	UserID    string    `gorm:"type:uuid;not null;"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:now();not null"`
}
