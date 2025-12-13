package dto

import (
	"time"

	"github.com/scmbr/oms/order-service/internal/models"
)

type OrderItemDTO struct {
	ItemID    string
	ProductID string
	Quantity  int32
	Price     float64
}

type OrderDTO struct {
	OrderID    string
	UserID     string
	Status     string
	TotalPrice float64
	Items      []OrderItemDTO
	CreatedAt  string
}

func ToOrderDTO(o *models.Order) *OrderDTO {
	var items []OrderItemDTO
	for _, i := range o.Items {
		items = append(items, OrderItemDTO{
			ItemID:    i.ItemID,
			ProductID: i.ProductID,
			Quantity:  int32(i.Quantity),
			Price:     i.Price,
		})
	}
	return &OrderDTO{
		OrderID:    o.OrderID,
		UserID:     o.UserID,
		Status:     o.Status,
		TotalPrice: o.TotalPrice,
		Items:      items,
		CreatedAt:  o.CreatedAt.Format(time.RFC3339),
	}
}
