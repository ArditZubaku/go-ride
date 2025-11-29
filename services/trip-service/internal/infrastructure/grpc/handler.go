package grpc

import (
	"context"

	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	pb.UnimplementedTripServiceServer

	service domain.TripService
}

func NewHandler(server *grpc.Server, service domain.TripService) *handler {
	// handler := new(handler)
	// handler.service = service
	handler := &handler{
		service: service,
	}

	// This way gRPC is going to be able to call the handler's methods
	pb.RegisterTripServiceServer(server, handler)

	return handler
}

func (h *handler) PreviewTrip(
	ctx context.Context,
	req *pb.PreviewTripReq,
) (*pb.PreviewTripRes, error) {
	pickup := &types.Coordinate{
		Latitude:  req.StartLocation.Latitude,
		Longitude: req.StartLocation.Longitude,
	}

	destination := &types.Coordinate{
		Latitude:  req.EndLocation.Latitude,
		Longitude: req.EndLocation.Longitude,
	}

	route, err := h.service.GetRoute(ctx, pickup, destination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get route: %v", err)
	}

	// Estimate the ride fares prices based on the route (ex. distance)
	estimatedFares := h.service.EstimaPkgsPriceWithRoute(route)
	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, req.UserID, route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate the ride fares: %v", err)
	}

	return &pb.PreviewTripRes{
		Route:     route.ToProto(),
		RideFares: domain.RideFareModelsToProtos(fares),
	}, nil
}

func (h *handler) CreateTrip(ctx context.Context, req *pb.CreateTripReq) (*pb.CreateTripRes, error) {
	fareID := req.GetRideFareID()
	userID := req.GetUserID()

	fare, err := h.service.GetFare(ctx, fareID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getFareErr: %v", err.Error())
	}

	rightFare, err := h.service.ValidateFare(fare, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "validateFareErr: %v", err.Error())
	}

	trip, err := h.service.CreateTrip(ctx, rightFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trip: %v", err)
	}

	return &pb.CreateTripRes{TripID: trip.ID.Hex()}, nil
}
