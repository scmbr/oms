package service

import (
	"context"

	"github.com/scmbr/oms/common/tx"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
	"gorm.io/gorm"
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

func (s *OutboxService) MarkAsProcessed(ctx context.Context, txManager tx.TxManager, eventID string) error {
	return txManager.WithTx(ctx, func(tx *gorm.DB) error {
		return s.repo.MarkAsSent(ctx, tx, eventID)
	})
}
