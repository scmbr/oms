package models

import "time"

type Order struct {
	OrderID    string      `gorm:"primaryKey;column:order_id"`
	UserID     string      `gorm:"column:user_id;not null"`
	Status     OrderStatus `gorm:"column:status;not null"`
	TotalPrice float64     `gorm:"column:total_price;not null"`
	Items      []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time   `gorm:"column:created_at;autoCreateTime"`
}
