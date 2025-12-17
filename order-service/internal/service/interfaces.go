package service

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
)

type Order interface {
	CreateOrder(ctx context.Context, userID string, items []dto.OrderItemDTO) (*dto.OrderDTO, error)
	GetOrder(ctx context.Context, orderID string) (*dto.OrderDTO, error)
	ListOrders(ctx context.Context, userID string) ([]dto.OrderDTO, error)
	UpdateStatus(ctx context.Context, orderID string, newStatus models.OrderStatus, eventID string) error
	ParseStatus(s string) (models.OrderStatus, error)
}

type Outbox interface {
	GetPending(ctx context.Context) ([]models.OutboxEvent, error)
	MarkAsProcessed(ctx context.Context, eventID string) error
}
type Services struct {
	Order  Order
	Outbox Outbox
}

func NewServices(repo *repository.Repositories) *Services {
	return &Services{
		Order:  NewOrderService(repo),
		Outbox: NewOutboxService(repo),
	}
}
