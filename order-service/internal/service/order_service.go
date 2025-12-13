package service

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepo
}

func NewOrderService(r repository.OrderRepo) *OrderService {
	return &OrderService{repo: r}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []models.OrderItem) (*dto.OrderDTO, error) {
	order, err := s.repo.Create(ctx, userID, items)
	if err != nil {
		return nil, err
	}
	return dto.ToOrderDTO(order), nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*dto.OrderDTO, error) {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return dto.ToOrderDTO(order), nil
}

func (s *OrderService) ListOrders(ctx context.Context, userID string) ([]dto.OrderDTO, error) {
	orders, err := s.repo.ListOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	var dtos []dto.OrderDTO
	for _, o := range orders {
		dtos = append(dtos, *dto.ToOrderDTO(&o))
	}
	return dtos, nil
}
