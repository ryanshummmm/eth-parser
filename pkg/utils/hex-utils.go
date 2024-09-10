package utils

import (
	"errors"
	"strconv"
)

func HexToInt(hex string) (int, error) {
	if len(hex) > 2 && hex[:2] == "0x" {
		parsed, err := strconv.ParseInt(hex[2:], 16, 64)
		if err != nil {
			return 0, err
		}
		return int(parsed), nil
	}
	return 0, errors.New("invalid hex string")
}
