package distance

import "math"

func EuclidianDistance(a []float64, b []float64) float64 {
	differenceOfXs := math.Pow(a[0]-b[0], 2)
	differenceOfYs := math.Pow(a[1]-b[1], 2)
	return math.Pow(differenceOfXs+differenceOfYs, 0.5)
}
