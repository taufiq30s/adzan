package utils

import "time"

type DateComponents struct {
	Year  int16
	Month int8
	Day   int8
}

func NewDateComponents(date time.Time) *DateComponents {
	return &DateComponents{
		Year:  int16(date.Year()),
		Month: int8(date.Month()),
		Day:   int8(date.Day()),
	}
}

func (date *DateComponents) ConvertToTime() time.Time {
	return time.Date(int(date.Year), time.Month(date.Month), int(date.Day), 0, 0, 0, 0, time.UTC)
}
