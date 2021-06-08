package utils_test

import (
	"testing"

	"groceryspend.io/server/utils"
)

func TestIsWithinTolerance(t *testing.T) {

	type test struct {
		expected   float64
		actual     float64
		tolarence  float64
		shouldPass bool
	}

	tests := []test{
		{expected: 10, actual: 10, tolarence: 0, shouldPass: true},
		{expected: 10, actual: 10, tolarence: 0.5, shouldPass: true},
		{expected: 100, actual: 10, tolarence: 0.1, shouldPass: false},
		{expected: 100, actual: 90, tolarence: 0.1, shouldPass: true},
		{expected: 100, actual: 110, tolarence: 0.1, shouldPass: true},
	}

	for _, testInstance := range tests {
		didPass := utils.IsWithinTolerance(testInstance.expected, testInstance.actual, testInstance.tolarence)
		if didPass != testInstance.shouldPass {
			t.Errorf("Expected e(%v) a(%v) with t(%v) to be (%t), but got %t",
				testInstance.expected, testInstance.actual, testInstance.tolarence, testInstance.shouldPass, didPass)
		}
	}
}

func TestIsLessThanWithinTolerance(t *testing.T) {

	type test struct {
		expected   float64
		actual     float64
		tolarence  float64
		shouldPass bool
	}

	tests := []test{
		{expected: 10, actual: 10, tolarence: 0, shouldPass: true},
		{expected: 10, actual: 10, tolarence: 0.5, shouldPass: true},
		{expected: 100, actual: 10, tolarence: 0.1, shouldPass: true},
		{expected: 100, actual: 90, tolarence: 0.1, shouldPass: true},
		{expected: 100, actual: 110, tolarence: 0.1, shouldPass: true},
		{expected: 100, actual: 210, tolarence: 0.1, shouldPass: false},
	}

	for _, testInstance := range tests {
		didPass := utils.IsLessThanWithinTolerance(testInstance.expected, testInstance.actual, testInstance.tolarence)
		if didPass != testInstance.shouldPass {
			t.Errorf("Expected e(%v) a(%v) with t(%v) to be (%t), but got %t",
				testInstance.expected, testInstance.actual, testInstance.tolarence, testInstance.shouldPass, didPass)
		}
	}
}

func TestIsGreaterThanWithinTolerance(t *testing.T) {

	type test struct {
		expected   float64
		actual     float64
		tolarence  float64
		shouldPass bool
	}

	tests := []test{
		{expected: 10, actual: 10, tolarence: 0, shouldPass: true},
		{expected: 10, actual: 10, tolarence: 0.5, shouldPass: true},
		{expected: 100, actual: 10, tolarence: 0.1, shouldPass: false},
		{expected: 100, actual: 90, tolarence: 0.1, shouldPass: false},
		{expected: 100, actual: 110, tolarence: 0.1, shouldPass: true},
		{expected: 100, actual: 210, tolarence: 0.1, shouldPass: true},
	}

	for _, testInstance := range tests {
		didPass := utils.IsGreaterThanWithinTolerance(testInstance.expected, testInstance.actual, testInstance.tolarence)
		if didPass != testInstance.shouldPass {
			t.Errorf("Expected e(%v) a(%v) with t(%v) to be (%t), but got %t",
				testInstance.expected, testInstance.actual, testInstance.tolarence, testInstance.shouldPass, didPass)
		}
	}
}
