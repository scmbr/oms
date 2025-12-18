package repository

import (
	"context"

	"github.com/scmbr/oms/inventory-service/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	return nil, nil

}
func (r *ProductRepository) GetById(ctx context.Context, productID string) (*models.Product, error) {
	return nil, nil
}
func (r *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	return nil, nil
}
func (r *ProductRepository) Delete(ctx context.Context, productID string) (*models.Product, error) {
	return nil, nil
}
