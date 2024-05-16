package distance

import (
	"math"
	"math/rand"
)

func PolarPoint(random *rand.Rand, r, a, b float64) Point {
	randomRadius := r * random.Float64()
	randomAngle := random.Float64() * 2 * math.Pi

	x := a + randomRadius*math.Cos(randomAngle)
	y := b + randomRadius*math.Sin(randomAngle)

	return Point{X: x, Y: y}
}
