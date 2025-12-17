package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
	"github.com/scmbr/oms/order-service/internal/service"
)

type MockOrderRepo struct{ mock.Mock }

func (m *MockOrderRepo) Create(ctx context.Context, userID string, items []dto.OrderItemDTO, total float64) (*models.Order, error) {
	orderID := uuid.New().String()
	var orderItems []models.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, models.OrderItem{
			ItemID:    uuid.New().String(),
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}

	order := &models.Order{
		OrderID:    orderID,
		UserID:     userID,
		Status:     "CREATED",
		TotalPrice: total,
		Items:      orderItems,
		CreatedAt:  time.Now(),
	}

	return order, nil
}
func (m *MockOrderRepo) UpdateStatus(ctx context.Context, orderID string, newStatus models.OrderStatus, eventID string) error {
	return m.Called(ctx, orderID, newStatus, eventID).Error(0)
}

func (m *MockOrderRepo) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepo) ListOrders(ctx context.Context, userID string) ([]models.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Order), args.Error(1)
}

type MockOutboxRepo struct{ mock.Mock }

func (m *MockOutboxRepo) Create(ctx context.Context, event *models.OutboxEvent) error {
	return m.Called(ctx, event).Error(0)
}

func (m *MockOutboxRepo) GetPending(ctx context.Context) ([]models.OutboxEvent, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.OutboxEvent), args.Error(1)
}

func (m *MockOutboxRepo) MarkAsSent(ctx context.Context, eventID string) error {
	return m.Called(ctx, eventID).Error(0)
}

type MockTxManager struct{ mock.Mock }

func (m *MockTxManager) Do(ctx context.Context, fn func(tx *gorm.DB) error) error {
	args := m.Called(ctx, fn)
	if fn != nil {
		_ = fn(nil)
	}
	return args.Error(0)
}

func (m *MockTxManager) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	args := m.Called(ctx, fn)
	if fn != nil {
		_ = fn(nil)
	}
	return args.Error(0)
}
func TestCreateOrder_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	orderRepo := new(MockOrderRepo)
	outboxRepo := new(MockOutboxRepo)

	repos := &repository.Repositories{Order: orderRepo, Outbox: outboxRepo}
	svc := service.NewOrderService(repos)
	itemsDTO := []dto.OrderItemDTO{
		{
			ItemID:    uuid.New().String(),
			ProductID: uuid.New().String(),
			Quantity:  2,
			Price:     50.0,
		},
		{
			ItemID:    uuid.New().String(),
			ProductID: uuid.New().String(),
			Quantity:  1,
			Price:     30.0,
		},
		{
			ItemID:    uuid.New().String(),
			ProductID: uuid.New().String(),
			Quantity:  5,
			Price:     10.0,
		},
	}
	res, err := svc.CreateOrder(ctx, userID, itemsDTO)

	assert.NoError(t, err)
	assert.Equal(t, userID, res.UserID)
	assert.Equal(t, 180.0, res.TotalPrice)
	assert.Equal(t, res.Status, models.StatusCreated)
}
func TestGetOrder_Success(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New().String()

	orderRepo := new(MockOrderRepo)
	outboxRepo := new(MockOutboxRepo)

	repos := &repository.Repositories{Order: orderRepo, Outbox: outboxRepo}
	svc := service.NewOrderService(repos)

	order := &models.Order{OrderID: orderID, UserID: uuid.New().String(), Status: models.StatusCreated}

	orderRepo.On("GetOrder", ctx, orderID).Return(order, nil)

	res, err := svc.GetOrder(ctx, orderID)
	assert.NoError(t, err)
	assert.Equal(t, orderID, res.OrderID)
	assert.Equal(t, models.StatusCreated, res.Status)
}

func TestListOrders_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()

	orderRepo := new(MockOrderRepo)
	outboxRepo := new(MockOutboxRepo)

	repos := &repository.Repositories{Order: orderRepo, Outbox: outboxRepo}
	svc := service.NewOrderService(repos)

	orders := []models.Order{
		{OrderID: uuid.New().String(), UserID: userID},
		{OrderID: uuid.New().String(), UserID: userID},
	}

	orderRepo.On("ListOrders", ctx, userID).Return(orders, nil)

	res, err := svc.ListOrders(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}
