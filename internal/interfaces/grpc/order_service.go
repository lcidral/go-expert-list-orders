package grpc

import (
	"context"

	"go-expert-list-orders/internal/pb"
	"go-expert-list-orders/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	listOrdersUseCase *usecase.ListOrdersUseCase
}

func NewOrderService(listOrdersUseCase *usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		listOrdersUseCase: listOrdersUseCase,
	}
}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.listOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:          order.ID,
			CustomerId:  order.CustomerID,
			Status:      order.Status,
			TotalAmount: order.TotalAmount,
			CreatedAt:   order.CreatedAt.String(),
			UpdatedAt:   order.UpdatedAt.String(),
		})
	}

	return &pb.ListOrdersResponse{
		Orders: pbOrders,
	}, nil
}
