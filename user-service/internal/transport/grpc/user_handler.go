package grpc

import (
	"context"
	"time"

	"github.com/scmbr/oms/user-service/internal/dto"
	pb "github.com/scmbr/oms/user-service/internal/pb"
	"github.com/scmbr/oms/user-service/internal/service"
	"github.com/scmbr/oms/user-service/pkg/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	input := dto.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := validation.Validate.Struct(&input); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err)
	}

	userDTO, err := h.userService.Register(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		UserId:           userDTO.UserID,
		Email:            userDTO.Email,
		Role:             userDTO.Role,
		RefreshToken:     userDTO.RefreshToken,
		RefreshExpiresAt: userDTO.RefreshExpiresAt.Format(time.RFC3339),
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	input := dto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := validation.Validate.Struct(&input); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err)
	}
	userDTO, err := h.userService.Login(ctx, input.Email, input.Password)
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
