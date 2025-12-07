package grpc

import (
	"context"

	pb "github.com/scmbr/oms/user-service/internal/pb"
	"github.com/scmbr/oms/user-service/internal/service"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	userDTO, err := h.userService.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		UserId:           userDTO.UserID,
		Email:            userDTO.Email,
		Role:             userDTO.Role,
		RefreshToken:     userDTO.RefreshToken,
		RefreshExpiresAt: userDTO.RefreshExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	userDTO, err := h.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		UserId:           userDTO.UserID,
		Email:            userDTO.Email,
		Role:             userDTO.Role,
		RefreshToken:     userDTO.RefreshToken,
		RefreshExpiresAt: userDTO.RefreshExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *UserHandler) ValidateRefreshToken(ctx context.Context, req *pb.ValidateRefreshRequest) (*pb.ValidateRefreshResponse, error) {
	userDTO, err := h.userService.ValidateRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateRefreshResponse{
		UserId: userDTO.UserID,
		Email:  userDTO.Email,
		Role:   userDTO.Role,
	}, nil
}
