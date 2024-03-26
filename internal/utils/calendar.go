package utils

import (
	"math"
	"time"
)

func IsLeapYear(year int) bool {
	return year%4 == 0 && !(year%100 == 0 && year%400 != 0)
}

func RoundToNearestMinutes(d time.Time) time.Time {
	seconds := float64(d.Second()) / 60.0
	nearestMinute := d.Minute() + int(math.Round(seconds))
	return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), nearestMinute, 0, 0, d.Location())
}
