package utils

import "fmt"

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func NewCoordinates(latitude float64, longitude float64) (*Coordinates, error) {
	if latitude > 90 || latitude < -90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	} else if longitude > 180 || longitude < -180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}

	return &Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}
