package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/scmbr/oms/common/tx"
	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
	"gorm.io/gorm"
)

var validTransitions = map[models.OrderStatus][]models.OrderStatus{
	models.StatusCreated:   {models.StatusReserved, models.StatusFailed, models.StatusCancelled},
	models.StatusReserved:  {models.StatusPaid, models.StatusFailed, models.StatusCancelled},
	models.StatusPaid:      {},
	models.StatusCancelled: {},
	models.StatusFailed:    {},
}

type OrderService struct {
	orderRepo  repository.OrderRepo
	outboxRepo repository.OutboxRepo
	txManager  tx.TxManager
}

func NewOrderService(r *repository.Repositories, txManager tx.TxManager) *OrderService {
	return &OrderService{orderRepo: r.Order, outboxRepo: r.Outbox, txManager: txManager}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []models.OrderItem) (*dto.OrderDTO, error) {
	var order *models.Order

	err := s.txManager.WithTx(ctx, func(tx *gorm.DB) error {
		o, err := s.orderRepo.Create(ctx, tx, userID, items)
		if err != nil {
			return err
		}

		outbox := &models.OutboxEvent{
			ExternalID: uuid.New().String(),
			EventType:  "order.created",
			OrderID:    o.OrderID,
			Payload:    marshalPayload(o),
		}

		if err := s.outboxRepo.Create(ctx, tx, outbox); err != nil {
			return err
		}

		order = o
		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.ToOrderDTO(order), nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*dto.OrderDTO, error) {
	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return dto.ToOrderDTO(order), nil
}

func (s *OrderService) ListOrders(ctx context.Context, userID string) ([]dto.OrderDTO, error) {
	orders, err := s.orderRepo.ListOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	var dtos []dto.OrderDTO
	for _, o := range orders {
		dtos = append(dtos, *dto.ToOrderDTO(&o))
	}
	return dtos, nil
}
func marshalPayload(order *models.Order) []byte {
	payload, err := json.Marshal(struct {
		OrderID    string             `json:"order_id"`
		UserID     string             `json:"user_id"`
		Status     models.OrderStatus `json:"status"`
		TotalPrice float64            `json:"total_price"`
	}{
		OrderID:    order.OrderID,
		UserID:     order.UserID,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
	})
	if err != nil {

		return nil
	}
	return payload
}
func (s *OrderService) UpdateStatus(ctx context.Context, orderID string, newStatus models.OrderStatus, eventID string) error {
	return s.txManager.WithTx(ctx, func(tx *gorm.DB) error {
		order, err := s.orderRepo.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}

		var exists bool
		if err := tx.Model(&models.OutboxEvent{}).Select("1").
			Where("external_id = ?", eventID).
			Limit(1).Scan(&exists).Error; err != nil {
			return err
		}
		if exists {
			return nil
		}

		validNext, ok := validTransitions[order.Status]
		if !ok {
			return fmt.Errorf("текущий статус %s неизвестен", order.Status)
		}

		allowed := false
		for _, st := range validNext {
			if st == newStatus {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("невозможный переход %s -> %s", order.Status, newStatus)
		}

		order.Status = newStatus
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		outbox := &models.OutboxEvent{
			ExternalID: eventID,
			EventType:  "order.status_changed",
			OrderID:    orderID,
			Payload:    marshalPayload(order),
		}
		if err := s.outboxRepo.Create(ctx, tx, outbox); err != nil {
			return err
		}

		return nil
	})
}
func (s *OrderService) ParseStatus(str string) (models.OrderStatus, error) {
	switch str {
	case string(models.StatusCreated):
		return models.StatusCreated, nil
	case string(models.StatusReserved):
		return models.StatusReserved, nil
	case string(models.StatusPaid):
		return models.StatusPaid, nil
	case string(models.StatusCancelled):
		return models.StatusCancelled, nil
	case string(models.StatusFailed):
		return models.StatusFailed, nil
	default:
		return "", fmt.Errorf("unknown order status: %s", s)
	}
}
