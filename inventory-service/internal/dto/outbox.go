package dto

import "time"

type OutboxResponse struct {
	ExternalID    string    `json:"external_id"`
	EventType     string    `json:"event_type"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	Payload       []byte    `json:"payload"`
	CreatedAt     time.Time `json:"created_at"`
}
