package pkg

import "testing"

func TestStringToUSDAmount(t *testing.T) {

	type test struct {
		input  string
		output float32
	}

	tests := []test{
		{input: "$123.45", output: 123.45},
		{input: "$123", output: 123.0},
		{input: " $123.45", output: 123.45},
		{input: " -$3.00", output: -3.00},
	}

	for _, test := range tests {
		actual, err := ParseStringToUSDAmount(test.input)

		if err != nil {
			t.Errorf("%v", err)
		}

		if actual != test.output {
			t.Errorf("Expected %v, got %v", test.output, actual)
		}
	}
}
