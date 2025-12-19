package service

import (
	"context"

	"github.com/scmbr/oms/inventory-service/internal/dto"
	"github.com/scmbr/oms/inventory-service/internal/models"
	"github.com/scmbr/oms/inventory-service/internal/repository"
)

type ReservationService struct {
	reservationRepo repository.Reservation
}

func NewReservationService(reservationRepo repository.Reservation) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
	}
}
func (s *ReservationService) Reserve(ctx context.Context) error {
	return nil
}
func (s *ReservationService) GetPending(ctx context.Context) ([]dto.ReservationResponse, error) {
	return nil, nil
}
func (s *ReservationService) UpdateStatus(ctx context.Context, id string, newStatus models.ReservationStatus) error {
	return nil
}
