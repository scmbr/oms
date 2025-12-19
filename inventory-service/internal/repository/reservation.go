package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/scmbr/oms/inventory-service/internal/models"
	"gorm.io/gorm"
)

type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}

}
func (r *ReservationRepository) Create(ctx context.Context, reservation *models.Reservation, externailID string) (*models.Reservation, error) {
	reservationID := uuid.New().String()
	reservation.ReservationID = reservationID
	exp := time.Now().UTC().Add(time.Minute * 15)
	reservation.ExpiredAt = &exp

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(reservation).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	payload, err := marshalPayload(reservation)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	outbox := models.OutboxEvent{
		ExternalID:    externailID,
		EventType:     "inventory.reserved",
		AggregateID:   reservationID,
		AggregateType: "reservation",
		Payload:       payload,
		Status:        models.OutboxPending,
	}
	if err := tx.Create(outbox).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return reservation, nil
}
func (r *ReservationRepository) GetById(ctx context.Context, reservationID string) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.WithContext(ctx).Where("reservation_id = ?", reservationID).First(&reservation).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}
func (r *ReservationRepository) GetAll(ctx context.Context) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.WithContext(ctx).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}
func (r *ReservationRepository) Delete(ctx context.Context, reservationID string) (*models.Reservation, error) {
	return nil, nil
}
func (r *ReservationRepository) UpdateStatus(ctx context.Context, reservationID string, newStatus models.ReservationStatus) (*models.Reservation, error) {
	return nil, nil
}
func marshalPayload(reservation *models.Reservation) ([]byte, error) {
	payload, err := json.Marshal(struct {
		ReservationID string                   `json:"reservation_id"`
		OrderID       string                   `json:"order_id"`
		ProductID     string                   `json:"product_id"`
		Quantity      uint                     `json:"quantity"`
		Status        models.ReservationStatus `json:"status"`
		CreatedAt     time.Time                `json:"created_at"`
		ExpiredAt     *time.Time               `json:"expired_at"`
	}{
		ReservationID: reservation.ReservationID,
		OrderID:       reservation.OrderID,
		ProductID:     reservation.ProductID,
		Quantity:      reservation.Quantity,
		Status:        reservation.Status,
		CreatedAt:     reservation.CreatedAt,
		ExpiredAt:     reservation.ExpiredAt,
	})
	if err != nil {

		return nil, nil
	}
	return payload, err
}
