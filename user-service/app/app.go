package app

import (
	"log"
	"net"
	"time"

	"github.com/scmbr/oms/common/config"
	pb "github.com/scmbr/oms/user-service/internal/pb"
	"github.com/scmbr/oms/user-service/internal/repository"
	"github.com/scmbr/oms/user-service/internal/service"
	grpc_handler "github.com/scmbr/oms/user-service/internal/transport/grpc"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() error {

	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	repos := repository.NewRepositories(db)

	userSvc := service.NewUserService(repos, 15*time.Minute, 24*time.Hour)

	handler := grpc_handler.NewUserHandler(userSvc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, handler)

	log.Printf("UserService gRPC running on :%s", cfg.Port)
	return grpcServer.Serve(lis)
}
