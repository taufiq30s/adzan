package utils

import (
	"fmt"
	"time"

	"github.com/ringsaturn/tzf"
)

type TimeZone struct {
	Time        string `json:"time"`
	GmtOffset   int8   `json:"gmtOffset"`
	TimezoneId  string `json:"timezoneId"`
	CountryCode string `json:"countryCode"`
}

var finder tzf.F

func init() {
	var err error
	defer catch(&err)
	finder, err = tzf.NewDefaultFinder()
	if err != nil {
		panic(err)
	}
}

func GetTimeZone(lat float64, lng float64) (timezone *time.Location, err error) {
	tzName := finder.GetTimezoneName(lng, lat)
	tz, err := time.LoadLocation(tzName)
	if err != nil {
		return nil, err
	}
	return tz, err
}

func catch(err *error) {
	if r := recover(); r != nil {
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}
