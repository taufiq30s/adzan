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
	Date      string `json:"date"`
	Formatted string `json:"formatted"`
	Timezone  string `json:"timezone"`
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

func ShowCurrentHijrDate(w http.ResponseWriter, r *http.Request) {
	rawLat := r.URL.Query().Get("lat")
	rawLng := r.URL.Query().Get("lng")
	currentTime := time.Now().UTC()
	var currentTz string
	w.Header().Set("Content-Type", "application/json")

	if rawLat != "" && rawLng != "" {
		coordinate, err := convertCoordinateToFloat64(rawLat, rawLng)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}

		timezone, err := utils.GetTimeZone(coordinate.Latitude, coordinate.Longitude)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		currentTime = time.Now().In(timezone)
		currentTz = timezone.String()
	} else {
		currentTz, _ = currentTime.Zone()
	}
	hijrDate := calc.ConvertGeorgianToHijr(*utils.NewDateComponents(currentTime))

	data := hijrData{
		Date: fmt.Sprintf(
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
	}
	fmt.Fprint(w, string(jsonData))
}
