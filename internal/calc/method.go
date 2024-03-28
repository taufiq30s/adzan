package calc

// This file contains methods that using to calculate
// pray time.
// Reference https://www.salahtimes.com/faq/twilight
type CalculationMethod int

const (
	OTHER CalculationMethod = iota
	// Muslim World League
	// Used Fajr angle of 18 and an Isha angle of 17
	// Main regions: Europe, Far East, parts of the USA
	MUSLIM_WORLD_LEAGUE CalculationMethod = iota

	// Egyptian General Authority of Survey
	// Used Fajr angle of 19.5 and an Isha angle of 17.5
	// Main regions: Africa, Syria, Iraq, Lebanon, Malaysia, parts of the USA
	EGYPTIAN

	// Islamic Society of North America (ISNA)
	// This method is included for completeness, but is not recommended.
	// Uses a Fajr angle of 15 and an Isha angle of 15
	// Main regions: Parts of the USA, Canada, parts of the UK
	NORTH_AMERICA

	// Islamic Organisations Union of France
	// Uses a Fajr angle of 12 and an Isha angle of 12
	// Main Region: France
	UOIF

	// Umm Al-Qurra
	// Uses a Fajr and Isha angle of 18.5
	// Main Region: The Arabian Peninsula
	UMM_AL_QURRA

	// University Of Islamic Sciences, Karachi
	// Uses a Fajr and an Isha angle of 18
	// Main Region: Pakistan, Bangladesh, India, Afghanistan, Parts of Europe
	KARACHI

	// Majlis Ugama Islam Singapura
	// Uses a Fajr angle of 20 and an Isha angle of 18
	// Main Region: Singapore, Malaysia, Brunei, Indonesia
	SINGAPORE

	// Kuwait
	// Uses a Fajr angle of 18 and an Isha angle of 17.5
	KUWAIT

	// Qatar
	// Modified version of Umm Al-Qurra that used Fajr angle of 18
	QATAR

	KEMENAG
	MUHAMMADIYAH
)

func GetCalculationMethod(method CalculationMethod) *CalculationParameters {
	param := NewCalculationParameter()
	switch method {
	case MUSLIM_WORLD_LEAGUE:
		param.SetFajrAngle(18.0).SetIshaAngle(17.0).SetMethodAjustment(PrayerAjustment{
			Dhuhr: 1,
		})
	case EGYPTIAN:
		param.SetFajrAngle(19.5).SetIshaAngle(17.5).SetMethodAjustment(PrayerAjustment{
			Dhuhr: 1,
		})
	case NORTH_AMERICA:
		param.SetFajrAngle(15.0).SetIshaAngle(15.0).SetMethodAjustment(PrayerAjustment{
			Dhuhr: 1,
		})
	case UOIF:
		param.SetFajrAngle(12.0).SetIshaAngle(12.0)
	case UMM_AL_QURRA:
		param.SetFajrAngle(18.5).SetIshaInterval(90)
	case KARACHI:
		param.SetFajrAngle(18.0).SetIshaAngle(18.0).SetMethodAjustment(PrayerAjustment{
			Dhuhr: 1,
		})
	case SINGAPORE:
		param.SetFajrAngle(20.0).SetIshaAngle(18.0).SetMethodAjustment(PrayerAjustment{
			Dhuhr: 1,
		})
	case KUWAIT:
		param.SetFajrAngle(18.0).SetIshaAngle(17.5)
	case QATAR:
		param.SetFajrAngle(18.0).SetIshaInterval(90)
	case KEMENAG:
		param.SetFajrAngle(20.0).SetIshaAngle(18.0).SetMethodAjustment(PrayerAjustment{
			Fajr:   2,
			Dhuhr:  2,
			Asr:    2,
			Magrib: 2,
			Isha:   2,
		})
	case MUHAMMADIYAH:
		param.SetFajrAngle(18.0).SetIshaAngle(18.0).SetMethodAjustment(PrayerAjustment{
			Fajr:   2,
			Dhuhr:  2,
			Asr:    2,
			Magrib: 2,
			Isha:   2,
		})
	}
	return param
}
