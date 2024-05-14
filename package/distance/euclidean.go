package distance

import "math"

type Point struct {
	X float64
	Y float64
}

func EuclidianDistance(a Point, b Point) float64 {
	differenceOfXs := math.Pow(a.X-b.X, 2)
	differenceOfYs := math.Pow(a.Y-b.Y, 2)

	return math.Pow(differenceOfXs+differenceOfYs, 0.5)
}
