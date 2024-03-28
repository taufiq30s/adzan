package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type TimeZone struct {
	Time        string `json:"time"`
	GmtOffset   int8   `json:"gmtOffset"`
	TimezoneId  string `json:"timezoneId"`
	CountryCode string `json:"countryCode"`
}

const uri string = "http://api.geonames.org/timezoneJSON?lat=%f&lng=%f&username=%s"

func GetTimeZone(lat float64, lng float64) (timezone *TimeZone, statusCode int, err error) {
	tzCh := make(chan *TimeZone)
	errCh := make(chan error)
	statusCodeCh := make(chan int)

	defer catch(&err)

	username, e := getUsername()
	if e != nil {
		return timezone, 500, e
	}

	url := fmt.Sprintf(uri, lat, lng, username)
	go func() {
		tz := TimeZone{}
		res, e := http.Get(url)
		if e != nil {
			statusCodeCh <- res.StatusCode
			errCh <- e
			return
		}
		defer res.Body.Close()

		raw, e := io.ReadAll(res.Body)
		if e != nil {
			statusCodeCh <- 500
			errCh <- e
			return
		}

		if e = json.Unmarshal(raw, &tz); e != nil {
			statusCodeCh <- 500
			errCh <- e
			return
		}

		if len(tz.TimezoneId) == 0 {
			statusCodeCh <- 404
			errCh <- fmt.Errorf("error: Timezone not found")
			return
		}
		tzCh <- &tz
	}()
	select {
	case timezone = <-tzCh:
		return timezone, 200, err
	case err = <-errCh:
		return timezone, <-statusCodeCh, err
	}
}

func getUsername() (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		return "", err
	}
	return os.Getenv("GEOLOCATION_USERNAME"), nil
}

func catch(err *error) {
	if r := recover(); r != nil {
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}
