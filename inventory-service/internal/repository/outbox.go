package repository

import (
	"context"

	"github.com/scmbr/oms/inventory-service/internal/models"
	"gorm.io/gorm"
)

type OutboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}
func (r *OutboxRepository) GetByStatus(ctx context.Context, status models.OutboxStatus) ([]models.OutboxEvent, error) {
	var outboxes []models.OutboxEvent
	if err := r.db.WithContext(ctx).Where("status = ?", status).Find(&outboxes).Error; err != nil {
		return nil, err
	}
	return outboxes, nil
}
func (r *OutboxRepository) UpdateStatus(ctx context.Context, externalID string, newStatus models.OutboxStatus) error {
	res := r.db.WithContext(ctx).Model(&models.OutboxEvent{}).Where("external_id = ?", externalID).Update("status", newStatus)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
