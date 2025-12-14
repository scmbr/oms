package dto

import "time"

type SagaEvent struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	OrderID    string    `json:"order_id"`
	OccurredAt time.Time `json:"occurred_at"`
}
