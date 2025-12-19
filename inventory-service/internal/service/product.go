package service

import (
	"context"

	"github.com/scmbr/oms/inventory-service/internal/dto"
	"github.com/scmbr/oms/inventory-service/internal/repository"
)

type ProductService struct {
	productRepo repository.Product
}

func NewProductService(productRepo repository.Product) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}
func (s *ProductService) GetProductPrices(ctx context.Context, productsIds []string) ([]dto.ProductsPricesResponse, error) {
	return nil, nil
}
