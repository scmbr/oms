package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/docker/go-connections/nat"
	"github.com/scmbr/oms/common/hasher"
	"github.com/scmbr/oms/user-service/internal/models"
	"github.com/scmbr/oms/user-service/internal/repository"
	"github.com/scmbr/oms/user-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	svc *service.UserService
	ctx context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	db, cleanup := setupPostgresContainer()
	defer cleanup()

	repos := repository.NewRepositories(db)
	h := hasher.BcryptHasher{}
	svc = service.NewUserService(repos, h, time.Hour*24, time.Hour*24)

	os.Exit(m.Run())
}

func setupPostgresContainer() (*gorm.DB, func()) {
	ctx := context.Background()

	sqlWait := wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
		return fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port())
	}).WithStartupTimeout(60 * time.Second)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: sqlWait,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		panic(err)
	}
	port, err := container.MappedPort(ctx, nat.Port("5432/tcp"))
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port())

	var db *gorm.DB
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}

	if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		panic(err)
	}

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			fmt.Println("Failed to terminate container:", err)
		}
	}

	return db, cleanup
}

// Тест регистрации нового пользователя
func TestRegister(t *testing.T) {
	email := "integration@mail.com"
	password := "Password123!"
	dto, err := svc.Register(ctx, email, password)
	assert.NoError(t, err)
	assert.Equal(t, email, dto.Email)
	assert.NotEmpty(t, dto.RefreshToken)
}

// Тест повторной регистрации (должна быть ошибка)
func TestRegisterDuplicate(t *testing.T) {
	email := "integration@mail.com"
	password := "Password123!"
	dto, err := svc.Register(ctx, email, password)
	assert.Error(t, err)
	assert.Nil(t, dto)
}

// Тест успешного логина
func TestLogin(t *testing.T) {
	email := "integration@mail.com"
	password := "Password123!"
	dto, err := svc.Login(ctx, email, password)
	assert.NoError(t, err)
	assert.Equal(t, email, dto.Email)
	assert.NotEmpty(t, dto.RefreshToken)
}

// Тест неправильного пароля
func TestLoginWrongPassword(t *testing.T) {
	email := "integration@mail.com"
	password := "WrongPassword!"
	dto, err := svc.Login(ctx, email, password)
	assert.Error(t, err)
	assert.Nil(t, dto)
}

// Тест неправильного email
func TestLoginNonExistentEmail(t *testing.T) {
	email := "notexist@mail.com"
	password := "Password123!"
	dto, err := svc.Login(ctx, email, password)
	assert.Error(t, err)
	assert.Nil(t, dto)
}

// Тест валидации refresh токена
func TestValidateRefreshToken(t *testing.T) {
	email := "integration@mail.com"
	password := "Password123!"
	loginDto, _ := svc.Login(ctx, email, password)
	dto, err := svc.ValidateRefreshToken(ctx, loginDto.RefreshToken)
	assert.NoError(t, err)
	assert.Equal(t, email, dto.Email)
}

// Тест неверного токена
func TestValidateInvalidRefreshToken(t *testing.T) {
	dto, err := svc.ValidateRefreshToken(ctx, "invalid-token")
	assert.Error(t, err)
	assert.Nil(t, dto)
}
