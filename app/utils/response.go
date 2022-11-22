package utils

import "strconv"

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func StrToUint(str string) uint {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return uint(i)
}
