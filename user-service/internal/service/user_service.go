package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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
	repos      *repository.Repositories
	tokenTTL   time.Duration
	refreshTTL time.Duration
}

func NewUserService(repos *repository.Repositories, tokenTTL, refreshTTL time.Duration) *UserService {
	return &UserService{
		repos:      repos,
		tokenTTL:   tokenTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *UserService) Register(ctx context.Context, email, password string) (*UserDTO, error) {
	existing, _ := s.repos.User.GetByEmail(ctx, email)
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
	if err := s.repos.User.Create(ctx, user); err != nil {
		return nil, err
	}

	refresh := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.repos.RefreshToken.Create(ctx, refresh); err != nil {
		return nil, err
	}

	return toUserDTO(user.ID.String(), user.Email, user.Role, refresh.Token, refresh.ExpiresAt), nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*UserDTO, error) {
	user, err := s.repos.User.GetByEmail(ctx, email)
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
	if err := s.repos.RefreshToken.Create(ctx, refresh); err != nil {
		return nil, err
	}

	return toUserDTO(user.ID.String(), user.Email, user.Role, refresh.Token, refresh.ExpiresAt), nil
}

func (s *UserService) ValidateRefreshToken(ctx context.Context, tokenStr string) (*UserDTO, error) {
	token, err := s.repos.RefreshToken.GetByToken(ctx, tokenStr)
	if err != nil || token == nil {
		return nil, ErrUserNotFound
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	user, err := s.repos.User.GetByID(ctx, token.UserID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	return toUserDTO(user.ID.String(), user.Email, user.Role, token.Token, token.ExpiresAt), nil
}
