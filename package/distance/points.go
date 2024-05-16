package distance

type Point struct {
	X float64
	Y float64
}

type CoordinationType uint8

const (
	EUCLIDIAN CoordinationType = 0
	POLAR     CoordinationType = 1
)
