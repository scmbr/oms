package service

import (
	"context"
	"encoding/json"

	"github.com/scmbr/oms/common/tx"
	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
	"gorm.io/gorm"
)

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
			EventType: "order.created",
			OrderID:   o.OrderID,
			Payload:   marshalPayload(o),
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
		OrderID    string  `json:"order_id"`
		UserID     string  `json:"user_id"`
		Status     string  `json:"status"`
		TotalPrice float64 `json:"total_price"`
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
