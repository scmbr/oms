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
	return nil, nil
}
func (r *StockRepository) GetById(ctx context.Context, stockID string) (*models.Stock, error) {
	return nil, nil
}
func (r *StockRepository) GetAll(ctx context.Context) ([]models.Stock, error) {
	return nil, nil
}
func (r *StockRepository) Delete(ctx context.Context, stockID string) (*models.Stock, error) {
	return nil, nil
}
func (r *StockRepository) UpdateQuantity(ctx context.Context, productID string, delta int) (*models.Stock, error) {
	return nil, nil
}
