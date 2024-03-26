package calc

import (
	"math"

	"github.com/taufiq30s/adzan/internal/utils"
)

type SolarTime struct {
	Transit float64
	Sunrise float64
	Sunset  float64

	Obsever            *utils.Coordinates
	Solar              *SolarCoordinates
	PrevSolar          *SolarCoordinates
	NextSolar          *SolarCoordinates
	ApproximateTransit float64
}

func NewSolarTime(date *utils.DateComponents, coordinate *utils.Coordinates) *SolarTime {
	jde := GetJulianDay(date, 0)

	solar := NewSolarCoordinates(jde)
	prevSolar := NewSolarCoordinates(jde - 1)
	nextSolar := NewSolarCoordinates(jde + 1)

	approximateTransit := ApproximateTransit(
		coordinate.Longitude, solar.ApparentSiderealTime, solar.RightAscension,
	)
	solarAltitude := -50.0 / 60.0

	transit := CorrectedTransit(
		approximateTransit, coordinate.Longitude, solar.ApparentSiderealTime,
		solar.RightAscension, prevSolar.RightAscension, nextSolar.RightAscension,
	)
	sunrise := CorrectedHourAngle(
		approximateTransit, solarAltitude, coordinate, false,
		solar.ApparentSiderealTime, solar.RightAscension, prevSolar.RightAscension,
		nextSolar.RightAscension, solar.Declination, prevSolar.Declination,
		nextSolar.Declination,
	)
	sunset := CorrectedHourAngle(
		approximateTransit, solarAltitude, coordinate, true,
		solar.ApparentSiderealTime, solar.RightAscension, prevSolar.RightAscension,
		nextSolar.RightAscension, solar.Declination, prevSolar.Declination,
		nextSolar.Declination,
	)

	return &SolarTime{
		Transit:            transit,
		Sunrise:            sunrise,
		Sunset:             sunset,
		Obsever:            coordinate,
		Solar:              solar,
		PrevSolar:          prevSolar,
		NextSolar:          nextSolar,
		ApproximateTransit: approximateTransit,
	}
}

func (solar *SolarTime) HourAngle(angle float64, afterTransit bool) float64 {
	return CorrectedHourAngle(
		solar.ApproximateTransit, angle, solar.Obsever, afterTransit,
		solar.Solar.ApparentSiderealTime, solar.Solar.RightAscension,
		solar.PrevSolar.RightAscension, solar.NextSolar.RightAscension,
		solar.Solar.Declination, solar.PrevSolar.Declination,
		solar.NextSolar.Declination,
	)
}

func (solar *SolarTime) Afternoon(sl utils.ShadowLength) float64 {
	tangent := math.Abs(solar.Obsever.Latitude - solar.Solar.Declination)
	inverse := utils.ShadowLengthToFloatMap[sl] + math.Tan(utils.Radians(tangent))
	angle := utils.Degrees(math.Atan(1.0 / inverse))

	return solar.HourAngle(angle, true)
}
