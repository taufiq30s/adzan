package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/taufiq30s/adzan/internal/calc"
	"github.com/taufiq30s/adzan/internal/utils"
)

type adzanData struct {
	Date     string `json:"date"`
	HijrDate string `json:"hijrDate"`
	Fajr     string `json:"fajr"`
	Sunrise  string `json:"sunrise"`
	Dhuhr    string `json:"dhuhr"`
	Ashr     string `json:"asr"`
	Magrib   string `json:"magrib"`
	Isha     string `json:"isha"`
}

func TodayAdzan(w http.ResponseWriter, r *http.Request) {
	rawLat := r.URL.Query().Get("lat")
	rawLng := r.URL.Query().Get("lng")

	if rawLat == "" || rawLng == "" {
		http.Error(w, "please input coordinate of location", 400)
	}

	coordinate, err := convertCoordinateToFloat64(rawLat, rawLng)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	timezone, statusCode, err := utils.GetTimeZone(coordinate.Latitude, coordinate.Longitude)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
	}

	currentTime, err := time.Parse("2006-01-02 15:04", timezone.Time)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
	}

	param := calc.GetCalculationMethod(calc.MUSLIM_WORLD_LEAGUE).SetMazhab(calc.SYAFI)
	if timezone.CountryCode == "ID" {
		param = calc.GetCalculationMethod(calc.KEMENAG).SetMazhab(calc.SYAFI)
	}

	adzan, err := calc.NewPrayerTimes(coordinate, utils.NewDateComponents(currentTime), param)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	adzan.SetTimeZone(timezone.TimezoneId)

	hijrDate := calc.ConvertGeorgianToHijr(*utils.NewDateComponents(currentTime))

	data := adzanData{
		Date: currentTime.Format("January 02, 2006"),
		HijrDate: fmt.Sprintf(
			"%v %v %v",
			hijrDate.Day,
			hijrMontName[int(hijrDate.Month)],
			hijrDate.Year,
		),
		Fajr:    adzan.Fajr.Format("15:04"),
		Sunrise: adzan.Sunrise.Format("15:04"),
		Dhuhr:   adzan.Dhuhr.Format("15:04"),
		Ashr:    adzan.Ashr.Format("15:04"),
		Magrib:  adzan.Magrib.Format("15:04"),
		Isha:    adzan.Isha.Format("15:04"),
	}

	jsonData, err := json.Marshal(utils.SuccessResponse(data))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	fmt.Fprint(w, string(jsonData))
}
