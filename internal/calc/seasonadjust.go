package calc

import (
	"math"
	"time"

	"github.com/taufiq30s/adzan/internal/utils"
)

type seasonAdjust struct {
	a       float64
	b       float64
	c       float64
	d       float64
	dyy     int
	sunTime time.Time
}

func DaysSinceSolstice(dayOfYear int, year int, latitude float64) int {
	daysSinceSolistice := 0
	northernOffset := 10
	isLeapYear := utils.IsLeapYear(year)
	southernOffset := 172
	daysInYear := 365
	if isLeapYear {
		southernOffset = 173
		daysInYear = 366
	}
	if latitude >= 0 {
		daysSinceSolistice = dayOfYear + northernOffset
		if daysSinceSolistice >= daysInYear {
			daysSinceSolistice -= daysInYear
		}
	} else {
		daysSinceSolistice = dayOfYear - southernOffset
		if daysSinceSolistice < 0 {
			daysSinceSolistice += daysInYear
		}
	}
	return daysSinceSolistice
}

func (seasonAdj *seasonAdjust) adjust() time.Time {
	adjustment := 0.0
	switch {
	case seasonAdj.dyy < 91:
		adjustment = seasonAdj.a + (seasonAdj.b-seasonAdj.a)/91*float64(seasonAdj.dyy)
	case seasonAdj.dyy < 137:
		adjustment = seasonAdj.b + (seasonAdj.c-seasonAdj.b)/46*(float64(seasonAdj.dyy)-91)
	case seasonAdj.dyy < 183:
		adjustment = seasonAdj.c + (seasonAdj.d-seasonAdj.c)/46*(float64(seasonAdj.dyy)-137)
	case seasonAdj.dyy < 229:
		adjustment = seasonAdj.d + (seasonAdj.c-seasonAdj.d)/46*(float64(seasonAdj.dyy)-183)
	case seasonAdj.dyy < 275:
		adjustment = seasonAdj.c + (seasonAdj.b-seasonAdj.c)/46*(float64(seasonAdj.dyy)-229)
	default:
		adjustment = seasonAdj.b + (seasonAdj.a-seasonAdj.b)/91*(float64(seasonAdj.dyy)-275)
	}
	return seasonAdj.sunTime.Add(time.Second * time.Duration(-1*math.Round(adjustment*60)))
}

func SeasonAdjustedMorningTwilight(latitude float64, day int, year int, sunrise time.Time) time.Time {
	a := 75 + ((28.65 / 55) * math.Abs(latitude))
	b := 75 + ((10.44 / 55) * math.Abs(latitude))
	c := 75 + ((19.44 / 55) * math.Abs(latitude))
	d := 75 + ((48.10 / 55) * math.Abs(latitude))
	dyy := DaysSinceSolstice(day, year, latitude)

	res := &seasonAdjust{
		a:       a,
		b:       b,
		c:       c,
		d:       d,
		dyy:     dyy,
		sunTime: sunrise,
	}
	return res.adjust()
}

func SeasonAdjustedEveningTwilight(latitude float64, day int, year int, sunset time.Time) time.Time {
	a := 75 + ((25.60 / 55) * math.Abs(latitude))
	b := 75 + ((2.05 / 55) * math.Abs(latitude))
	c := 75 + ((9.21 / 55) * math.Abs(latitude))
	d := 75 + ((6.14 / 55) * math.Abs(latitude))
	dyy := DaysSinceSolstice(day, year, latitude)

	res := &seasonAdjust{
		a:       a,
		b:       b,
		c:       c,
		d:       d,
		dyy:     dyy,
		sunTime: sunset,
	}
	return res.adjust()
}
