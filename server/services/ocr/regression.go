package ocr

// LinearRegression represents a regression line with a slope and y-intercept
type LinearRegression struct {
	Slope     float64
	Intercept float64
}

// CalculateExpectedY calculates the expected Y value using a value and a regression line
func CalculateExpectedY(x float64, lr *LinearRegression) float64 {
	return x*lr.Slope + lr.Intercept
}

// NewLinearRegression creates a linear regression based on two points
func NewLinearRegression(x1 float64, y1 float64, x2 float64, y2 float64) *LinearRegression {
	slope := (y2 - y1) / (x2 - x1)
	intercept := y1 - x1*slope
	return &LinearRegression{
		Slope:     slope,
		Intercept: intercept,
	}
}
