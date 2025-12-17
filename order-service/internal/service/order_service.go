package service

import (
	"context"
	"fmt"
	"log"

	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
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
}

func NewOrderService(r *repository.Repositories) *OrderService {
	return &OrderService{orderRepo: r.Order, outboxRepo: r.Outbox}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []dto.OrderItemDTO) (*dto.OrderDTO, error) {
	var order *models.Order
	var total float64
	for i := range items {
		total += items[i].Price * float64(items[i].Quantity)
	}

	order, err := s.orderRepo.Create(ctx, userID, items, total)
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

func (s *OrderService) UpdateStatus(ctx context.Context, orderID string, newStatus models.OrderStatus, eventID string) error {
	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	allowed := validTransitions[order.Status]
	isValid := false
	for _, st := range allowed {
		if st == newStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		log.Printf("Invalid status transition: %s → %s", order.Status, newStatus)
		return nil
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, newStatus, eventID)
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
