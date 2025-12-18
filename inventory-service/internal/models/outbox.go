package models

import "time"

type OutboxStatus string

var (
	OutboxPending OutboxStatus = "PENDING"
	OutboxSent    OutboxStatus = "SENT"
	OutboxFailed  OutboxStatus = "FAILED"
)

type OutboxEvent struct {
	ID            uint         `gorm:"primaryKey"`
	ExternalID    string       `gorm:"type:uuid;uniqueIndex"`
	EventType     string       `gorm:"type:varchar(64);not null"`
	AggregateID   string       `gorm:"type:uuid;not null"`
	AggregateType string       `gorm:"not null"`
	Payload       []byte       `gorm:"type:jsonb;not null"`
	Status        OutboxStatus `gorm:"type:varchar(16);not null;default:'pending';index"`
	CreatedAt     time.Time    `gorm:"autoCreateTime"`
	ProcessedAt   *time.Time
}
