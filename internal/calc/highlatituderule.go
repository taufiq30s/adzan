package calc

/*
High Latitude Rule

The purpose of this rule is to determine Fajr and Isha in a location with a high latitude
where the usual formula cannot be to determine it

Reference: http://praytimes.org/wiki/Calculation#Higher_Latitudes
*/

type NightPortion struct {
	fajr float32
	Isha float32
}

type HighLatitudeRule int8

const (
	NONE HighLatitudeRule = iota

	// In this method, the period from sunset to sunrise is divided into two halves.
	// The first half is considered to be the "night" and the other half as "day break".
	// Fajr and Isha in this method are assumed to be at mid-night during the abnormal periods.
	MIDDLE_OF_THE_NIGHT

	// In this method, the period between sunset and sunrise is divided into seven parts.
	// Isha begins after the first one-seventh part, and Fajr is at the beginning of the seventh part.
	ONE_SEVENTH_OF_THE_NIGHT

	// This is an intermediate solution, used by some recent prayer time calculators.
	// Let α be the twilight angle for Isha, and let t = α/60. The period between sunset and
	// sunrise is divided into t parts. Isha begins after the first part. \
	// For example, if the twilight angle for Isha is 15, then Isha begins at the end of the first quarter (15/60)
	// of the night. Time for Fajr is calculated similarly.
	ANGLE_BASED_METHOD
)

func NewNightPortion(fajr float32, isha float32) *NightPortion {
	return &NightPortion{
		fajr: fajr,
		Isha: isha,
	}
}
