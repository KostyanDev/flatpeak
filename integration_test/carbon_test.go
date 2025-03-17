package integration_test

import (
	"app/internal/domain"
	"app/internal/transport/converters"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/mock"
)

func (s *TestSuite) TestGetCarbonIntensity() {

	mockResponse := []domain.Carbon{
		{
			ValidFrom: parseTime("2025-03-17T17:00:00Z"),
			ValidTo:   parseTime("2025-03-17T17:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 120},
		},
		{
			ValidFrom: parseTime("2025-03-17T17:30:00Z"),
			ValidTo:   parseTime("2025-03-17T18:00:00Z"),
			Intensity: domain.IntensityData{Forecast: 110},
		},
		{
			ValidFrom: parseTime("2025-03-17T18:30:00Z"),
			ValidTo:   parseTime("2025-03-17T19:00:00Z"),
			Intensity: domain.IntensityData{Forecast: 90},
		},
		{
			ValidFrom: parseTime("2025-03-17T19:00:00Z"),
			ValidTo:   parseTime("2025-03-17T19:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 80}, // +
		},
		{
			ValidFrom: parseTime("2025-03-17T19:30:00Z"),
			ValidTo:   parseTime("2025-03-17T20:00:00Z"),
			Intensity: domain.IntensityData{Forecast: 85}, // +
		},
		{
			ValidFrom: parseTime("2025-03-17T20:00:00Z"),
			ValidTo:   parseTime("2025-03-17T20:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 75}, // +
		},
		{
			ValidFrom: parseTime("2025-03-17T20:30:00Z"),
			ValidTo:   parseTime("2025-03-17T21:00:00Z"),
			Intensity: domain.IntensityData{Forecast: 80}, // +
		},
		{
			ValidFrom: parseTime("2025-03-17T21:00:00Z"),
			ValidTo:   parseTime("2025-03-17T21:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 100},
		},
		{
			ValidFrom: parseTime("2025-03-17T21:30:00Z"),
			ValidTo:   parseTime("2025-03-17T22:00:00Z"),
			Intensity: domain.IntensityData{Forecast: 95},
		},
	}

	s.mockClient.EXPECT().FetchCarbonForecast(mock.Anything).Return(mockResponse, nil)

	req, err := http.NewRequest(http.MethodGet, s.server.URL+"/slots?duration=120&continuous=true", nil)
	s.Require().NoError(err)

	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, req)

	s.Equal(http.StatusOK, recorder.Code, "Expected HTTP 200 response")

	var response converters.SlotsResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	s.Require().NoError(err, "Failed to decode response body")

	s.Require().Greater(len(response.Slots), 0, "Expected at least one slot in the response")

	expectedFrom := parseTime("2025-03-17T19:00:00Z")
	expectedTo := parseTime("2025-03-17T21:00:00Z")
	expectedIntensity := 80 // 80, 85, 75, 80

	s.Equal(expectedFrom, response.Slots[0].ValidFrom, "ValidFrom does not match expected optimal slot")
	s.Equal(expectedTo, response.Slots[0].ValidTo, "ValidTo does not match expected optimal slot")
	s.Equal(expectedIntensity, response.Slots[0].Intensity, "Intensity average does not match expected value")
}

func (s *TestSuite) TestGetCarbonIntensity_NonContinuous() {

	mockResponse := []domain.Carbon{
		{
			ValidFrom: parseTime("2025-03-17T10:00:00Z"),
			ValidTo:   parseTime("2025-03-17T10:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 120},
		},
		{
			ValidFrom: parseTime("2025-03-17T10:30:00Z"),
			ValidTo:   parseTime("2025-03-17T11:00:00Z"),
			Intensity: domain.IntensityData{Forecast: 110},
		},
		{
			ValidFrom: parseTime("2025-03-17T12:00:00Z"),
			ValidTo:   parseTime("2025-03-17T12:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 70}, // +
		},
		{
			ValidFrom: parseTime("2025-03-17T15:00:00Z"),
			ValidTo:   parseTime("2025-03-17T15:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 65}, // +
		},
		{
			ValidFrom: parseTime("2025-03-17T19:00:00Z"),
			ValidTo:   parseTime("2025-03-17T19:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 60}, // +
		},

		{
			ValidFrom: parseTime("2025-03-17T20:00:00Z"),
			ValidTo:   parseTime("2025-03-17T20:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 90},
		},
		{
			ValidFrom: parseTime("2025-03-17T22:00:00Z"),
			ValidTo:   parseTime("2025-03-17T22:30:00Z"),
			Intensity: domain.IntensityData{Forecast: 100},
		},
	}

	s.mockClient.EXPECT().FetchCarbonForecast(mock.Anything).Return(mockResponse, nil)

	req, err := http.NewRequest(http.MethodGet, s.server.URL+"/slots?duration=90&continuous=false", nil)
	s.Require().NoError(err)

	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, req)

	s.Equal(http.StatusOK, recorder.Code, "Expected HTTP 200 response")

	var response converters.SlotsResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	s.Require().NoError(err, "Failed to decode response body")

	s.Require().Equal(3, len(response.Slots), "Expected exactly 3 optimal slots")

	expectedSlots := []domain.LowestCarbonPeriod{
		{
			ValidFrom: parseTime("2025-03-17T19:00:00Z"),
			ValidTo:   parseTime("2025-03-17T19:30:00Z"),
			Intensity: 60,
		},
		{
			ValidFrom: parseTime("2025-03-17T15:00:00Z"),
			ValidTo:   parseTime("2025-03-17T15:30:00Z"),
			Intensity: 65,
		},
		{
			ValidFrom: parseTime("2025-03-17T12:00:00Z"),
			ValidTo:   parseTime("2025-03-17T12:30:00Z"),
			Intensity: 70,
		},
	}

	for i, expected := range expectedSlots {
		s.Equal(expected.ValidFrom, response.Slots[i].ValidFrom, "ValidFrom mismatch")
		s.Equal(expected.ValidTo, response.Slots[i].ValidTo, "ValidTo mismatch")
		s.Equal(expected.Intensity, response.Slots[i].Intensity, "Intensity mismatch")
	}
}

func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse time: %s, error: %v", timeStr, err))
	}
	return t
}
