package calc

import "fmt"

type CalculationParameters struct {
	// Method that will be used
	Method CalculationMethod

	// The angle of sun to calculate Fajr
	FajrAngle float32

	// The angle of sun to calculate Isha
	IshaAngle float32

	// Minutes after magrib (Time of Isha = Magrib + IshaInterval)
	IshaInterval int8

	// The Juristic method to calculate ashr
	Mazhab Mazhab

	HighLatitudeRule HighLatitudeRule

	// Manual Ajustment
	Ajustment PrayerAjustment

	// Ajustment that set by a calculation method
	MethodAjustment PrayerAjustment
}

func NewCalculationParameter() *CalculationParameters {
	return &CalculationParameters{
		Method:           OTHER,
		FajrAngle:        0.0,
		IshaAngle:        0.0,
		IshaInterval:     0,
		Mazhab:           SYAFI,
		HighLatitudeRule: MIDDLE_OF_THE_NIGHT,
		Ajustment:        PrayerAjustment{},
		MethodAjustment:  PrayerAjustment{},
	}
}

func (param *CalculationParameters) SetMethod(method CalculationMethod) *CalculationParameters {
	param.Method = method
	return param
}

func (param *CalculationParameters) SetFajrAngle(angle float32) *CalculationParameters {
	param.FajrAngle = angle
	return param
}

func (param *CalculationParameters) SetIshaAngle(angle float32) *CalculationParameters {
	param.IshaAngle = angle
	return param
}

func (param *CalculationParameters) SetIshaInterval(interval int8) *CalculationParameters {
	param.IshaInterval = interval
	return param
}

func (param *CalculationParameters) SetMazhab(mazhab Mazhab) *CalculationParameters {
	param.Mazhab = mazhab
	return param
}

func (param *CalculationParameters) SetAjustment(ajustment PrayerAjustment) *CalculationParameters {
	param.Ajustment = ajustment
	return param
}

func (param *CalculationParameters) SetHighLatitudeRule(highLatitudeRule HighLatitudeRule) *CalculationParameters {
	param.HighLatitudeRule = highLatitudeRule
	return param
}

func (param *CalculationParameters) SetMethodAjustment(ajusment PrayerAjustment) *CalculationParameters {
	param.MethodAjustment = ajusment
	return param
}

func (param *CalculationParameters) GetNightPortion() (*NightPortion, error) {
	switch param.HighLatitudeRule {
	case MIDDLE_OF_THE_NIGHT:
		return NewNightPortion(1.0/2.0, 1.0/2.0), nil
	case ONE_SEVENTH_OF_THE_NIGHT:
		return NewNightPortion(1.0/7.0, 1.0/7.0), nil
	case ANGLE_BASED_METHOD:
		return NewNightPortion(param.FajrAngle/60.0, param.IshaAngle/60.0), nil
	default:
		return nil, fmt.Errorf("invalid high latitude rule")
	}
}
