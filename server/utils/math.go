package utils

import "math"

// IsWithinTolerance - given an expected and actual value, determine if the value is within the desired tolerance
func IsWithinTolerance(expected float64, actual float64, tolerance float64) bool {
	off := (actual - expected) / expected

	return math.Abs(off) <= tolerance
}

// IsLessThanWithinTolerance - given an expected value, is the actual value under the expected with some tolerance?
func IsLessThanWithinTolerance(expected float64, actual float64, tolerance float64) bool {
	return (1.0-tolerance)*actual <= expected
}

// IsGreaterThanWithinTolerance - given an expected value, is the actual value under the expected with some tolerance?
func IsGreaterThanWithinTolerance(expected float64, actual float64, tolerance float64) bool {
	return (1.0+tolerance)*actual >= expected
}
