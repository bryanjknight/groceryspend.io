package utils

import (
	"strconv"
	"strings"
)

// ParseStringToUSDAmount parse a string value into a float32 assuming the value is something like $123.45
func ParseStringToUSDAmount(s string) (float32, error) {
	// strip the $ out, convert to a float32
	sNoDollarSign := strings.ReplaceAll(s, "$", "")
	sTrimmedNoDollarSign := strings.TrimSpace(sNoDollarSign)
	val, err := strconv.ParseFloat(sTrimmedNoDollarSign, 32)
	if err != nil {
		return 0, err
	}
	return float32(val), nil
}
