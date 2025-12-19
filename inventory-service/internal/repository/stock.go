package repository

import (
	"context"
	"fmt"

	"github.com/scmbr/oms/inventory-service/internal/models"
	"gorm.io/gorm"
)

type StockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}

}
func (r *StockRepository) Create(ctx context.Context, stock *models.Stock) (*models.Stock, error) {
	if err := r.db.WithContext(ctx).Create(stock).Error; err != nil {
		return nil, err
	}

	return stock, nil
}
func (r *StockRepository) GetById(ctx context.Context, productID string) (*models.Stock, error) {
	var stock models.Stock
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&stock).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}
func (r *StockRepository) GetAll(ctx context.Context) ([]models.Stock, error) {
	var stock []models.Stock
	if err := r.db.WithContext(ctx).Find(&stock).Error; err != nil {
		return nil, err
	}
	return stock, nil
}
func (r *StockRepository) Delete(ctx context.Context, stockID string) error {
	if err := r.db.WithContext(ctx).Where("stock_id = ?", stockID).Delete(&models.Reservation{}).Error; err != nil {
		return err
	}
	return nil
}
func (r *StockRepository) UpdateQuantity(ctx context.Context, productID string, delta int) error {
	if delta == 0 {
		return nil
	}

	var res *gorm.DB
	if delta > 0 {
		res = r.db.WithContext(ctx).
			Model(&models.Stock{}).
			Where("product_id = ?", productID).
			Update("available", gorm.Expr("available + ?", delta))
	} else {
		res = r.db.WithContext(ctx).
			Model(&models.Stock{}).
			Where("product_id = ? AND available >= ?", productID, uint(-delta)).
			Update("available", gorm.Expr("available - ?", uint(-delta)))
	}

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("product %s: not enough available stock or not found", productID)
	}

	return nil
}
