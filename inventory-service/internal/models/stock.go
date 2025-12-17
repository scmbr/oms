package models

type Stock struct {
	ProductID string `gorm:"column:product_id;not null"`
	Available uint   `gorm:"column:available;not null"`
}
