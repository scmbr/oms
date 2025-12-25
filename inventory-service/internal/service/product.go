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
	products, err := s.productRepo.GetAllByIds(ctx, productsIds)
	if err != nil {
		return nil, err
	}
	res := make([]dto.ProductsPricesResponse, 0, len(products))
	for _, p := range products {
		res = append(res, dto.ProductsPricesResponse{
			ProductID: p.ProductID,
			Price:     p.Price,
		})
	}
	return res, nil
}
