package grpc

import (
	"context"
	"fmt"
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
	fmt.Printf("ROUTE: %+v", route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get route: %v", err)
	}

	return &pb.PreviewTripRes{
		Route:     route.ToProto(),
		TripID:    "",
		RideFares: make([]*pb.RideFare, 0),
	}, nil
}
