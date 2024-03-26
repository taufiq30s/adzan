package utils

import "math"

// Normalizing an angle within a bound.
func NormalizeWithBound(angle float64, max float64) float64 {
	return angle - (max * math.Floor(angle/max))
}

// Normalizing an angle within a single revolution (360).
func UnwindAngle(angle float64) float64 {
	return NormalizeWithBound(angle, 360.0)
}

func ClosestAngle(angle float64) float64 {
	if angle >= -180 && angle <= 180 {
		return angle
	}
	return angle - (360.0 * math.Round(angle/360.0))
}

func Interpolate(y2 float64, y1 float64, y3 float64, n float64) float64 {
	a := y2 - y1
	b := y3 - y2
	c := b - a
	return y2 + (n/2)*(a+b+n*c)
}

func InterpolateAngles(y2 float64, y1 float64, y3 float64, n float64) float64 {
	a := UnwindAngle(y2 - y1)
	b := UnwindAngle(y3 - y2)
	c := b - a
	return y2 + (n/2)*(a+b+n*c)
}

// Convert Degree to Radian
func Radians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

// Convert Radian to Degree
func Degrees(radians float64) float64 {
	return radians * (180.0 / math.Pi)
}
