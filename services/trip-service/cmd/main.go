package main

import (
	"context"
	"log"
	"net"
	"net/http"
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
	svc := service.NewService(inMemRepo)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
		<-shutdown
		cancel()
	}()

	listener, err := net.Listen("tcp", GRPC_ADDR)
	if err != nil {
		log.Fatalf("Failed to start gRPC listener on %s: %v", GRPC_ADDR, err)
	}

	grpcServer := grpc.NewServer()
	// TODO: init grpc handler impl

	log.Printf("Starting the gRPC server of Trip Service on addr: %s", listener.Addr().String())

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC server encountered an error while serving: %v", err)
			cancel()
		}
	}()

	// Wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shutting down the gRPC server...")
	grpcServer.GracefulStop()
}

func handleShutdown(
	server *http.Server,
	serverErrorsChan chan error,
	shutdown chan os.Signal,
) {
	select {
	case err := <-serverErrorsChan:
		log.Printf("Error starting the server: %v\n", err)
	case sig := <-shutdown:
		log.Printf("Server is shutting down due to: %v signal\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil { // Waits until all handlers finish or cancels based on the ctx
			log.Printf("Could not shutdown the server gracefully because: %v\n", err)
			if err := server.Close(); err != nil {
				log.Printf("Could not close the server connections because: %v\n", err)
			}
		}
	}
}
