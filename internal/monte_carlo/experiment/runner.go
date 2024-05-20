package experiment

import (
	"fmt"
	"log"
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

func Run(global *GlobalOptions, r uint8, wg *sync.WaitGroup, PRChannel chan []string, RRChannel chan []string) {
	defer wg.Done()

	log.Printf("Starting Simulation...\n\tscenario:%d\n\tr: %d\n\tk:%d", global.CType, r, global.NumberOfPoints)

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
	log.Printf("Finishing Simulation...\n\tscenario:%d\n\tr: %d\n\tk:%d", global.CType, r, global.NumberOfPoints)
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
			point := distance.EuclidianPoint(e.globalVars.RNG, e.globalVars.Radius, CIRCLE_CENTER.X, CIRCLE_CENTER.Y)
			e.mutex.Unlock()

			new = sample{
				Point:            point,
				distanceToCenter: distance.EuclidianDistance(point, CIRCLE_CENTER),
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
		csv_string := fmt.Sprintf("%d,%d,%d,%f,%d,%f,%f,%f",
			e.globalVars.CType,
			e.globalVars.NumberOfPoints,
			e.r,
			e.globalVars.Radius,
			e.globalVars.Seed,
			distances[i].Point.X,
			distances[i].Point.Y,
			distances[i].distanceToCenter,
		)

		register := strings.Split(csv_string, ",")

		pointsRegisters <- register
	}

	csv_string := fmt.Sprintf("%d,%d,%d,%f,%d,%f,%f,%f",
		e.globalVars.CType,
		e.globalVars.NumberOfPoints,
		e.r,
		e.globalVars.Radius,
		e.globalVars.Seed,
		meanDistance,
		variance,
		stdDeviation,
	)

	log.Printf("Execution %d of Scenario %d and R %d finished. \n\tmean:%f\n\tvariance:%f\n\tstdDeviation:%f\n.",
		rid,
		e.globalVars.CType,
		e.r,
		meanDistance,
		variance,
		stdDeviation,
	)

	resultsRegisters <- strings.Split(csv_string, ",")
}
