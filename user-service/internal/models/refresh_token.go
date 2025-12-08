package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;column:token_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index;column:user_id"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
