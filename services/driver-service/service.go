package main

import (
	math "math/rand/v2"
	"sync"

	pb "ride-sharing/shared/proto/driver"
	"ride-sharing/shared/util"

	"github.com/mmcloughlin/geohash"
)

type driverInMap struct {
	Driver *pb.Driver
	// Index int
	// TODO: route
}

type Service struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

func newService() *Service {
	return &Service{
		drivers: make([]*driverInMap, 0),
	}
}

func (s *Service) RegisterDriver(driverID string, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := math.IntN(len(PredefinedRoutes))
	randomRoute := PredefinedRoutes[randomIndex]

	randomPlate := GenerateRandomPlate()
	randomAvatar := util.GetRandomAvatar(randomIndex)

	// we can ignore this property for now, but it must be sent to the frontend.
	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &pb.Driver{
		Id:             driverID,
		Geohash:        geohash,
		Location:       &pb.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
		Name:           "Lando Norris",
		PackageSlug:    packageSlug,
		ProfilePicture: randomAvatar,
		CarPlate:       randomPlate,
	}

	s.drivers = append(s.drivers, &driverInMap{Driver: driver})

	return driver, nil
}

func (s *Service) UnregisterDriver(driverID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := s.drivers[:0] // reuse underlying array

	for _, d := range s.drivers {
		if d.Driver.Id != driverID {
			filtered = append(filtered, d)
		}
	}

	s.drivers = filtered
}
