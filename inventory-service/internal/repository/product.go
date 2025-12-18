package repository

import (
	"context"

	"github.com/google/uuid"
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
	productID := uuid.New().String()
	product.ProductID = productID
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
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
