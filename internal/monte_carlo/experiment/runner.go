package experiment

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"

	"github.com/namelew/MonteCarloPointDistribution/package/distance"
)

type GlobalOptions struct {
	NumberOfPoints uint16
	Seed           uint32
	Radius         float64
	CType          distance.CoordinationType
	RNG            *rand.Rand
}

type experiment struct {
	r          uint8
	globalVars *GlobalOptions
	wg         *sync.WaitGroup
	mutex      *sync.Mutex
}

type sample struct {
	distance.Point
	distanceToCenter float64
}

var CIRCLE_CENTER = distance.Point{X: 0, Y: 0}

func (s *sample) inValid(radius float64) bool {
	return s.distanceToCenter <= radius
}

func Run(global *GlobalOptions, r uint8, wg *sync.WaitGroup, PRChannel chan []string, RRChannel chan []string) {
	defer wg.Done()

	numberOfRuns := int(math.Pow10(int(r)))

	wgRuns := sync.WaitGroup{}
	mutexRuns := sync.Mutex{}

	wgRuns.Add(numberOfRuns)

	current := experiment{
		r:          r,
		globalVars: global,
		wg:         &wgRuns,
		mutex:      &mutexRuns,
	}

	for i := 0; i < numberOfRuns; i++ {
		go current.Run(i, PRChannel, RRChannel)
	}

	wgRuns.Wait()
}

func (e *experiment) Run(rid int, pointsRegisters chan []string, resultsRegisters chan []string) {
	defer e.wg.Done()

	distances := make([]*sample, e.globalVars.NumberOfPoints)
	var sumPointDistance float64 = 0

	for i := range distances {

		var new sample

		switch e.globalVars.CType {
		case distance.EUCLIDIAN:
			e.mutex.Lock()
			point := distance.Point{
				X: e.globalVars.RNG.Float64(),
				Y: e.globalVars.RNG.Float64(),
			}
			e.mutex.Unlock()

			new = sample{
				Point:            point,
				distanceToCenter: distance.EuclidianDistance(point, CIRCLE_CENTER),
			}
			// Validation Step
			for {
				if new.inValid(e.globalVars.Radius) {
					break
				}

				e.mutex.Lock()
				point := distance.Point{
					X: e.globalVars.RNG.Float64(),
					Y: e.globalVars.RNG.Float64(),
				}
				e.mutex.Unlock()

				new = sample{
					Point:            point,
					distanceToCenter: distance.EuclidianDistance(point, CIRCLE_CENTER),
				}
			}
		case distance.POLAR:
			e.mutex.Lock()
			point := distance.PolarPoint(e.globalVars.RNG, e.globalVars.Radius, CIRCLE_CENTER.X, CIRCLE_CENTER.Y)
			e.mutex.Unlock()

			new = sample{
				Point:            point,
				distanceToCenter: distance.EuclidianDistance(point, CIRCLE_CENTER),
			}
		}

		distances[i] = &new
		sumPointDistance += new.distanceToCenter
	}

	meanDistance := sumPointDistance / float64(e.globalVars.NumberOfPoints)

	var sumLocalVariation float64 = 0.0

	for i := range distances {
		sumLocalVariation += math.Pow(distances[i].distanceToCenter-meanDistance, 2)
	}

	variance := sumLocalVariation / (float64(e.globalVars.NumberOfPoints) - 1)
	stdDeviation := math.Sqrt(variance)

	for i := range distances {
		csv_string := fmt.Sprintf("%d,%d,%d,%d,%f,%d,%f,%f,%f",
			1,
			e.globalVars.NumberOfPoints,
			e.r,
			rid,
			e.globalVars.Radius,
			e.globalVars.Seed,
			distances[i].Point.X,
			distances[i].Point.Y,
			distances[i].distanceToCenter,
		)

		register := strings.Split(csv_string, ",")

		pointsRegisters <- register
	}

	csv_string := fmt.Sprintf("%d,%d,%d,%d,%f,%d,%f,%f,%f",
		1,
		e.globalVars.NumberOfPoints,
		e.r,
		rid,
		e.globalVars.Radius,
		e.globalVars.Seed,
		meanDistance,
		variance,
		stdDeviation,
	)

	resultsRegisters <- strings.Split(csv_string, ",")
}
