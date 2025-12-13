package service

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
)

type Orders interface {
	CreateOrder(ctx context.Context, userID string, items []models.OrderItem) (*models.Order, error)
	GetOrder(ctx context.Context, orderID string) (*models.Order, error)
	ListOrders(ctx context.Context, userID string) ([]models.Order, error)
}

type Services struct {
	Order Orders
}

func NewServicesitories(repo repository.OrderRepo) *Services {
	return &Services{
		Order: NewOrderService(repo),
	}
}
