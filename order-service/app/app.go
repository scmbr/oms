package app

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/scmbr/oms/common/config"
	"github.com/scmbr/oms/common/tx"
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
		URL:   cfg.RabbitMQURL,
		Queue: "orders",
	}
	conn, err := rabbit.NewConnection(rabbitCfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	publisher := rabbit.NewPublisher(conn, rabbitCfg.Queue)

	repos := repository.NewRepositories(db)
	txManager := tx.NewTxManager(db)
	orderSvc := service.NewServices(repos, txManager)

	outboxWorker := worker.NewOutboxWorker(orderSvc.Outbox, publisher, 5*time.Second, txManager)

	go outboxWorker.Start(ctx)
	handler := grpc_handler.NewOrderHandler(orderSvc.Order)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, handler)

	log.Printf("OrderService gRPC running on :%s", cfg.Port)
	return grpcServer.Serve(lis)
}
