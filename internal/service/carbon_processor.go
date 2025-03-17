package service

import (
	"app/internal/domain"
	"context"
	"math"
	"sort"
	"time"
)

func (s Service) GetWeightedCarbonIntensity(ctx context.Context, filter domain.GetSlots) ([]domain.LowestCarbonPeriod, error) {
	data, err := s.carbonClient.FetchCarbonForecast(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to fetch carbon intensity forecast")
		return nil, err
	}

	return calculateWeightedAverage(data, filter), nil
}

func calculateWeightedAverage(data []domain.Carbon, filter domain.GetSlots) []domain.LowestCarbonPeriod {
	if len(data) == 0 {
		return nil
	}

	countOfSteps := int(math.Ceil(float64(filter.Duration) / 30.0))
	if countOfSteps > len(data) {
		countOfSteps = len(data)
	}

	if filter.Continuous == false {
		retVal := make([]domain.LowestCarbonPeriod, countOfSteps)
		sort.Slice(data, func(i, j int) bool {
			if data[i].Intensity.Forecast == data[j].Intensity.Forecast {
				return data[i].ValidFrom.Before(data[j].ValidFrom)
			}
			return data[i].Intensity.Forecast < data[j].Intensity.Forecast
		})
		for i := 0; i < countOfSteps; i++ {
			retVal[i] = domain.LowestCarbonPeriod{
				ValidFrom: data[i].ValidFrom,
				ValidTo:   data[i].ValidTo,
				Intensity: data[i].Intensity.Forecast,
			}
		}
		return retVal
	}

	return findOptimalContinuousPeriod(data, countOfSteps)
}

func calculateAverageIntensity(data []domain.Carbon) domain.LowestCarbonPeriod {
	if len(data) == 0 {
		return domain.LowestCarbonPeriod{}
	}

	var totalForecast, count int
	var from, to time.Time

	from = data[0].ValidFrom
	to = data[0].ValidTo

	for _, entry := range data {
		totalForecast += entry.Intensity.Forecast
		count++

		if entry.ValidFrom.Before(from) {
			from = entry.ValidFrom
		}
		if entry.ValidTo.After(to) {
			to = entry.ValidTo
		}
	}

	avgIntensity := totalForecast / count

	return domain.LowestCarbonPeriod{
		ValidFrom: from,
		ValidTo:   to,
		Intensity: avgIntensity,
	}
}

func findOptimalContinuousPeriod(data []domain.Carbon, countOfSteps int) []domain.LowestCarbonPeriod {
	if len(data) == 0 || countOfSteps > len(data) {
		return nil
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].ValidFrom.Before(data[j].ValidFrom)
	})

	totalIntensity := 0
	minAvgIntensity := math.MaxFloat64
	bestStartIndex := 0

	for i := 0; i < countOfSteps; i++ {
		totalIntensity += data[i].Intensity.Forecast
	}

	currentAvgIntensity := float64(totalIntensity) / float64(countOfSteps)
	minAvgIntensity = currentAvgIntensity

	for i := countOfSteps; i < len(data); i++ {
		totalIntensity -= data[i-countOfSteps].Intensity.Forecast
		totalIntensity += data[i].Intensity.Forecast

		currentAvgIntensity = float64(totalIntensity) / float64(countOfSteps)

		if currentAvgIntensity < minAvgIntensity {
			minAvgIntensity = currentAvgIntensity
			bestStartIndex = i - countOfSteps + 1
		}
	}

	bestSlots := data[bestStartIndex : bestStartIndex+countOfSteps]
	optimalPeriod := calculateAverageIntensity(bestSlots)

	return []domain.LowestCarbonPeriod{optimalPeriod}
}
