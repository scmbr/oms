package repository

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/models"
	"gorm.io/gorm"
)

type OutboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}
func (r *OutboxRepository) Create(ctx context.Context, tx *gorm.DB, event *models.OutboxEvent) error {
	return tx.WithContext(ctx).Create(event).Error
}

func (r *OutboxRepository) GetPending(ctx context.Context) ([]models.OutboxEvent, error) {
	var events []models.OutboxEvent
	if err := r.db.WithContext(ctx).
		Where("processed_at IS NULL").
		Order("created_at ASC").
		Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *OutboxRepository) MarkAsSent(ctx context.Context, tx *gorm.DB, eventID string) error {
	return tx.WithContext(ctx).
		Model(&models.OutboxEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"status":       "processed",
			"processed_at": gorm.Expr("NOW()"),
		}).Error
}
