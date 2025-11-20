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
