package app

import (
	"context"
	"net"
	"time"

	"github.com/scmbr/oms/common/config"
	"github.com/scmbr/oms/common/logger"
	pb "github.com/scmbr/oms/order-service/internal/pb"
	"github.com/scmbr/oms/order-service/internal/repository"
	"github.com/scmbr/oms/order-service/internal/service"
	grpc_handler "github.com/scmbr/oms/order-service/internal/transport/grpc"
	"github.com/scmbr/oms/order-service/internal/transport/rabbit"
	"github.com/scmbr/oms/order-service/internal/worker"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()
	logger.Init("order-service")
	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		logger.Error("database connection failed", err)
		return err
	}
	logger.Info("database connected successfully")
	rabbitCfg := rabbit.Config{
		URL:      cfg.RabbitMQURL,
		Exchange: "saga.events",
		Queue:    "order-service",
	}
	conn, err := rabbit.NewConnection(rabbitCfg)
	if err != nil {
		logger.Error("rabbitmq connection failed", err)
		return err
	}
	logger.Info("rabbitmq connected succesfully")
	defer conn.Close()

	publisher := rabbit.NewPublisher(conn, rabbitCfg.Exchange)

	repos := repository.NewRepositories(db)
	orderSvc := service.NewServices(repos)

	outboxWorker := worker.NewOutboxWorker(orderSvc.Outbox, publisher, 5*time.Second)
	logger.Info("outbox worker is running")
	go outboxWorker.Start(ctx)
	handler := grpc_handler.NewOrderHandler(orderSvc.Order)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return err
	}
	consumer, err := rabbit.NewConsumer(conn, rabbitCfg.Queue, orderSvc.Order)
	if err != nil {
		return err
	}
	logger.Info("consumer worker is running")
	go consumer.Start(ctx)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, handler)

	logger.Info("order-service gRPC server started", map[string]interface{}{
		"port":     cfg.Port,
		"protocol": "gRPC",
		"service":  "order-service",
	})
	return grpcServer.Serve(lis)
}
