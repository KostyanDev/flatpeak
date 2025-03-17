package converter

import (
	"app/internal/domain"
	"errors"
	"time"
)

type CustomTime struct {
	time.Time
}

const iso8601Format = "2006-01-02T15:04Z"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	parsedTime, err := time.Parse(`"`+iso8601Format+`"`, string(b))
	if err != nil {
		return errors.New("invalid datetime format, expected ISO8601 like 2025-03-17T00:00Z")
	}
	ct.Time = parsedTime
	return nil
}

type CarbonIntensityResponse struct {
	Data []CarbonIntensityData `json:"data"`
}

type CarbonIntensityData struct {
	ValidFrom CustomTime    `json:"from"`
	ValidTo   CustomTime    `json:"to"`
	Intensity IntensityData `json:"intensity"`
}

type IntensityData struct {
	Forecast int    `json:"forecast"`
	Actual   int    `json:"actual"`
	Index    string `json:"index"`
}

func (c CarbonIntensityData) toDomain() domain.Carbon {
	return domain.Carbon{
		ValidFrom: c.ValidFrom.Time,
		ValidTo:   c.ValidTo.Time,
		Intensity: domain.IntensityData{
			Forecast: c.Intensity.Forecast,
			Actual:   c.Intensity.Actual,
			Index:    c.Intensity.Index,
		},
	}
}

func ToDomain(data CarbonIntensityResponse) []domain.Carbon {
	retVal := make([]domain.Carbon, len(data.Data))
	for i := range data.Data {
		retVal[i] = data.Data[i].toDomain()
	}
	return retVal
}
