package models

type Stock struct {
	ProductID string `gorm:"primaryKey;column:product_id;not null"`
	Available uint   `gorm:"column:available;not null"`
}
