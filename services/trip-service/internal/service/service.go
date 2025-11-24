package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/env"
	"ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(
	ctx context.Context,
	fare *domain.RideFareModel,
) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
		Driver:   &trip.TripDriver{},
	}

	return s.repo.CreateTrip(ctx, t)
}

func (s *service) GetRoute(
	ctx context.Context,
	pickup,
	destination *types.Coordinate,
) (*tripTypes.OsrmAPIResponse, error) {
	baseURL := env.GetString("OSRM_API", "http://router.project-osrm.org")

	url := fmt.Sprintf(
		"%s/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		baseURL,
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	routeRes := new(tripTypes.OsrmAPIResponse)
	if err := json.Unmarshal(body, routeRes); err != nil {
		return nil, fmt.Errorf("failed to parse route response: %v", err)
	}

	return routeRes, nil
}

func (s *service) EstimaPkgsPriceWithRoute(
	route *tripTypes.OsrmAPIResponse,
) []*domain.RideFareModel {
	baseFares := s.getBaseFares()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))

	for idx, fare := range baseFares {
		estimatedFares[idx] = s.estimateFareRoute(fare, route)
	}

	return estimatedFares
}

func (s *service) GenerateTripFares(
	ctx context.Context,
	rideFares []*domain.RideFareModel,
	userID string,
	route *tripTypes.OsrmAPIResponse,
) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))

	for idx, fare := range rideFares {
		fare := &domain.RideFareModel{
			ID:                primitive.NewObjectID(),
			UserID:            userID,
			TotalPriceInCents: fare.TotalPriceInCents,
			PackageSlug:       fare.PackageSlug,
			Route:             route,
		}

		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare: %w", err)
		}

		fares[idx] = fare
	}

	return fares, nil
}

func (s *service) GetFare(
	ctx context.Context,
	fareID string,
) (*domain.RideFareModel, error) {
	fare, err := s.repo.GetRideFareByID(ctx, fareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip fare: %w", err)
	}

	return fare, nil
}

func (s *service) ValidateFare(
	fare *domain.RideFareModel,
	userID string,
) (*domain.RideFareModel, error) {
	if fare.UserID != userID {
		return nil, fmt.Errorf("fare does not belong to the user")
	}
	return fare, nil
}

func (s *service) estimateFareRoute(
	fare *domain.RideFareModel,
	route *tripTypes.OsrmAPIResponse,
) *domain.RideFareModel {
	pricingCfg := tripTypes.GetDefaultPricingConfig()

	carPkgPrice := fare.TotalPriceInCents

	distanceKm := route.Routes[0].Distance
	durationMin := route.Routes[0].Duration

	// distance
	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	// time
	timeFare := durationMin * pricingCfg.PricingPerMinute
	// car price
	totalPrice := carPkgPrice + distanceFare + timeFare

	// return &domain.RideFareModel{
	// 	TotalPriceInCents: totalPrice,
	// 	PackageSlug:       fare.PackageSlug,
	// }

	fare.TotalPriceInCents = totalPrice

	return fare
}

func (s *service) getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{PackageSlug: "suv", TotalPriceInCents: 200},
		{PackageSlug: "sedan", TotalPriceInCents: 350},
		{PackageSlug: "van", TotalPriceInCents: 400},
		{PackageSlug: "luxury", TotalPriceInCents: 1_000},
	}
}
