package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
)

type inMemRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInMemRepository() *inMemRepository {
	return &inMemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inMemRepository) CreateTrip(
	ctx context.Context,
	trip *domain.TripModel,
) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *inMemRepository) SaveRideFare(
	ctx context.Context,
	fare *domain.RideFareModel,
) error {
	r.rideFares[fare.ID.Hex()] = fare
	return nil
}

func (r *inMemRepository) GetRideFareByID(
	ctx context.Context,
	id string,
) (*domain.RideFareModel, error) {
	fare, ok := r.rideFares[id]
	if !ok {
		return nil, fmt.Errorf("ride fare with id %s doesn't exist!", id)
	}

	return fare, nil
}
