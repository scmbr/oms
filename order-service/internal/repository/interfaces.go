package repository

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/models"
	"gorm.io/gorm"
)

type OrderRepo interface {
	Create(ctx context.Context, userID string, items []models.OrderItem) (*models.Order, error)
	GetOrder(ctx context.Context, orderID string) (*models.Order, error)
	ListOrders(ctx context.Context, userID string) ([]models.Order, error)
}

type Repositories struct {
	Order OrderRepo
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Order: NewOrderRepository(db),
	}
}
