package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/scmbr/oms/user-service/internal/models"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type RefreshTokenRepo interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	GetByToken(ctx context.Context, tokenStr string) (*models.RefreshToken, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repositories struct {
	User         UserRepo
	RefreshToken RefreshTokenRepo
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		RefreshToken: NewRefreshTokenRepository(db),
	}
}
