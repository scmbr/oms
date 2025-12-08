package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"not null;default:user"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
