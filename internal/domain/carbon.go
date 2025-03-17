package domain

import "time"

type LowestCarbonPeriod struct {
	ValidFrom time.Time
	ValidTo   time.Time
	Intensity int
}
type Carbon struct {
	ValidFrom time.Time
	ValidTo   time.Time
	Intensity IntensityData
}

type IntensityData struct {
	Forecast int
	Actual   int
	Index    string
}
