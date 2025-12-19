package repository

import (
	"context"

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
func (r *StockRepository) UpdateQuantity(ctx context.Context, productID string, delta int) (*models.Stock, error) {
	return nil, nil
}
