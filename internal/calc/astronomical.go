package calc

import (
	"math"

	"github.com/taufiq30s/adzan/internal/utils"
)

// Calculate Mean Solar Longitude of Sun in degree
// where t is julian century
// The accuracy of this method is low (accuracy 0.01 degree)
//
// Reference: Astronomical Algorithm Chapter 25 Page 163
func MeanSolarLongitude(t float64) float64 {
	L0 := float64(280.4664567) + (float64(36000.76983) * t) + (float64(0.0003032) * math.Pow(t, 2))
	return utils.UnwindAngle(L0)
}

// Calculate Mean Lunar
//
// Reference: Astronomical Algorithm Chapter 22 Page 144
func MeanLunarLongitude(t float64) float64 {
	Lp := float64(218.3165) + (float64(481267.8813) * t)
	return utils.UnwindAngle(Lp)
}

// Ascending Lunar Node Longitude
//
// Reference: Chapter 22 Page 144
func AscendingLunarNodeLongitude(t float64) float64 {
	Ω := float64(125.04452) -
		(float64(1934.136261) * t) +
		(float64(0.0020708) * math.Pow(t, 2)) +
		(math.Pow(t, 3) / float64(450000))
	return utils.UnwindAngle(Ω)
}

// The Sun's equation of the center in degrees.
//
// Where t is Julian Century and m is mean
// annomaly of the sun in degree.
//
// Reference: Chapter 25 Page 164
func SolarEquitionOfTheCenter(t float64, m float64) float64 {
	mRad := utils.Radians(m)
	return ((1.914602 - (0.004817 * t) - (0.000014 * math.Pow(t, 2))) * math.Sin(mRad)) +
		((0.019993 - (0.000101 * t)) * math.Sin(2*mRad)) +
		(0.000289 * math.Sin(3*mRad))
}

// The Mean anomaly of the Sun
//
// Where 't' is Julian Century
//
// Reference: Chapter 25 Page 163
func MeanSolarAnomaly(t float64) float64 {
	M := 357.52911 + (35999.05029 * t) - (0.0001537 * math.Pow(t, 2))
	return utils.UnwindAngle(M)
}

// The Apperent Longitude of the Sun
// Referred to the true equinox of the date.
//
// Where 't' is Julian Century and 'L0' is mean
// longitude of the Sun.
//
// Reference: Chapter 25 Page 164
func ApparentSolarLongitude(t float64, L0 float64) float64 {
	longitude := L0 + SolarEquitionOfTheCenter(t, MeanSolarAnomaly(t))
	tetha := 125.04 - (1934.136 * t)
	lambda := longitude - 0.00569 - (0.00478 * math.Sin(utils.Radians(tetha)))
	return utils.UnwindAngle(lambda)
}

// Mean Sidereal Time.
// The hour angle of the vernal equinox, in degrees
//
// Reference: Chapter 25 Page 165
func MeanSiderealTime(t float64) float64 {
	JD := (t * 36525) + 2451545.0
	tetha := 280.46061837 +
		(360.98564736629 * (JD - 2451545)) +
		(0.000387933 * math.Pow(t, 2)) -
		(math.Pow(t, 3) / 38710000)
	return utils.UnwindAngle(tetha)
}

// Nutation in Longitude
// returns the nutation in longitude
//
// Given 'L0', the solar longitude.
// 'Lp', the lunar longitude.
// 'tetha', the ascending node
//
// Reference: Chapter 22 Page 144
func NutationInLongitude(L0 float64, Lp float64, omega float64) float64 {
	return ((-17.2 / 3600.0) * math.Sin(utils.Radians(omega))) -
		((1.32 / 3600.0) * math.Sin(2*utils.Radians(L0))) -
		((0.23 / 3600.0) * math.Sin(2*utils.Radians(Lp))) +
		((0.21 / 3600.0) * math.Sin(2*utils.Radians(omega)))
}

// Nutation in Obliquity
// returns the nutation in obliquity
//
// Given 'L0', the solar longitude.
// 'Lp', the lunar longitude.
// 'tetha', the ascending node
//
// Reference: Chapter 22 Page 144
func NutationInObliquity(L0 float64, Lp float64, tetha float64) float64 {
	return ((9.2 / 3600) * math.Cos(utils.Radians(tetha))) +
		((0.57 / 3600) * math.Cos(2*utils.Radians(L0))) +
		((0.10 / 3600) * math.Cos(2*utils.Radians(Lp))) -
		((0.09 / 3600) * math.Cos(2*utils.Radians(tetha)))
}

// Mean Obliquity of The Ecliptic
// returns the mean obliquity of the ecliptic in degrees
//
// Given 't', the julian century
//
// Referenc: Chapter 22 Page 147
func MeanObliquityOfTheEcliptic(t float64) float64 {
	return 23.439291 -
		(0.013004167 * t) -
		(0.0000001639 * math.Pow(t, 2)) +
		(0.0000005036 * math.Pow(t, 3))
}

// Apparent Obliquity of The Ecliptic
// returns the mean obliquity of the ecliptic,
// corrected for calculating the apparent position
// of the sun in degrees
//
// Given 't', the julian century and
// 'meanObliquityEcliptic', mean obliquity of the ecliptic
//
// Reference: Chapter 25 Page 165
func ApparentObliquityOfTheEcliptic(t float64, meanObliquityEcliptic float64) float64 {
	O := 125.04 - (1934.136 * t)
	return meanObliquityEcliptic + (0.00256 * math.Cos(utils.Radians(O)))
}

// Altitude of Celestial Body
// returns the altitude of the celestial body
func AltitudeOfCelestialBody(observerLatitude float64, declination float64, H float64) float64 {
	return utils.Degrees(
		math.Asin(
			(math.Sin(utils.Radians(observerLatitude)) * math.Sin(utils.Radians(declination))) +
				(math.Cos(utils.Radians(observerLatitude)) * math.Cos(utils.Radians(declination)) *
					math.Cos(utils.Radians(H))),
		),
	)
}

// Approximate Transit
// returns the approximate transite
func ApproximateTransit(L float64, siderealTime float64, rightAscension float64) float64 {
	Lw := L * -1
	return utils.NormalizeWithBound((rightAscension+Lw-siderealTime)/360, 1)
}

func CorrectedTransit(approximateTransit float64, L float64, siderealTime float64,
	rightAscension float64, prevRightAscension float64,
	nextRightAscension float64) float64 {
	Lw := L * -1
	theta := utils.UnwindAngle(siderealTime + (360.985647 * approximateTransit))
	a := utils.UnwindAngle(utils.InterpolateAngles(rightAscension, prevRightAscension, nextRightAscension, approximateTransit))
	H := utils.ClosestAngle(theta - Lw - a)
	dm := H / -360
	return (approximateTransit + dm) * 24
}

func CorrectedHourAngle(
	approximateTransit float64,
	angle float64,
	coordinate *utils.Coordinates,
	afterTransit bool,
	siderealTime float64,
	rightAscension float64,
	prevRightAscension float64,
	nextRightAscension float64,
	declination float64,
	prevDeclination float64,
	nextDeclination float64,
) float64 {
	Lw := coordinate.Longitude * -1
	term1 := math.Sin(utils.Radians(angle)) - (math.Sin(utils.Radians(coordinate.Latitude)) * math.Sin(utils.Radians(declination)))
	term2 := math.Cos(utils.Radians(coordinate.Latitude)) * math.Cos(utils.Radians(declination))
	H0 := utils.Degrees(math.Acos(term1 / term2))
	m := approximateTransit + (H0 / 360)
	if !afterTransit {
		m = approximateTransit - (H0 / 360)
	}
	theta := utils.UnwindAngle(siderealTime + (360.985647 * m))
	a := utils.UnwindAngle(utils.InterpolateAngles(rightAscension, prevRightAscension, nextRightAscension, m))
	delta := utils.Interpolate(declination, prevDeclination, nextDeclination, m)
	H := theta - Lw - a
	h := AltitudeOfCelestialBody(
		coordinate.Latitude,
		delta,
		H,
	)
	dm := (h - angle) / (360 * math.Cos(utils.Radians(delta)) *
		math.Cos(utils.Radians(coordinate.Latitude)) *
		math.Sin(utils.Radians(H)))
	return (m + dm) * 24
}
