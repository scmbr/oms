package models

type OrderItem struct {
	ItemID    string  `gorm:"primaryKey;column:item_id"`
	OrderID   string  `gorm:"column:order_id;not null;index"`
	ProductID string  `gorm:"column:product_id;not null"`
	Quantity  int     `gorm:"column:quantity;not null"`
	Price     float64 `gorm:"column:price;not null"`
}
