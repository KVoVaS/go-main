package grpcsrv

import (
	"context"
	warehousev1 "warehouse/api/gen/warehouse/v1"
	"warehouse/internal/service"
)

type WarehouseServer struct {
	warehousev1.UnimplementedWarehouseServiceServer
	svc *service.OrderService
}

func NewWarehouseServer(svc *service.OrderService) *WarehouseServer {
	return &WarehouseServer{svc: svc}
}

func (s *WarehouseServer) CreateOrder(ctx context.Context, req *warehousev1.CreateOrderRequest) (*warehousev1.CreateOrderResponse, error) {
	order, err := s.svc.CreateOrder(ctx, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}
	return &warehousev1.CreateOrderResponse{
		OrderId: order.ID,
		Status:  order.Status,
	}, nil
}

func (s *WarehouseServer) GetOrderStatus(ctx context.Context, req *warehousev1.GetOrderStatusRequest) (*warehousev1.GetOrderStatusResponse, error) {
	order, err := s.svc.GetOrderStatus(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &warehousev1.GetOrderStatusResponse{
		OrderId: order.ID,
		Status:  order.Status,
	}, nil
}
