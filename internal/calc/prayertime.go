package calc

import (
	"time"

	"github.com/taufiq30s/adzan/internal/utils"
)

type PrayerTimes struct {
	Imsak             time.Time
	Fajr              time.Time
	Sunrise           time.Time
	Dhuhr             time.Time
	Ashr              time.Time
	Magrib            time.Time
	Isha              time.Time
	Coordinates       *utils.Coordinates
	DateComponent     *utils.DateComponents
	CalculationParams *CalculationParameters
}

func createDateComponents(d float64, date *utils.DateComponents) (time.Time, error) {
	timeComponents, err := utils.NewTimeComponents(d)
	if err != nil {
		return time.Time{}, err
	}
	return timeComponents.DateComponents(date), nil
}

func NewPrayerTimes(coords *utils.Coordinates, date *utils.DateComponents, params *CalculationParameters) (*PrayerTimes, error) {
	currentDate := date.ConvertToTime()

	tommorowDate := utils.NewDateComponents(currentDate.AddDate(0, 0, 1))

	solarTime := NewSolarTime(date, coords)

	tempDhuhr, err := createDateComponents(solarTime.Transit, date)
	if err != nil {
		return nil, err
	}

	tempSunrise, err := createDateComponents(solarTime.Sunrise, date)
	if err != nil {
		return nil, err
	}

	tempMaghrib, err := createDateComponents(solarTime.Sunset, date)
	if err != nil {
		return nil, err
	}

	tommorowSolarTime := NewSolarTime(tommorowDate, coords)
	tommorowSunriseComponents, err := utils.NewTimeComponents(tommorowSolarTime.Sunrise)
	if err != nil {
		return nil, err
	}

	tempAshr, err := createDateComponents(solarTime.Afternoon(MazhabToShadowLengthMap[params.Mazhab]), date)
	if err != nil {
		return nil, err
	}
	tommorowSunrise := tommorowSunriseComponents.DateComponents(tommorowDate)
	night := tommorowSunrise.Sub(tempMaghrib) * 1000

	// Fajr Calculation
	tempFajr, err := createDateComponents(solarTime.HourAngle(-1*float64(params.FajrAngle), false), date)
	if err != nil {
		return nil, err
	}
	nightPortion, err := params.GetNightPortion()
	if err != nil {
		return nil, err
	}
	portion := nightPortion.fajr
	nightFraction := (int64)(portion * float32(night.Seconds()) / 1000)
	safeFajr := tempSunrise.Add(time.Second * time.Duration(-1*(int)(nightFraction)))

	if tempFajr.IsZero() || tempFajr.Before(safeFajr) {
		tempFajr = safeFajr
	}

	// Isha Calculation with check againts safe value
	var tempIsha time.Time
	if params.IshaInterval > 0 {
		tempIsha = tempMaghrib.Add(time.Second * time.Duration(params.IshaInterval*60))
	} else {
		tempIsha, err = createDateComponents(solarTime.HourAngle(-1*float64(params.IshaAngle), true), date)
		if err != nil {
			return nil, err
		}

		portion = nightPortion.Isha
		nightFraction = int64(portion * float32(night.Seconds()) / 1000)
		safeIsha := tempMaghrib.Add(time.Second * time.Duration(nightFraction))

		if tempIsha.IsZero() || tempIsha.After(safeIsha) {
			tempIsha = safeIsha
		}
	}

	// Assign final times to public struct members with all offsets
	fajr := utils.RoundToNearestMinutes(
		tempFajr.Add(
			time.Minute * time.Duration(params.Ajustment.Fajr+params.MethodAjustment.Fajr),
		),
	)
	sunrise := utils.RoundToNearestMinutes(
		tempSunrise.Add(
			time.Minute * time.Duration(params.Ajustment.Sunrise+params.MethodAjustment.Sunrise),
		),
	)
	dhuhr := utils.RoundToNearestMinutes(
		tempDhuhr.Add(
			time.Minute * time.Duration(params.Ajustment.Dhuhr+params.MethodAjustment.Dhuhr),
		),
	)
	ashr := utils.RoundToNearestMinutes(
		tempAshr.Add(
			time.Minute * time.Duration(params.Ajustment.Asr+params.MethodAjustment.Asr),
		),
	)
	maghrib := utils.RoundToNearestMinutes(
		tempMaghrib.Add(
			time.Minute * time.Duration(params.Ajustment.Magrib+params.MethodAjustment.Magrib),
		),
	)
	isha := utils.RoundToNearestMinutes(
		tempIsha.Add(
			time.Minute * time.Duration(params.Ajustment.Isha+params.MethodAjustment.Isha),
		),
	)

	return &PrayerTimes{
		Imsak:             fajr.Add(-time.Minute * 10),
		Fajr:              fajr,
		Sunrise:           sunrise,
		Dhuhr:             dhuhr,
		Ashr:              ashr,
		Magrib:            maghrib,
		Isha:              isha,
		Coordinates:       coords,
		DateComponent:     date,
		CalculationParams: params,
	}, nil
}

func (prayer *PrayerTimes) CurrentPrayer() Prayer {
	currentTime := time.Now().Unix()
	switch {
	case prayer.Isha.Unix()-currentTime <= 0:
		return ISHA
	case prayer.Magrib.Unix()-currentTime <= 0:
		return MAGRIB
	case prayer.Ashr.Unix()-currentTime <= 0:
		return ASR
	case prayer.Dhuhr.Unix()-currentTime <= 0:
		return DHUHR
	case prayer.Fajr.Unix()-currentTime <= 0:
		return FAJR
	case prayer.Imsak.Unix()-currentTime <= 0:
		return IMSAK
	default:
		return NO_PRAYER
	}
}

func (prayer *PrayerTimes) NextPrayer() Prayer {
	currentPrayer := prayer.CurrentPrayer()
	return (currentPrayer + 1) % 7
}

func (pray *PrayerTimes) TimePray(prayer Prayer) time.Time {
	switch prayer {
	case IMSAK:
		return pray.Imsak
	case FAJR:
		return pray.Fajr
	case SUNRISE:
		return pray.Sunrise
	case DHUHR:
		return pray.Dhuhr
	case ASR:
		return pray.Ashr
	case MAGRIB:
		return pray.Magrib
	case ISHA:
		return pray.Isha
	default:
		return time.Time{}
	}
}

func (pray *PrayerTimes) SetTimeZone(tz string) error {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}
	pray.Imsak = pray.Imsak.In(loc)
	pray.Fajr = pray.Fajr.In(loc)
	pray.Sunrise = pray.Sunrise.In(loc)
	pray.Dhuhr = pray.Dhuhr.In(loc)
	pray.Ashr = pray.Ashr.In(loc)
	pray.Magrib = pray.Magrib.In(loc)
	pray.Isha = pray.Isha.In(loc)

	return nil
}
