package models

import "time"

type Product struct {
	ProductID string    `gorm:"primaryKey;column:product_id"`
	Title     string    `gorm:"column:title"`
	Sku       string    `gorm:"column:sku;uniqueIndex"`
	Price     float64   `gorm:"column:price;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}
