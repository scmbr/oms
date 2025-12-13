package service

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
)

type Orders interface {
	CreateOrder(ctx context.Context, userID string, items []models.OrderItem) (*dto.OrderDTO, error)
	GetOrder(ctx context.Context, orderID string) (*dto.OrderDTO, error)
	ListOrders(ctx context.Context, userID string) ([]dto.OrderDTO, error)
}

type Services struct {
	Order Orders
}

func NewServices(repo *repository.Repositories) *Services {
	return &Services{
		Order: NewOrderService(repo.Order),
	}
}
