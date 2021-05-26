package utils

import "math"

// IsWithinTolerance - given an expected and actual value, determine if the value is within the desired tolerance
func IsWithinTolerance(expected float64, actual float64, tolerance float64) bool {
	off := (actual - expected) / expected

	return math.Abs(off) <= tolerance
}
