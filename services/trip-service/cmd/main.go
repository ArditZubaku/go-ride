package main

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

func main() {
	ctx := context.Background()

	inMemRepo := repository.NewInMemRepository()

	svc := service.NewService(inMemRepo)
	t, err := svc.CreateTrip(ctx, &domain.RideFareModel{UserID: "42"})
	if err != nil {
		log.Printf("ERROR: %+v", err)
	}

	log.Println(t)

	select {}
}
