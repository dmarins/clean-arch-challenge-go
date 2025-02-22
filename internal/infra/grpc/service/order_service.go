package service

import (
	"context"

	"github.com/dmarins/clean-arch-challenge-go/internal/infra/grpc/pb"
	"github.com/dmarins/clean-arch-challenge-go/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.Blank) (*pb.OrderList, error) {
	output, err := s.ListOrderUseCase.Execute()
	if err != nil {
		return nil, err
	}

	orders := make([]*pb.CreateOrderResponse, 0)
	for _, item := range output {
		order := pb.CreateOrderResponse{
			Id:         item.ID,
			Price:      float32(item.Price),
			Tax:        float32(item.Tax),
			FinalPrice: float32(item.FinalPrice),
		}

		orders = append(orders, &order)
	}

	orderList := &pb.OrderList{
		Orders: orders,
	}

	return orderList, nil
}
