package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(9 * time.Second)
	// We let the mux handle this
	// if r.Method != http.MethodPost {
	// 	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	// 	return
	// }
	// var reqBody previewTripRequest
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

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, "Failed to marshal request body to JSON", http.StatusInternalServerError)
		return
	}
	reader := bytes.NewReader(jsonBody)

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	// tripService.Client.PreviewTrip()

	resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	if err != nil {
		log.Printf("Error contacting trip service: %v", err)
		return
	}

	defer resp.Body.Close()

	var resBody any
	if err := json.NewDecoder(resp.Body).Decode(resBody); err != nil {
		http.Error(w, "Failed to parse JSON data from trip service", http.StatusBadRequest)
		return
	}

	res := contracts.APIResponse{Data: "OK", Error: nil}

	writeJSON(w, http.StatusCreated, res)
}
