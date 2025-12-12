package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/scmbr/oms/user-service/internal/dto"
	"github.com/scmbr/oms/user-service/internal/models"
	"github.com/scmbr/oms/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists      = errors.New("user with this email already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
)

type UserService struct {
	users         repository.UserRepo
	refreshTokens repository.RefreshTokenRepo
	tokenTTL      time.Duration
	refreshTTL    time.Duration
}

func NewUserService(r *repository.Repositories, tokenTTL, refreshTTL time.Duration) *UserService {
	return &UserService{
		users:         r.User,
		refreshTokens: r.RefreshToken,
		tokenTTL:      tokenTTL,
		refreshTTL:    refreshTTL,
	}
}

func (s *UserService) Register(ctx context.Context, email, password string) (*dto.UserDTO, error) {
	existing, _ := s.users.GetByEmail(ctx, email)
	if existing != nil {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user",
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	refresh := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.refreshTokens.Create(ctx, refresh); err != nil {
		return nil, err
	}
	return dto.ToUserDTO(user, refresh.Token, refresh.ExpiresAt), nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*dto.UserDTO, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	refresh := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.refreshTokens.Create(ctx, refresh); err != nil {
		return nil, err
	}

	return dto.ToUserDTO(user, refresh.Token, refresh.ExpiresAt), nil
}

func (s *UserService) ValidateRefreshToken(ctx context.Context, tokenStr string) (*dto.UserDTO, error) {
	token, err := s.refreshTokens.GetByToken(ctx, tokenStr)
	if err != nil || token == nil {
		return nil, ErrUserNotFound
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	user, err := s.users.GetByID(ctx, token.UserID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	return dto.ToUserDTO(user, token.Token, token.ExpiresAt), nil
}
