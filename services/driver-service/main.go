package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"

	"google.golang.org/grpc"
)

var GRPCAddr = ":9082"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc := newService()

	rabbitMQURI := env.GetString(env.RabbitMQ.URI, env.RabbitMQDefaults.URI)
	rabbitMQ, err := messaging.NewRabbitMQ(rabbitMQURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQ.Close()

	log.Println("Successfully connected to RabbitMQ")

	consumer := NewTripConsumer(rabbitMQ)
	go func(consumer Consumer) {
		if err := consumer.Listen(); err != nil {
			log.Fatalf("Failed to listen to the RabbitMQ messages: %v", err)
		}
	}(consumer)

	// Starting the gRPC server
	grpcServer := grpc.NewServer()
	NewGrpcHandler(grpcServer, svc)

	log.Printf("Starting Driver service gRPC server on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
