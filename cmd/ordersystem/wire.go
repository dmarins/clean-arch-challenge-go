//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/dmarins/clean-arch-challenge-go/internal/entity"
	"github.com/dmarins/clean-arch-challenge-go/internal/event"
	"github.com/dmarins/clean-arch-challenge-go/internal/infra/database"
	"github.com/dmarins/clean-arch-challenge-go/internal/infra/web"
	"github.com/dmarins/clean-arch-challenge-go/internal/usecase"
	"github.com/dmarins/clean-arch-challenge-go/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)

	return &usecase.CreateOrderUseCase{}
}

func NewListOrderUseCase(db *sql.DB) *usecase.ListOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrderUseCase,
	)

	return &usecase.ListOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)

	return &web.WebOrderHandler{}
}
