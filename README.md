# Clean Architecture - FullCycle Go Expert Challenge

https://goexpert.fullcycle.com.br/pos-goexpert/

[![Go](https://img.shields.io/badge/go-1.22.4-informational?logo=go)](https://go.dev)

## Clone the project

```
$ git clone https://github.com/dmarins/clean-arch-challenge-go.git
$ cd clean-arch-challenge-go
```

## Download dependencies

```
$ go mod tidy
```

## Containers Up

```
$ make dc-up
```

## Database Init

```
$ make db-init
```

## Run Project

```
$ make run
```

## REST API Requests
```
HTTP requests inside /api directory
```

## GraphQL API Requests
```
access http://localhost:8080

Create Order:

mutation createOrder {
  createOrder(input: {id: "lalala", Price: 22.22, Tax: 33.33}) {
    id
    Price
    Tax
    FinalPrice
  }
}

List Orders:

query listOrders{
  listOrders {
    id
    Price
    Tax
    FinalPrice
  }
}

```

## GRPC Calls
```
use evans as client GRPC:
$ evans -r repl

$ package pb

$ service OrderService

$ call CreateOrder (provide the data)

$ call ListOrders

```