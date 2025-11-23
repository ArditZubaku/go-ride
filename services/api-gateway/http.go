package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	reqBody := new(previewTripRequest)
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validation
	if len(reqBody.UserID) <= 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// TODO: This can be done better - don't create a new connection for each req
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.toProto())
	if err != nil {
		errMsg := "Failed to preview the trip"
		log.Printf("%s: %v", errMsg, err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	res := contracts.APIResponse{Data: tripPreview, Error: nil}

	writeJSON(w, http.StatusCreated, res)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	reqBody := new(startTripRequest)

	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Why we need to create a new client for each connection:
	// because if a service is down, we don't want to block the whole application
	// so we create a new client for each connection
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(r.Context(), reqBody.toProto())
	if err != nil {
		errMsg := "Failed to start the trip"
		log.Printf("%s: %v", errMsg, err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: trip, Error: nil}

	writeJSON(w, http.StatusCreated, response)
}
