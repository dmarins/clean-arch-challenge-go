package usecase

import (
	"github.com/dmarins/clean-arch-challenge-go/internal/entity"
	"github.com/dmarins/clean-arch-challenge-go/pkg/events"
)

type OrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
	EventDispatcher   events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	orderRepository entity.OrderRepositoryInterface,
	orderCreatedEvent events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository:   orderRepository,
		OrderCreatedEvent: orderCreatedEvent,
		EventDispatcher:   eventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}

	order.CalculateFinalPrice()

	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}

	c.OrderCreatedEvent.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreatedEvent)

	return dto, nil
}
