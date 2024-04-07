package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/taufiq30s/adzan/internal/calc"
	"github.com/taufiq30s/adzan/internal/utils"
)

type hijrData struct {
	GregorianDate string `json:"gregorianDate"`
	HijrDate      string `json:"hijrDate"`
	Formatted     string `json:"formatted"`
	Timezone      string `json:"timezone"`
}

var monthName = map[int]string{
	1:  "Muharram",
	2:  "Safar",
	3:  "Rabiul Awal",
	4:  "Rabiul Akhir",
	5:  "Jumadil Awal",
	6:  "Jumadil Akhir",
	7:  "Rajab",
	8:  "Syaban",
	9:  "Ramadhan",
	10: "Syawal",
	11: "Zulkaidah",
	12: "Dzulhijjah",
}

func convertCoordinateToFloat64(rawLat string, rawLng string) (*utils.Coordinates, error) {
	lat, err := strconv.ParseFloat(rawLat, 64)
	if err != nil {
		return nil, err
	}

	lng, err := strconv.ParseFloat(rawLng, 64)
	if err != nil {
		return nil, err
	}

	coordinate, err := utils.NewCoordinates(lat, lng)
	if err != nil {
		return nil, err
	}
	return coordinate, nil
}

func setTimeBasedOnCoordinate(rawLat string, rawLng string, date time.Time) (time.Time, string, error) {
	var currentTz string

	if rawLat != "" && rawLng != "" {
		coordinate, err := convertCoordinateToFloat64(rawLat, rawLng)
		if err != nil {
			return time.Time{}, "", err
		}

		timezone, err := utils.GetTimeZone(coordinate.Latitude, coordinate.Longitude)
		if err != nil {
			return time.Time{}, "", err
		}
		date = date.In(timezone)
		currentTz = timezone.String()
	} else {
		currentTz, _ = date.Zone()
	}
	return date, currentTz, nil
}

func ShowCurrentHijrDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawLat := r.URL.Query().Get("lat")
	rawLng := r.URL.Query().Get("lng")
	var rawDate time.Time
	var err error

	if r.URL.Query().Get("date") != "" {
		rawDate, err = time.Parse("2006-01-02", r.URL.Query().Get("date"))
		if err != nil {
			http.Error(w, "Invalid date format. Valid format (YYYY-MM-DD). Example: 2024-04-01", 400)
			return
		}
	} else {
		rawDate = time.Now().UTC()
	}
	date, currentTz, err := setTimeBasedOnCoordinate(rawLat, rawLng, rawDate)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	hijrDate := calc.ConvertGeorgianToHijr(*utils.NewDateComponents(date))

	data := hijrData{
		GregorianDate: fmt.Sprintf(
			"%d-%d-%d",
			date.Year(),
			int(date.Month()),
			date.Day(),
		),
		HijrDate: fmt.Sprintf(
			"%v-%v-%v",
			hijrDate.Year,
			hijrDate.Month,
			hijrDate.Day,
		),
		Formatted: fmt.Sprintf(
			"%v %v %v H",
			hijrDate.Day,
			monthName[int(hijrDate.Month)],
			hijrDate.Year,
		),
		Timezone: currentTz,
	}
	jsonData, err := json.Marshal(utils.SuccessResponse(data))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprint(w, string(jsonData))
}
