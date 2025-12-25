package tx

import (
	"context"

	"gorm.io/gorm"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type txManager struct {
	db *gorm.DB
}
type txKey struct{}

func TxKey() any {
	return txKey{}
}
func NewTxManager(db *gorm.DB) *txManager {
	return &txManager{db: db}
}

func (m *txManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxTx := context.WithValue(ctx, txKey{}, tx)
		return fn(ctxTx)
	})
}
