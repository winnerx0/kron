package user

import (
	"time"

	refreshtoken "github.com/winnerx0/kron/internal/refresh_token"
)

type User struct {
	ID             string                   `gorm:"type:uuid;primaryKey,default:gen_random_uuidv7()"`
	Username       string                      `gorm:"type:varchar(50);unique;not null"`
	Email          string                      `gorm:"type:varchar(255);unique;not null"`
	ProfilePicture string                      `gorm:"type:varchar(255);not null"`
	JoinedAt       time.Time                   `gorm:"joined_at;default:now();not null"`
	RefreshTokens  []refreshtoken.RefreshToken `gorm:"foreignKey:UserID"`
}
