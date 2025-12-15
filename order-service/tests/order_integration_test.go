package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/scmbr/oms/common/tx"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/scmbr/oms/order-service/internal/repository"
	"github.com/scmbr/oms/order-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	orderSvc  *service.OrderService
	outboxSvc *service.OutboxService
	ctx       context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	db, cleanup := setupPostgresContainer()
	defer cleanup()

	repos := repository.NewRepositories(db)
	txManager := tx.NewTxManager(db)
	orderSvc = service.NewOrderService(repos, txManager)
	outboxSvc = service.NewOutboxService(repos)
	os.Exit(m.Run())
}

func setupPostgresContainer() (*gorm.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port())
		}).WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432/tcp")
	dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port())

	var db *gorm.DB
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}

	if err := db.AutoMigrate(&models.Order{}, &models.OrderItem{}, &models.OutboxEvent{}); err != nil {
		panic(err)
	}

	cleanup := func() {
		_ = container.Terminate(ctx)
	}

	return db, cleanup
}

// ---------------------- Тесты ----------------------

// Создание заказа
func TestCreateOrder(t *testing.T) {
	userID := uuid.New().String()
	items := []models.OrderItem{
		{ProductID: uuid.New().String(), Quantity: 2, Price: 10},
	}

	order, err := orderSvc.CreateOrder(ctx, userID, items)
	assert.NoError(t, err)
	assert.Equal(t, userID, order.UserID)
	assert.Equal(t, len(items), len(order.Items))
	assert.Equal(t, models.StatusCreated, order.Status)
}

// Получение заказа
func TestGetOrder(t *testing.T) {
	userID := uuid.New().String()
	items := []models.OrderItem{
		{ProductID: uuid.New().String(), Quantity: 1, Price: 5},
	}

	createdOrder, _ := orderSvc.CreateOrder(ctx, userID, items)
	got, err := orderSvc.GetOrder(ctx, createdOrder.OrderID)
	assert.NoError(t, err)
	assert.Equal(t, createdOrder.OrderID, got.OrderID)
}

// Список заказов пользователя
func TestListOrders(t *testing.T) {
	userID := uuid.New().String()
	// Создаем 2 заказа
	_, _ = orderSvc.CreateOrder(ctx, userID, []models.OrderItem{{ProductID: uuid.New().String(), Quantity: 1, Price: 5}})
	_, _ = orderSvc.CreateOrder(ctx, userID, []models.OrderItem{{ProductID: uuid.New().String(), Quantity: 2, Price: 10}})

	list, err := orderSvc.ListOrders(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

// Обновление статуса заказа
func TestUpdateStatus_ValidTransition(t *testing.T) {
	userID := uuid.New().String()
	items := []models.OrderItem{{ProductID: uuid.New().String(), Quantity: 1, Price: 5}}
	order, _ := orderSvc.CreateOrder(ctx, userID, items)

	err := orderSvc.UpdateStatus(ctx, order.OrderID, models.StatusReserved, uuid.New().String())
	assert.NoError(t, err)

	updated, _ := orderSvc.GetOrder(ctx, order.OrderID)
	assert.Equal(t, models.StatusReserved, updated.Status)
}

// Невалидная транзакция (нельзя напрямую в PAID из CREATED)
func TestUpdateStatus_InvalidTransition(t *testing.T) {
	userID := uuid.New().String()
	items := []models.OrderItem{{ProductID: uuid.New().String(), Quantity: 1, Price: 5}}
	order, _ := orderSvc.CreateOrder(ctx, userID, items)

	err := orderSvc.UpdateStatus(ctx, order.OrderID, models.StatusPaid, uuid.New().String())
	assert.Error(t, err)

	unchanged, _ := orderSvc.GetOrder(ctx, order.OrderID)
	assert.Equal(t, models.StatusCreated, unchanged.Status)
}

// Получение непереданных outbox событий
func TestGetPendingOutboxEvents(t *testing.T) {
	events, err := outboxSvc.GetPending(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, events)

	for _, e := range events {
		assert.NotEmpty(t, e.ExternalID)
		assert.Equal(t, "pending", e.Status)
	}
}
