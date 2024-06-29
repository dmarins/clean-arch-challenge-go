package main

import (
	"database/sql"
	"fmt"

	"github.com/dmarins/clean-arch-challenge-go/configs"
	"github.com/dmarins/clean-arch-challenge-go/internal/event/handler"
	"github.com/dmarins/clean-arch-challenge-go/internal/infra/web/webserver"
	"github.com/dmarins/clean-arch-challenge-go/pkg/events"
	"github.com/streadway/amqp"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Envs
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// Db
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	rabbitMQChannel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	// createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)

	// HttpServer
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	fmt.Println("Starting web server on port", configs.WebServerPort)

	webserver.Start()
}
