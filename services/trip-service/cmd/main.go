package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"
)

const HTTP_ADDR = ":8083"

func main() {
	inMemRepo := repository.NewInMemRepository()
	svc := service.NewService(inMemRepo)
	mux := http.NewServeMux()

	httpHandler := h.HttpHandler{Service: svc}

	mux.HandleFunc("POST /preview", httpHandler.HandleTripPreview)

	server := &http.Server{
		Addr:    HTTP_ADDR,
		Handler: mux,
	}
	serverErrorsChan := make(chan error, 1)

	go func() {
		log.Printf("Server listening on: %s", HTTP_ADDR)
		serverErrorsChan <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	handleShutdown(server, serverErrorsChan, shutdown)
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
