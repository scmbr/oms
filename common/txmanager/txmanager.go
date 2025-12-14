package tx

import (
	"context"

	"gorm.io/gorm"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type txManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *txManager {
	return &txManager{db: db}
}

func (m *txManager) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.db.WithContext(ctx).Transaction(fn)
}
