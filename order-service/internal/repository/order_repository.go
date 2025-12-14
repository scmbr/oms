package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/scmbr/oms/order-service/internal/models"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, tx *gorm.DB, userID string, items []models.OrderItem) (*models.Order, error) {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}

	order := &models.Order{
		OrderID:    uuid.New().String(),
		UserID:     userID,
		Status:     "CREATED",
		TotalPrice: total,
		Items:      items,
		CreatedAt:  time.Now(),
	}

	if err := tx.Create(order).Error; err != nil {
		return nil, err
	}

	eventPayload, _ := json.Marshal(struct {
		OrderID string `json:"order_id"`
		UserID  string `json:"user_id"`
		Status  string `json:"status"`
	}{
		OrderID: order.OrderID,
		UserID:  userID,
		Status:  order.Status,
	})

	outbox := &models.OutboxEvent{
		EventType: "order.created",
		OrderID:   order.OrderID,
		Payload:   eventPayload,
	}

	if err := tx.Create(outbox).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).Preload("Items").First(&order, "order_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListOrders(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
