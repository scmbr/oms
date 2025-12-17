package repository

import (
	"context"
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

func (r *OrderRepository) Create(ctx context.Context, userID string, items []dto.OrderItemDTO, total float64) (*models.Order, error) {
	orderID := uuid.New().String()

	var total float64
	for i := range items {
		items[i].ItemID = uuid.New().String()
		items[i].OrderID = orderID
		total += items[i].Price * float64(items[i].Quantity)
	}

	order := &models.Order{
		OrderID:    orderID,
		UserID:     userID,
		Status:     "CREATED",
		TotalPrice: total,
		Items:      items,
		CreatedAt:  time.Now(),
	}

	if err := tx.Create(order).Error; err != nil {
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
