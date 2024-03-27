package utils

type success struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func SuccessResponse(data any) *success {
	return &success{
		Success: true,
		Data:    data,
	}
}
