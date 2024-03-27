package routes

import (
	"net/http"

	"github.com/taufiq30s/adzan/internal/api"
)

func NewRoute() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hijr", api.ShowCurrentHijrDate)
	return mux
}

// func dummy(w http.ResponseWriter, r *http.Request) {
// 	currentDate := time.Now()
// 	utc := currentDate.UTC()
// 	date := utils.NewDateComponents(utc)
// 	params := calc.GetCalculationMethod(calc.KEMENAG)
// 	params.Mazhab = calc.SYAFI

// 	cords, err := utils.NewCoordinates(3.58333333, 98.65000000)
// 	if err != nil {
// 		fmt.Fprintf(w, "%v", err.Error())
// 		return
// 	}

// 	prayerTimes, err := calc.NewPrayerTimes(cords, date, params)
// 	if err != nil {
// 		fmt.Fprintln(w, err.Error())
// 		return
// 	}

// 	fmt.Fprintf(w, "Current Time:\nLocal: %+v\nUTC: %+v\n\n", currentDate, utc)
// 	fmt.Fprintf(w, "UTC: \n")
// 	fmt.Fprintf(w, "Fajr : %+v\n", prayerTimes.Fajr)
// 	fmt.Fprintf(w, "Sunrise : %+v\n", prayerTimes.Sunrise)
// 	fmt.Fprintf(w, "Dhuhr : %+v\n", prayerTimes.Dhuhr)
// 	fmt.Fprintf(w, "Asr : %+v\n", prayerTimes.Ashr)
// 	fmt.Fprintf(w, "Magrib : %+v\n", prayerTimes.Magrib)
// 	fmt.Fprintf(w, "Isha : %+v\n\n", prayerTimes.Isha)

// 	err = prayerTimes.SetTimeZone("Asia/Jakarta")
// 	if err != nil {
// 		fmt.Fprintln(w, err.Error())
// 		return
// 	}

// 	fmt.Fprintf(w, "Local: \n")
// 	fmt.Fprintf(w, "Fajr : %+v\n", prayerTimes.Fajr)
// 	fmt.Fprintf(w, "Sunrise : %+v\n", prayerTimes.Sunrise)
// 	fmt.Fprintf(w, "Dhuhr : %+v\n", prayerTimes.Dhuhr)
// 	fmt.Fprintf(w, "Asr : %+v\n", prayerTimes.Ashr)
// 	fmt.Fprintf(w, "Magrib : %+v\n", prayerTimes.Magrib)
// 	fmt.Fprintf(w, "Isha : %+v\n", prayerTimes.Isha)
// }
