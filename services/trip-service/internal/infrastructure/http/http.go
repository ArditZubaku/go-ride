package http

import (
	"encoding/json"
	"log"
	"net/http"

	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
)

type HttpHandler struct {
	Service domain.TripService
}

type previewTripRequest struct {
	UserID      string            `json:"userID"`
	PickUp      *types.Coordinate `json:"pickup"`
	Destination *types.Coordinate `json:"destination"`
}

func (h *HttpHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	reqBody := new(previewTripRequest)
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	t, err := h.Service.GetRoute(ctx, reqBody.PickUp, reqBody.Destination)
	if err != nil {
		log.Printf("ERROR HandleTripPreview: %+v", err)
	}

	writeJSON(w, http.StatusOK, t)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
