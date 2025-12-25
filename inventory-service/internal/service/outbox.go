package service

import (
	"context"

	"github.com/scmbr/oms/inventory-service/internal/dto"
	"github.com/scmbr/oms/inventory-service/internal/models"
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
	outboxes, err := s.outboxRepo.GetByStatus(ctx, models.OutboxPending)
	if err != nil {
		return nil, err
	}

	res := make([]dto.OutboxResponse, 0, len(outboxes))
	for _, o := range outboxes {
		res = append(res, dto.OutboxResponse{
			ExternalID:    o.ExternalID,
			EventType:     o.EventType,
			AggregateID:   o.AggregateID,
			AggregateType: o.AggregateType,
			Payload:       o.Payload,
			CreatedAt:     o.CreatedAt,
		})
	}
	return res, nil
}

func (s *OutboxService) MarkAsSent(ctx context.Context, externalID string) error {
	return s.outboxRepo.UpdateStatus(ctx, externalID, models.OutboxSent)
}
