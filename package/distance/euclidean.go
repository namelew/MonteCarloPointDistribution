package distance

import (
	"math"
	"math/rand"
)

type Point struct {
	X float64
	Y float64
}

func EuclidianPoint(seed uint32, r, a, b float64) Point {
	random := rand.New(rand.NewSource(int64(seed)))

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
