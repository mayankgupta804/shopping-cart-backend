package main

import (
	"context"
	"log"
	"shopping-cart-backend/internal/repository"
	"shopping-cart-backend/internal/serializer"
	"shopping-cart-backend/internal/service"
	"shopping-cart-backend/pkg/database"
)

func main() {
	cfg := database.Config{
		Name:     "bike_station",
		User:     "mayank",
		Password: "secret",
		Host:     "localhost",
		Port:     "5432",
	}

	d, err := database.NewFromEnv(context.Background(), cfg.DatabaseConfig())
	if err != nil {
		log.Fatal(err)
	}
	a := repository.NewAccountRepository(d)
	s := service.NewAccountService(a)
	if err := s.Create(serializer.CreateAccountRequest{
		Name:     "awesome",
		Email:    "dasda@dasd.com",
		Password: "dasdasd",
		Role:     "user",
	}); err != nil {
		log.Fatal(err)
	}

}
