package models

import "time"

type OutboxEvent struct {
	ID          uint      `gorm:"primaryKey"`
	ExternalID  string    `gorm:"type:uuid;uniqueIndex"`
	EventType   string    `gorm:"type:varchar(64);not null"`
	OrderID     string    `gorm:"type:uuid;not null"`
	Payload     []byte    `gorm:"type:jsonb;not null"`
	Status      string    `gorm:"type:varchar(16);not null;default:'pending'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	ProcessedAt *time.Time
}
