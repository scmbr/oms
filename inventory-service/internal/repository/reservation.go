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
func (r *ReservationRepository) Create(ctxTx context.Context, reservation *models.Reservation, externailID string, payload []byte, eventType string) (*models.Reservation, error) {
	tx, ok := ctxTx.Value(tx.TxKey()).(*gorm.DB)
	if !ok {
		tx = r.db
	}
	if err := tx.Create(reservation).Error; err != nil {
		return nil, err
	}
	outbox := models.OutboxEvent{
		ExternalID:    externailID,
		EventType:     models.EventTypeReserved,
		AggregateID:   reservation.ReservationID,
		AggregateType: "reservation",
		Payload:       payload,
		Status:        models.OutboxPending,
	}
	if err := tx.Create(outbox).Error; err != nil {
		return nil, err
	}
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
func (r *ReservationRepository) Delete(ctx context.Context, reservationID string) error {
	if err := r.db.WithContext(ctx).Where("reservation_id = ?", reservationID).Delete(&models.Reservation{}).Error; err != nil {
		return err
	}
	return nil
}
func (r *ReservationRepository) UpdateStatus(ctx context.Context, reservationID string, newStatus models.ReservationStatus) error {
	res := r.db.WithContext(ctx).Model(&models.Reservation{}).Where("reservation_id = ?", reservationID).Update("status", newStatus)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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
