package main

import (
	"encoding/json"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
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

	// TODO: Call trip service

	res := contracts.APIResponse{Data: "OK", Error: nil}

	writeJSON(w, http.StatusCreated, res)
}
