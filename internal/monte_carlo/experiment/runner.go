package experiment

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"

	"github.com/namelew/MonteCarloPointDistribution/package/distance"
)

type experiment struct {
	r                     uint8
	k                     uint16
	seed                  uint32
	ctype                 distance.CoordinationType
	randomNumberGenerator *rand.Rand
	radius                float64
	wg                    *sync.WaitGroup
	mutex                 *sync.Mutex
}

type sample struct {
	distance.Point
	distanceToCenter float64
}

var CIRCLE_CENTER = distance.Point{X: 0, Y: 0}

func (s *sample) inValid(radius float64) bool {
	return s.distanceToCenter <= radius
}

func Run(k uint16, r uint8, seed uint32, ctype distance.CoordinationType, radius float64, wg *sync.WaitGroup, PRChannel chan []string, RRChannel chan []string, rgn *rand.Rand) {
	defer wg.Done()

	numberOfRuns := int(math.Pow10(int(r)))

	wgRuns := sync.WaitGroup{}
	mutexRuns := sync.Mutex{}

	wgRuns.Add(numberOfRuns)

	current := experiment{
		r:                     r,
		k:                     k,
		seed:                  seed,
		ctype:                 ctype,
		randomNumberGenerator: rgn,
		radius:                radius,
		wg:                    &wgRuns,
		mutex:                 &mutexRuns,
	}

	for i := 0; i < numberOfRuns; i++ {
		go current.Run(i, PRChannel, RRChannel)
	}

	wgRuns.Wait()
}

func (e *experiment) Run(rid int, pointsRegisters chan []string, resultsRegisters chan []string) {
	defer e.wg.Done()

	distances := make([]*sample, e.k)
	var sumPointDistance float64 = 0

	for i := range distances {

		var new sample

		switch e.ctype {
		case distance.EUCLIDIAN:
			e.mutex.Lock()
			point := distance.Point{
				X: e.randomNumberGenerator.Float64(),
				Y: e.randomNumberGenerator.Float64(),
			}
			e.mutex.Unlock()

			new = sample{
				Point:            point,
				distanceToCenter: distance.EuclidianDistance(point, CIRCLE_CENTER),
			}
			// Validation Step
			for {
				if new.inValid(e.radius) {
					break
				}

				e.mutex.Lock()
				point := distance.Point{
					X: e.randomNumberGenerator.Float64(),
					Y: e.randomNumberGenerator.Float64(),
				}
				e.mutex.Unlock()

				new = sample{
					Point:            point,
					distanceToCenter: distance.EuclidianDistance(point, CIRCLE_CENTER),
				}
			}
		case distance.POLAR:
			e.mutex.Lock()
			point := distance.PolarPoint(e.randomNumberGenerator, e.radius, CIRCLE_CENTER.X, CIRCLE_CENTER.Y)
			e.mutex.Unlock()

			new = sample{
				Point:            point,
				distanceToCenter: distance.EuclidianDistance(point, CIRCLE_CENTER),
			}
		}

		distances[i] = &new
		sumPointDistance += new.distanceToCenter
	}

	meanDistance := sumPointDistance / float64(e.k)

	var sumLocalVariation float64 = 0.0

	for i := range distances {
		sumLocalVariation += math.Pow(distances[i].distanceToCenter-meanDistance, 2)
	}

	variance := sumLocalVariation / (float64(e.k) - 1)
	stdDeviation := math.Sqrt(variance)

	for i := range distances {
		csv_string := fmt.Sprintf("%d,%d,%d,%d,%f,%d,%f,%f,%f",
			1,
			e.k,
			e.r,
			rid,
			e.radius,
			e.seed,
			distances[i].Point.X,
			distances[i].Point.Y,
			distances[i].distanceToCenter,
		)

		register := strings.Split(csv_string, ",")

		pointsRegisters <- register
	}

	csv_string := fmt.Sprintf("%d,%d,%d,%d,%f,%d,%f,%f,%f",
		1,
		e.k,
		e.r,
		rid,
		e.radius,
		e.seed,
		meanDistance,
		variance,
		stdDeviation,
	)

	resultsRegisters <- strings.Split(csv_string, ",")
}
