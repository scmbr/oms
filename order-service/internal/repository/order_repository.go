package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, userID string, items []dto.OrderItemDTO, total float64) (*models.Order, error) {
	orderID := uuid.New().String()
	var orderItems []models.OrderItem
	for i := range items {
		var orderItem models.OrderItem
		orderItem.ItemID = uuid.New().String()
		orderItem.OrderID = orderID
		orderItem.Price = items[i].Price
		orderItem.Quantity = int(items[i].Quantity)
		orderItem.ProductID = items[i].ProductID

		orderItems = append(orderItems, orderItem)
	}

	order := &models.Order{
		OrderID:    orderID,
		UserID:     userID,
		Status:     "CREATED",
		TotalPrice: total,
		Items:      orderItems,
		CreatedAt:  time.Now(),
	}

	outbox := &models.OutboxEvent{
		ExternalID: uuid.New().String(),
		EventType:  "order.created",
		OrderID:    order.OrderID,
		Payload:    marshalPayload(order),
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Create(outbox).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
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

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID string, newStatus models.OrderStatus, eventID string) error {
	res := r.db.WithContext(ctx).Model(&models.Order{}).Where("order_id = ?", orderID).Update("status", newStatus)
	if res.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return res.Error

}
func marshalPayload(order *models.Order) []byte {
	payload, err := json.Marshal(struct {
		OrderID    string             `json:"order_id"`
		UserID     string             `json:"user_id"`
		Status     models.OrderStatus `json:"status"`
		TotalPrice float64            `json:"total_price"`
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
