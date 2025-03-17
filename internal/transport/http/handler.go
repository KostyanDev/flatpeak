package http

import (
	_ "app/internal/domain"
	"app/internal/transport/converters"
	"encoding/json"
	"net/http"
)

// GetOptimalSlots
//
// @Summary Get Optimal Carbon Intensity Slots
// @Description Retrieves the best time slots with the lowest carbon intensity for the given duration.
// @Tags Carbon Intensity
// @Accept  json
// @Produce  json
// @Param duration query int false "Duration in minutes (default: 30, max: 1440)"
// @Param continuous query bool false "Whether to return a single continuous slot (default: false)"
// @Success 200 {array} converters.SlotsResponse "List of optimal time slots"
// @Failure 400 {object} map[string]string "Bad Request: Invalid query parameters"
// @Failure 500 {object} map[string]string "Internal Server Error: Failed to retrieve carbon intensity data"
// @Router /slots [get]
func (h *Handler) GetOptimalSlots(w http.ResponseWriter, r *http.Request) {
	filters, err := converters.ParseGetSlotsRequest(r.URL.Query())
	if err != nil {
		http.Error(w, "invalid query parameters", http.StatusBadRequest)
		return
	}

	slots, err := h.service.GetWeightedCarbonIntensity(r.Context(), filters)
	if err != nil {
		h.log.Error("Error retrieving optimal slots: ", err)
		http.Error(w, "Failed to retrieve carbon intensity data", http.StatusInternalServerError)
		return
	}

	responseSlots := converters.ToCarbonSlotResponseList(slots)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseSlots); err != nil {
		http.Error(w, "Failed to retrieve carbon intensity slots", http.StatusInternalServerError)
	}
}
