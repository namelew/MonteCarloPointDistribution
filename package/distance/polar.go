package distance

import (
	"math"
	"math/rand"
)

func PolarPoint(seed uint32, r, a, b float64) Point {
	random := rand.New(rand.NewSource(int64(seed)))

	randomRadius := r * math.Sqrt(random.Float64())
	randomAngle := random.Float64() * 2 * math.Pi

	x := a + randomRadius*math.Cos(randomAngle)
	y := b + randomRadius*math.Sin(randomAngle)

	return Point{X: x, Y: y}
}
