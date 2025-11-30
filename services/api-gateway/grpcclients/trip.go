package grpcclients

import (
	"ride-sharing/shared/env"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type tripServiceClient struct {
	Client pb.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient() (*tripServiceClient, error) {
	tripServiceURL := env.GetString("TRIP_SERVICE_URL", "trip-service:9083")

	conn, err := grpc.NewClient(
		tripServiceURL,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewTripServiceClient(conn)

	return &tripServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *tripServiceClient) Close() {
	if c.conn != nil {
		util.CloseAndLog(c.conn, "gRPC connection for trip service client")
	}
}
