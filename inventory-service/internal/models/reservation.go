package models

import "time"

type ReservationStatus string

const (
	ReservationPending   ReservationStatus = "PENDING"
	ReservationReserved  ReservationStatus = "RESERVED"
	ReservationFailed    ReservationStatus = "FAILED"
	ReservationCancelled ReservationStatus = "CANCELLED"
	ReservationExpired   ReservationStatus = "EXPIRED"
)

type Reservation struct {
	ReservationID string            `gorm:"primaryKey;column:reservation_id"`
	OrderID       string            `gorm:"column:order_id;not null"`
	ProductID     string            `gorm:"column:product_id;not null"`
	Quantity      uint              `gorm:"column:quantity;not null"`
	Status        ReservationStatus `gorm:"column:status;not null"`
	CreatedAt     time.Time         `gorm:"column:created_at;autoCreateTime"`
	ExpiredAt     *time.Time        `gorm:"column:expired_at"`
}
