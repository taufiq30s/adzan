package utils

import (
	"fmt"
	"math"
	"time"
)

type TimeComponents struct {
	Hours   int16
	Minutes int16
	Seconds int16
}

func NewTimeComponents(d float64) (*TimeComponents, error) {
	if math.IsInf(d, 0) {
		return nil, fmt.Errorf("given value is infinite")
	}
	if math.IsNaN(d) {
		return nil, fmt.Errorf("given value is NaN")
	}

	hours := math.Floor(d)
	minutes := math.Floor((d - hours) * 60)
	seconds := math.Floor((d - (hours + minutes/60)) * 60 * 60)

	return &TimeComponents{
		Hours:   int16(hours),
		Minutes: int16(minutes),
		Seconds: int16(seconds),
	}, nil
}

func (t *TimeComponents) DateComponents(d *DateComponents) time.Time {
	date := time.Date(
		int(d.Year),
		time.Month(d.Month),
		int(d.Day),
		0,
		int(t.Minutes),
		int(t.Seconds),
		0,
		time.UTC,
	)
	return date.Add(time.Hour * time.Duration(t.Hours))
}
