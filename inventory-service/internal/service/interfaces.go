package service

import (
	"context"

	"github.com/scmbr/oms/inventory-service/internal/dto"
	"github.com/scmbr/oms/inventory-service/internal/models"
	"github.com/scmbr/oms/inventory-service/internal/repository"
)

type Product interface {
	GetProductPrices(ctx context.Context, productsIds []string) ([]dto.ProductsPricesResponse, error)
}
type Reservation interface {
	Reserve(ctx context.Context) error
	GetPending(ctx context.Context) ([]dto.ReservationResponse, error)
	UpdateStatus(ctx context.Context, id string, newStatus models.ReservationStatus) error
}

type Outbox interface {
	GetPending(ctx context.Context) ([]dto.OutboxResponse, error)
	MarkAsSent(ctx context.Context, externalID string) error
}
type Service struct {
	product     Product
	reservation Reservation
	outbox      Outbox
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		product:     NewProductService(repo.Product),
		reservation: NewReservationService(repo.Reservation),
		outbox:      NewOutboxService(repo.Outbox),
	}

}
