package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/scmbr/oms/common/hasher"
	"github.com/scmbr/oms/user-service/internal/models"
	"github.com/scmbr/oms/user-service/internal/repository"
	"github.com/scmbr/oms/user-service/internal/service"
)

type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) Create(ctx context.Context, u *models.User) error {
	return m.Called(ctx, u).Error(0)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockRefreshRepo struct{ mock.Mock }

func (m *MockRefreshRepo) Create(ctx context.Context, r *models.RefreshToken) error {
	return m.Called(ctx, r).Error(0)
}

func (m *MockRefreshRepo) GetByToken(ctx context.Context, tokenStr string) (*models.RefreshToken, error) {
	args := m.Called(ctx, tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func TestRegister_Success(t *testing.T) {
	userRepo := new(MockUserRepo)
	refreshRepo := new(MockRefreshRepo)

	repos := &repository.Repositories{
		User:         userRepo,
		RefreshToken: refreshRepo,
	}

	h := hasher.BcryptHasher{}
	svc := service.NewUserService(repos, h, time.Hour*24, time.Hour*24)

	ctx := context.Background()

	userRepo.On("GetByEmail", ctx, "test@mail.com").Return(nil, service.ErrUserNotFound)
	userRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)
	refreshRepo.On("Create", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	dto, err := svc.Register(ctx, "test@mail.com", "Password123!")

	assert.NoError(t, err)
	assert.Equal(t, "test@mail.com", dto.Email)
	assert.NotEmpty(t, dto.RefreshToken)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), dto.RefreshExpiresAt, 5*time.Second)
}

func TestLogin_Success(t *testing.T) {
	userRepo := new(MockUserRepo)
	refreshRepo := new(MockRefreshRepo)

	repos := &repository.Repositories{
		User:         userRepo,
		RefreshToken: refreshRepo,
	}

	h := hasher.BcryptHasher{}
	svc := service.NewUserService(repos, h, time.Hour*24, time.Hour*24)

	ctx := context.Background()

	pass := "correct-password"
	hash, _ := h.Hash(pass)

	user := &models.User{
		ID:           uuid.New(),
		Email:        "login@mail.com",
		PasswordHash: hash,
		Role:         "user",
	}

	userRepo.On("GetByEmail", ctx, "login@mail.com").Return(user, nil)
	refreshRepo.On("Create", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	dto, err := svc.Login(ctx, "login@mail.com", pass)

	assert.NoError(t, err)
	assert.Equal(t, "login@mail.com", dto.Email)
	assert.NotEmpty(t, dto.RefreshToken)
}
