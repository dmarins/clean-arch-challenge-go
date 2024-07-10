package usecase

import (
	"github.com/dmarins/clean-arch-challenge-go/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	orderRepository entity.OrderRepositoryInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *ListOrderUseCase) Execute() ([]OrderOutputDTO, error) {

	orders, err := c.OrderRepository.GetAll()
	if err != nil {
		return nil, err
	}

	dtos := []OrderOutputDTO{}
	for _, item := range orders {
		dto := OrderOutputDTO{
			ID:         item.ID,
			Price:      item.Price,
			Tax:        item.Tax,
			FinalPrice: item.FinalPrice,
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}
