package grpc

import (
	"context"

	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	pb "github.com/scmbr/oms/order-service/internal/pb"
	"github.com/scmbr/oms/order-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	orderService service.Order
}

func NewOrderHandler(svc service.Order) *OrderHandler {
	return &OrderHandler{orderService: svc}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var items []dto.OrderItemDTO
	for _, i := range req.Items {
		items = append(items, dto.OrderItemDTO{
			ProductID: i.ProductId,
			Quantity:  i.Quantity,
		})
	}

	orderDTO, err := h.orderService.CreateOrder(ctx, req.UserId, itemsToModels(items))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &pb.CreateOrderResponse{
		OrderId: orderDTO.OrderID,
		Order:   dtoToProto(orderDTO),
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	orderDTO, err := h.orderService.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}

	return &pb.GetOrderResponse{
		Order: dtoToProto(orderDTO),
	}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orderDTOs, err := h.orderService.ListOrders(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	var protoOrders []*pb.Order
	for _, o := range orderDTOs {
		protoOrders = append(protoOrders, dtoToProto(&o))
	}

	return &pb.ListOrdersResponse{Orders: protoOrders}, nil
}

func itemsToModels(items []dto.OrderItemDTO) []models.OrderItem {
	var result []models.OrderItem
	for _, i := range items {
		result = append(result, models.OrderItem{
			ItemID:    i.ItemID,
			ProductID: i.ProductID,
			Quantity:  int(i.Quantity),
			Price:     0,
		})
	}
	return result
}

func dtoToProto(o *dto.OrderDTO) *pb.Order {
	var protoItems []*pb.OrderItem
	for _, i := range o.Items {
		protoItems = append(protoItems, &pb.OrderItem{
			ItemId:    i.ItemID,
			ProductId: i.ProductID,
			Quantity:  i.Quantity,
			Price:     i.Price,
		})
	}

	return &pb.Order{
		OrderId:    o.OrderID,
		UserId:     o.UserID,
		Status:     string(o.Status),
		TotalPrice: o.TotalPrice,
		Items:      protoItems,
		CreatedAt:  o.CreatedAt,
	}
}
