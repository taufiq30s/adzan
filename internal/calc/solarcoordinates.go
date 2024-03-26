package calc

import (
	"math"

	"github.com/taufiq30s/adzan/internal/utils"
)

type SolarCoordinates struct {
	Declination          float64
	RightAscension       float64
	ApparentSiderealTime float64
}

func NewSolarCoordinates(jde float64) *SolarCoordinates {
	T := GetJulianCentury(jde)
	L0 := MeanSolarLongitude(T)
	Lp := MeanLunarLongitude(T)
	tetha := AscendingLunarNodeLongitude(T)
	lambda := utils.Radians(ApparentSolarLongitude(T, L0))
	theta0 := MeanSiderealTime(T)
	dPsi := NutationInLongitude(L0, Lp, tetha)
	dEpsilon := NutationInObliquity(L0, Lp, tetha)
	epsilon0 := MeanObliquityOfTheEcliptic(T)
	epsilonApp := utils.Radians(ApparentObliquityOfTheEcliptic(T, epsilon0))

	declination := utils.Degrees(math.Asin(math.Sin(epsilonApp) * math.Sin(lambda)))

	rightAscension := utils.UnwindAngle(
		utils.Degrees(math.Atan2(math.Cos(epsilonApp)*math.Sin(lambda), math.Cos(lambda))),
	)

	apparentSiderealTime := theta0 + (((dPsi * 3600) * math.Cos(utils.Radians(epsilon0+dEpsilon))) / 3600)

	return &SolarCoordinates{
		Declination:          declination,
		RightAscension:       rightAscension,
		ApparentSiderealTime: apparentSiderealTime,
	}
}
