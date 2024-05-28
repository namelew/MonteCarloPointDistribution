package distance

import (
	"math"
	"math/rand"
)

func EuclidianPoint(random *rand.Rand, r, a, b float32) Point {
	randomRadius := float64(r) * math.Sqrt(random.Float64())
	randomAngle := random.Float64() * 2 * math.Pi

	x := a + float32(randomRadius*math.Cos(randomAngle))
	y := b + float32(randomRadius*math.Sin(randomAngle))

	return Point{X: x, Y: y}
}

func EuclidianDistance(a Point, b Point) float32 {
	differenceOfXs := math.Pow(float64(a.X-b.X), 2)
	differenceOfYs := math.Pow(float64(a.Y-b.Y), 2)

	return float32(math.Pow(differenceOfXs+differenceOfYs, 0.5))
}
