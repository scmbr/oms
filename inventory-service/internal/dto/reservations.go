package dto

import (
	"time"

	"github.com/scmbr/oms/inventory-service/internal/models"
)

type ReservationResponse struct {
	ReservationID string                   `json:"reservation_id"`
	OrderID       string                   `json:"order_id"`
	ProductID     string                   `json:"product_id"`
	Quantity      uint                     `json:"quantity"`
	Status        models.ReservationStatus `json:"status"`
	CreatedAt     time.Time                `json:"created_at"`
	ExpiredAt     *time.Time               `json:"expired_at"`
}
