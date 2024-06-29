package main

import "github.com/dmarins/clean-arch-challenge-go/configs"

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
}
