package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dmarins/clean-arch-challenge-go/configs"
	"github.com/dmarins/clean-arch-challenge-go/internal/event/handler"
	"github.com/dmarins/clean-arch-challenge-go/internal/infra/graph"
	"github.com/dmarins/clean-arch-challenge-go/internal/infra/grpc/pb"
	"github.com/dmarins/clean-arch-challenge-go/internal/infra/grpc/service"
	"github.com/dmarins/clean-arch-challenge-go/internal/infra/web/webserver"
	"github.com/dmarins/clean-arch-challenge-go/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db)

	// Http Server
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("post", "/order", webOrderHandler.Create)
	webserver.AddHandler("get", "/order", webOrderHandler.List)
	fmt.Println("Starting HTTP server on port", configs.WebServerPort)

	go webserver.Start()

	// GRPC Server
	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting GRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// GraphQL Server
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}
