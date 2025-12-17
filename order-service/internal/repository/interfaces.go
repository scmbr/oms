package repository

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	"gorm.io/gorm"
)

type OrderRepo interface {
	Create(ctx context.Context, userID string, items []dto.OrderItemDTO, total float64) (*models.Order, error)
	GetOrder(ctx context.Context, orderID string) (*models.Order, error)
	ListOrders(ctx context.Context, userID string) ([]models.Order, error)
	UpdateStatus(ctx context.Context, orderID string, newStatus models.OrderStatus, eventID string) error
}
type OutboxRepo interface {
	GetPending(ctx context.Context) ([]models.OutboxEvent, error)
	MarkAsSent(ctx context.Context, eventID string) error
}
type Repositories struct {
	Order  OrderRepo
	Outbox OutboxRepo
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Order:  NewOrderRepository(db),
		Outbox: NewOutboxRepository(db),
	}
}
