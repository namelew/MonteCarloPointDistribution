package distance

type Point struct {
	X float32
	Y float32
}

type CoordinationType uint8

const (
	EUCLIDIAN CoordinationType = 0
	POLAR     CoordinationType = 1
)
