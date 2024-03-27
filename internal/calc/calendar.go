package calc

import (
	"fmt"
	"math"

	"github.com/taufiq30s/adzan/internal/utils"
)

func generateDate(b float64) utils.DateComponents {
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	var month, year float64
	day := b - d - math.Floor(30.6001*e)
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}

	if month < 3 {
		year = c - 4715
	} else {
		year = c - 4716
	}
	return utils.DateComponents{
		Year:  int16(year),
		Month: int8(month),
		Day:   int8(day),
	}
}

func ConvertHijrToGeorgian(date *utils.DateComponents) (utils.DateComponents, error) {
	N := float64(date.Day) + math.Floor(29.5001*(float64(date.Month)-1)+0.99)
	Q := math.Floor(float64(date.Year) / 30)
	R := math.Mod(float64(date.Year), 30)
	A := math.Floor((11*R + 3) / 30)
	W := (404 * Q) + (354 * R) + 208 + A
	Q1 := math.Floor(W / 1461)
	Q2 := math.Mod(W, 1461)
	G := 621 + 4*math.Floor((7*Q)+Q1)
	K := math.Floor(Q2 / 365.2422)
	E := math.Floor(365.2422 * K)
	J := Q2 - E + N - 1
	X := G + K

	if J < 365 {
		return utils.DateComponents{}, fmt.Errorf("invalid date")
	}
	if J > 366 && math.Mod(X, 4) == 0 {
		J -= 366
		X++
	} else if J > 365 && math.Mod(X, 4) > 0 {
		J -= 365
		X++
	}

	// Convert Julian Calendar to Georgian Calendar
	JD := math.Floor(365.25*(X-1)) + 1721423 + J
	alpha := math.Floor((JD - 1867216.25) / 36524.25)
	beta := 0.0

	if JD < 2299161 {
		beta = JD
	} else {
		beta = JD + 1 + alpha - math.Floor(alpha/4)
	}
	b := beta + 1524
	return generateDate(b), nil
}

func ConvertGeorgianToHijr(date utils.DateComponents) utils.DateComponents {
	// Convert Georgian to Julian
	if date.Month < 3 {
		date.Year--
		date.Month += 12
	}

	alpha := math.Floor(float64(date.Year) / 100)
	beta := 2 - alpha + math.Floor(alpha/4)
	b := math.Floor(365.25*float64(date.Year)) +
		math.Floor(30.6001*(float64(date.Month)+1)) +
		float64(date.Day) + 1722519 + beta
	julianDate := generateDate(b)

	W := 1.0
	if julianDate.Year%4 > 0 {
		W++
	}
	N := math.Floor((275*float64(julianDate.Month))/9) - W*
		math.Floor((float64(julianDate.Month)+9)/12) +
		float64(julianDate.Day) - 30
	A := float64(julianDate.Year - 623)
	B := math.Floor(A / 4)
	C := math.Mod(A, 4)
	C1 := 365.2501 * C
	C2 := math.Floor(C1)
	if (C1 - C2) > 0.5 {
		C2++
	}
	D := 1461*B + 170 + C2
	Q := math.Floor(D / 10631)
	R := math.Mod(D, 10631)
	J := math.Floor(R / 354)
	K := math.Mod(R, 354)
	O := math.Floor(((11 * J) + 14) / 30)
	H := (30 * Q) + J + 1
	JJ := K - O + N - 1 // Number Days of Hijr

	// Check Hijr is Leap year or not and substract it
	// with 354 if common year and 355 if leap year
	if JJ > 354 {
		CL := math.Mod(H, 30)
		DL := math.Mod((11*CL)+3, 30)
		if DL < 19 {
			JJ -= 354
		} else {
			JJ -= 355
		}
		H++

		if JJ == 0 {
			JJ += 355
			H--
		}
	}

	// Convert Month and Day from day Number JJ
	S := math.Floor((JJ - 1) / 29.5)
	return utils.DateComponents{
		Year:  int16(H),
		Month: int8(1 + S),
		Day:   int8(JJ - (29.5 * S)),
	}
}

func GetJulianDay(date *utils.DateComponents, hours float64) float64 {
	var A, B float64
	if date.Month < 3 {
		date.Year--
		date.Month += 12
	}
	A = math.Floor(float64(date.Year) / 100)
	B = 2 - A + math.Floor(A/4)
	return float64(int(365.25*float64(date.Year+4716))+int(30.6001*float64(date.Month+1))) + float64(date.Day) + B - 1524.5
}

func GetJulianCentury(jd float64) float64 {
	return (jd - 2451545) / 36525
}
