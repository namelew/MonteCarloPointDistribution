package distance

import (
	"math"
	"math/rand"
)

func PolarPoint(random *rand.Rand, r, a, b float32) Point {
	randomRadius := float64(r) * random.Float64()
	randomAngle := random.Float64() * 2 * math.Pi

	x := a + float32(randomRadius*math.Cos(randomAngle))
	y := b + float32(randomRadius*math.Sin(randomAngle))

	return Point{X: x, Y: y}
}
