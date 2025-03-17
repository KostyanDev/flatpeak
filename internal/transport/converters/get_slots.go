package converters

import (
	"app/internal/domain"
	"errors"
	"net/url"
	"strconv"
	"time"
)

type SlotResponse struct {
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
	Intensity int       `json:"intensity"`
}

type SlotsResponse struct {
	Slots []SlotResponse `json:"slots"`
}

func ToCarbonSlotResponse(data domain.LowestCarbonPeriod) SlotResponse {
	return SlotResponse{
		ValidFrom: data.ValidFrom,
		ValidTo:   data.ValidTo,
		Intensity: data.Intensity,
	}
}

func ToCarbonSlotResponseList(slots []domain.LowestCarbonPeriod) SlotsResponse {

	responseList := make([]SlotResponse, len(slots))
	for i := range slots {
		responseList[i] = ToCarbonSlotResponse(slots[i])
	}

	return SlotsResponse{Slots: responseList}
}

func ParseGetSlotsRequest(query url.Values) (domain.GetSlots, error) {
	var request domain.GetSlots

	durationStr := query.Get("duration")
	if durationStr == "" {
		return request, errors.New("missing required parameter: duration")
	}
	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration < 1 || duration > 1440 {
		return request, errors.New("invalid value for parameter: duration (must be between 1 and 1440 minutes)")
	}
	request.Duration = duration

	continuousStr := query.Get("continuous")
	if continuousStr == "" {
		return request, errors.New("missing required parameter: continuous")
	}
	continuous, err := strconv.ParseBool(continuousStr)
	if err != nil {
		return request, errors.New("invalid value for parameter: continuous (must be 'true' or 'false')")
	}
	request.Continuous = continuous

	return request, nil
}
