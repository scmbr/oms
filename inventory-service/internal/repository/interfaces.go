package repository

import (
	"context"

	"github.com/scmbr/oms/inventory-service/internal/models"
	"gorm.io/gorm"
)

type Product interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	GetById(ctx context.Context, productID string) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	Delete(ctx context.Context, productID string) error
}
type Stock interface {
	Create(ctx context.Context, stock *models.Stock) (*models.Stock, error)
	GetById(ctx context.Context, productID string) (*models.Stock, error)
	GetAll(ctx context.Context) ([]models.Stock, error)
	Delete(ctx context.Context, stockID string) error
	UpdateQuantity(ctx context.Context, productID string, delta int) (*models.Stock, error)
}
type Reservation interface {
	Create(ctx context.Context, reservation *models.Reservation, externalID string) (*models.Reservation, error)
	GetById(ctx context.Context, reservationID string) (*models.Reservation, error)
	GetAll(ctx context.Context) ([]models.Reservation, error)
	Delete(ctx context.Context, reservationID string) error
	UpdateStatus(ctx context.Context, reservationID string, newStatus models.ReservationStatus) error
}
type Outbox interface {
	GetByStatus(ctx context.Context, status models.OutboxStatus) ([]models.OutboxEvent, error)
	UpdateStatus(ctx context.Context, externalID string, newStatus models.OutboxStatus) error
}
type Repository struct {
	product     Product
	stock       Stock
	reservation Reservation
	outbox      Outbox
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		product:     NewProductRepository(db),
		stock:       NewStockRepository(db),
		reservation: NewReservationRepository(db),
		outbox:      NewOutboxRepository(db),
	}
}
