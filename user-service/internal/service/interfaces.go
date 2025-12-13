package service

import (
	"context"
	"time"

	"github.com/scmbr/oms/common/hasher"
	"github.com/scmbr/oms/user-service/internal/dto"
	"github.com/scmbr/oms/user-service/internal/repository"
)

type Users interface {
	Register(ctx context.Context, email, password string) (*dto.UserDTO, error)
	Login(ctx context.Context, email, password string) (*dto.UserDTO, error)
	ValidateRefreshToken(ctx context.Context, tokenStr string) (*dto.UserDTO, error)
}

type Services struct {
	User Users
}

func NewServices(r *repository.Repositories,
	h hasher.PasswordHasher,
	tokenTTL,
	refreshTTL time.Duration) *Services {
	return &Services{
		User: NewUserService(r, h, tokenTTL, refreshTTL),
	}
}
