package distance

import (
	"math"
	"math/rand"
)

func EuclidianPoint(random *rand.Rand, r, a, b float64) Point {
	randomRadius := r * math.Sqrt(random.Float64())
	randomAngle := random.Float64() * 2 * math.Pi

	x := a + randomRadius*math.Cos(randomAngle)
	y := b + randomRadius*math.Sin(randomAngle)

	return Point{X: x, Y: y}
}

func EuclidianDistance(a Point, b Point) float64 {
	differenceOfXs := math.Pow(a.X-b.X, 2)
	differenceOfYs := math.Pow(a.Y-b.Y, 2)

	return math.Pow(differenceOfXs+differenceOfYs, 0.5)
}
