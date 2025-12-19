package service

import (
	"context"

	"github.com/scmbr/oms/inventory-service/internal/dto"
	"github.com/scmbr/oms/inventory-service/internal/repository"
)

type OutboxService struct {
	outboxRepo repository.Outbox
}

func NewOutboxService(outboxRepo repository.Outbox) *OutboxService {
	return &OutboxService{
		outboxRepo: outboxRepo,
	}
}

func (s *OutboxService) GetPending(ctx context.Context) ([]dto.OutboxResponse, error) {
	return nil, nil
}
func (s *OutboxService) MarkAsSent(ctx context.Context, externalID string) error {
	return nil
}
