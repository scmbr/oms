package service

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepo
}

func NewOrderService(repo repository.OrderRepo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []models.OrderItem) (*models.Order, error) {
	return s.repo.Create(ctx, userID, items)
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	return s.repo.GetOrder(ctx, orderID)
}

func (s *OrderService) ListOrders(ctx context.Context, userID string) ([]models.Order, error) {
	return s.repo.ListOrders(ctx, userID)
}
