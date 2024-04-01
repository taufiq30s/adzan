package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/taufiq30s/adzan/internal/calc"
	"github.com/taufiq30s/adzan/internal/utils"
)

type adzanData struct {
	Date     string `json:"date"`
	HijrDate string `json:"hijrDate"`
	Imsak    string `json:"imsak"`
	Fajr     string `json:"fajr"`
	Sunrise  string `json:"sunrise"`
	Dhuhr    string `json:"dhuhr"`
	Ashr     string `json:"asr"`
	Magrib   string `json:"magrib"`
	Isha     string `json:"isha"`
}

func validateYearAndMonthParameter(rawYear string, rawMonth string) bool {
	yearR, _ := regexp.Compile("^[12][0-9]{3}$")
	monthR, _ := regexp.Compile(`\b(?:1[0-2]|[1-9])\b`)

	if !yearR.MatchString(rawYear) {
		return false
	}
	if !monthR.MatchString(rawMonth) {
		return false
	}
	return true
}

func getTimeZoneAndCurrentDate(coordinate *utils.Coordinates) (*utils.TimeZone, *time.Time, int, error) {
	timezone, statusCode, err := utils.GetTimeZone(coordinate.Latitude, coordinate.Longitude)
	if err != nil {
		return nil, nil, statusCode, err
	}

	currentDate, err := time.Parse("2006-01-02 15:04", timezone.Time)
	if err != nil {
		return nil, nil, 500, err
	}

	return timezone, &currentDate, 200, nil
}

func getAdzanData(date *time.Time, coordinate *utils.Coordinates, timezone *utils.TimeZone) (*calc.PrayerTimes, int, error) {
	param := calc.GetCalculationMethod(calc.MUSLIM_WORLD_LEAGUE).SetMazhab(calc.SYAFI)
	if timezone.CountryCode == "ID" {
		param = calc.GetCalculationMethod(calc.KEMENAG).SetMazhab(calc.SYAFI)
	}

	adzan, err := calc.NewPrayerTimes(coordinate, utils.NewDateComponents(*date), param)
	if err != nil {
		return nil, 500, err
	}
	adzan.SetTimeZone(timezone.TimezoneId)

	return adzan, 200, nil
}

func TodayAdzan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rawLat := r.URL.Query().Get("lat")
	rawLng := r.URL.Query().Get("lng")

	if rawLat == "" || rawLng == "" {
		http.Error(w, "please input coordinate of location", 400)
		return
	}

	coordinate, err := convertCoordinateToFloat64(rawLat, rawLng)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	timezone, currentDate, statusCode, err := getTimeZoneAndCurrentDate(coordinate)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	adzan, statusCode, err := getAdzanData(currentDate, coordinate, timezone)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	hijrDate := calc.ConvertGeorgianToHijr(*utils.NewDateComponents(*currentDate))
	formattedDate := currentDate.Format("January 02, 2006")
	if timezone.CountryCode == "ID" {
		formattedDate = currentDate.Format("02 January 2006")
	}

	data := adzanData{
		Date: formattedDate,
		HijrDate: fmt.Sprintf(
			"%v %v %v",
			hijrDate.Day,
			hijrMontName[int(hijrDate.Month)],
			hijrDate.Year,
		),
		Imsak:   adzan.Imsak.Format("15:04"),
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

func MonthlyAdzan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rawLat := r.URL.Query().Get("lat")
	rawLng := r.URL.Query().Get("lng")
	rawYear := r.URL.Query().Get("year")
	rawMonth := r.URL.Query().Get("month")

	if rawLat == "" || rawLng == "" {
		http.Error(w, "please input coordinate of location", 400)
		return
	}

	if !validateYearAndMonthParameter(rawYear, rawMonth) {
		http.Error(w, "please input valid year and month (ex: 2024 for year and 05 or 5 for month)", 400)
		return
	}

	coordinate, err := convertCoordinateToFloat64(rawLat, rawLng)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	timezone, _, statusCode, err := getTimeZoneAndCurrentDate(coordinate)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	if len(rawMonth) == 1 {
		rawMonth = "0" + rawMonth
	}
	startDate, _ := time.Parse("2006-01-02", rawYear+"-"+rawMonth+"-01")
	var prayerTimes []adzanData

	for d := startDate; d.Month() == startDate.Month(); d = d.AddDate(0, 0, 1) {
		prayerTime, statusCode, err := getAdzanData(&d, coordinate, timezone)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		hijrDate := calc.ConvertGeorgianToHijr(*utils.NewDateComponents(d))
		formattedDate := d.Format("January 02, 2006")
		if timezone.CountryCode == "ID" {
			formattedDate = d.Format("02 January 2006")
		}

		prayerTimes = append(prayerTimes, adzanData{
			Date: formattedDate,
			HijrDate: fmt.Sprintf(
				"%v %v %v",
				hijrDate.Day,
				hijrMontName[int(hijrDate.Month)],
				hijrDate.Year,
			),
			Imsak:   prayerTime.Imsak.Format("15:04"),
			Fajr:    prayerTime.Fajr.Format("15:04"),
			Sunrise: prayerTime.Sunrise.Format("15:04"),
			Dhuhr:   prayerTime.Dhuhr.Format("15:04"),
			Ashr:    prayerTime.Ashr.Format("15:04"),
			Magrib:  prayerTime.Magrib.Format("15:04"),
			Isha:    prayerTime.Isha.Format("15:04"),
		})
	}

	jsonData, err := json.Marshal(utils.SuccessResponse(prayerTimes))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	fmt.Fprint(w, string(jsonData))
}
