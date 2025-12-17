package service

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
)

type OutboxService struct {
	repo repository.OutboxRepo
}

func NewOutboxService(r *repository.Repositories) *OutboxService {
	return &OutboxService{repo: r.Outbox}
}

func (s *OutboxService) GetPending(ctx context.Context) ([]models.OutboxEvent, error) {
	return s.repo.GetPending(ctx)
}

func (s *OutboxService) MarkAsProcessed(ctx context.Context, eventID string) error {
	return s.repo.MarkAsSent(ctx, eventID)
}
