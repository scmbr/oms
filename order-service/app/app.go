package app

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/scmbr/oms/common/config"
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

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		return err
	}
	rabbitCfg := rabbit.Config{
		URL:      cfg.RabbitMQURL,
		Exchange: "saga.events",
		Queue:    "order-service",
	}
	conn, err := rabbit.NewConnection(rabbitCfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	publisher := rabbit.NewPublisher(conn, rabbitCfg.Exchange)

	repos := repository.NewRepositories(db)
	orderSvc := service.NewServices(repos)

	outboxWorker := worker.NewOutboxWorker(orderSvc.Outbox, publisher, 5*time.Second)

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
	go consumer.Start(ctx)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, handler)

	log.Printf("OrderService gRPC running on :%s", cfg.Port)
	return grpcServer.Serve(lis)
}
