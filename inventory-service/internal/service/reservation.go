package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scmbr/oms/common/tx"
	"github.com/scmbr/oms/inventory-service/internal/dto"
	"github.com/scmbr/oms/inventory-service/internal/models"
	"github.com/scmbr/oms/inventory-service/internal/repository"
)

type ReservationService struct {
	reservationRepo repository.Reservation
	stockRepo       repository.Stock
	outboxRepo      repository.Outbox
	tx              tx.TxManager
}

func NewReservationService(reservationRepo repository.Reservation, stockRepo repository.Stock, outboxRepo repository.Outbox, tx tx.TxManager) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		stockRepo:       stockRepo,
		outboxRepo:      outboxRepo,
		tx:              tx,
	}
}
func (s *ReservationService) Reserve(ctx context.Context, productID string, orderID string, externalID string, quantity int) error {
	return s.tx.WithTx(ctx, func(ctxTx context.Context) error {
		if quantity <= 0 {
			return fmt.Errorf("quantity is less than zero or equal zero for productID: %s and orderID: %s", productID, orderID)
		}
		outbox, err := s.outboxRepo.GetByExtID(ctx, externalID)
		if err != nil {
			return err
		}

		if outbox != nil && outbox.Status != models.OutboxFailed {
			return fmt.Errorf("already processed for externalID: %s", externalID)
		}
		eventType := models.EventTypeReserved
		status := models.ReservationReserved
		if err := s.stockRepo.UpdateQuantity(ctxTx, productID, -quantity); err != nil {
			status = models.ReservationFailed
			eventType = models.EventTypeFailed
		}

		exp := time.Now().UTC().Add(time.Minute * 15)
		reservation := &models.Reservation{
			ReservationID: uuid.NewString(),
			OrderID:       orderID,
			ProductID:     productID,
			Quantity:      uint(quantity),
			Status:        status,
			CreatedAt:     time.Now(),
			ExpiredAt:     &exp,
		}
		payload, err := marshalPayload(reservation)
		if err != nil {
			return err
		}
		if _, err := s.reservationRepo.Create(ctxTx, reservation, externalID, payload, eventType); err != nil {
			return err
		}
		return nil
	})
}
func (s *ReservationService) GetPending(ctx context.Context) ([]dto.ReservationResponse, error) {
	return nil, nil
}
func (s *ReservationService) UpdateStatus(ctx context.Context, id string, newStatus models.ReservationStatus) error {
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
