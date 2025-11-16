package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

const GRPC_ADDR = ":9083"

func main() {
	inMemRepo := repository.NewInMemRepository()
	_ = service.NewService(inMemRepo) // TODO: use svc when grpc handler is implemented

	listener, err := net.Listen("tcp", GRPC_ADDR)
	if err != nil {
		log.Fatalf("Failed to start gRPC listener on %s: %v", GRPC_ADDR, err)
	}

	grpcServer := grpc.NewServer()
	// TODO: init grpc handler impl

	log.Printf("Starting the gRPC server of Trip Service on addr: %s", listener.Addr().String())

	serverErrorsChan := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			serverErrorsChan <- err
		}
	}()

	handleShutdown(grpcServer, serverErrorsChan, shutdown)
}

func handleShutdown(
	server *grpc.Server,
	serverErrorsChan chan error,
	shutdown chan os.Signal,
) {
	select {
	case err := <-serverErrorsChan:
		log.Printf("Error starting the server: %v\n", err)
		log.Println("Shutting down the gRPC server...")
		// Stop receiving signals since we're shutting down
		signal.Stop(shutdown)
		close(serverErrorsChan)
		server.GracefulStop()
	case sig := <-shutdown:
		log.Printf("Server is shutting down due to: %v signal\n", sig)
		// Stop receiving signals since we're shutting down
		signal.Stop(shutdown)
		close(serverErrorsChan)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// An empty struct channel to signal when GracefulStop completes so we can log that
		done := make(chan struct{})
		go func() {
			log.Println("Shutting down the gRPC server gracefully...")
			server.GracefulStop()
			close(done)
		}()

		// Wait for graceful stop or timeout
		select {
		case <-done:
			log.Println("gRPC server shut down gracefully")
		case <-ctx.Done():
			log.Printf("Graceful shutdown timeout exceeded, forcing stop: %v\n", ctx.Err())
			server.Stop()
		}
	}
}
