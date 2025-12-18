package repository

import (
	"context"
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
func (r *ReservationRepository) Create(ctx context.Context, reservation *models.Reservation) (*models.Reservation, error) {
	reservationID := uuid.New().String()
	reservation.ReservationID = reservationID
	exp := time.Now().UTC().Add(time.Minute * 15)
	reservation.ExpiredAt = &exp
	if err := r.db.WithContext(ctx).Create(reservation).Error; err != nil {
		return nil, err
	}

	return reservation, nil
}
func (r *ReservationRepository) GetById(ctx context.Context, reservationID string) (*models.Reservation, error) {
	return nil, nil
}
func (r *ReservationRepository) GetAll(ctx context.Context) ([]models.Reservation, error) {
	return nil, nil
}
func (r *ReservationRepository) Delete(ctx context.Context, reservationID string) (*models.Reservation, error) {
	return nil, nil
}
func (r *ReservationRepository) UpdateStatus(ctx context.Context, reservationID string, newStatus models.ReservationStatus) (*models.Reservation, error) {
	return nil, nil
}
