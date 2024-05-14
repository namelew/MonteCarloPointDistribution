package experiment

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/namelew/MonteCarloPointDistribution/package/distance"
)

type experiment struct {
	points_file  *os.File
	results_file *os.File
	r            uint8
	k            uint16
	seed         uint32
	radius       float64
	wg           *sync.WaitGroup
	mutex        *sync.Mutex
}

type sample struct {
	distance.Point
	distanceToCenter float64
}

var CIRCLE_CENTER = distance.Point{X: 0, Y: 0}

func (s *sample) inValid(radius float64) bool {
	return s.distanceToCenter <= radius
}

func Run(k uint16, r uint8, seed uint32, radius float64, wg *sync.WaitGroup, points_file *os.File, results_file *os.File) {
	defer wg.Done()

	numberOfRuns := int(math.Pow10(int(r)))

	wgRuns := sync.WaitGroup{}
	mutexRuns := sync.Mutex{}

	wgRuns.Add(numberOfRuns)

	for i := 0; i < numberOfRuns; i++ {
		current := experiment{
			points_file:  points_file,
			results_file: results_file,
			r:            r,
			k:            k,
			seed:         seed,
			radius:       radius,
			wg:           &wgRuns,
			mutex:        &mutexRuns,
		}

		go current.Run()
	}

	wgRuns.Wait()
}

func (e *experiment) Run() {
	defer e.wg.Done()

	distances := make([]*sample, e.k)
	var sumPointDistance float64 = 0

	randomNumberGenerator := rand.New(rand.NewSource(int64(e.seed)))

	for i := range distances {
		point := distance.Point{
			X: randomNumberGenerator.Float64(),
			Y: randomNumberGenerator.Float64(),
		}

		new := sample{
			Point:            point,
			distanceToCenter: distance.EuclidianDistance(point, CIRCLE_CENTER),
		}

		// Validation Step
		for {
			if new.inValid(e.radius) {
				break
			}

			point := distance.Point{
				X: randomNumberGenerator.Float64(),
				Y: randomNumberGenerator.Float64(),
			}

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

	e.mutex.Lock()
	pointsWriter := csv.NewWriter(e.points_file)
	resultsWriter := csv.NewWriter(e.results_file)

	for i := range distances {
		csv_string := fmt.Sprintf("%d,%d,%d,%f,%d,%f,%f,%f",
			1,
			e.k,
			e.r,
			e.radius,
			e.seed,
			distances[i].Point.X,
			distances[i].Point.Y,
			distances[i].distanceToCenter,
		)

		err := pointsWriter.Write(strings.Split(csv_string, ","))

		if err != nil {
			log.Println("Unable to write point line on points file: ", err)
			continue
		}
	}

	pointsWriter.Flush()

	csv_string := fmt.Sprintf("%d,%d,%d,%f,%d,%f,%f,%f",
		1,
		e.k,
		e.r,
		e.radius,
		e.seed,
		meanDistance,
		variance,
		stdDeviation,
	)

	err := resultsWriter.Write(strings.Split(csv_string, ","))

	if err != nil {
		log.Println("Unable to write simulation result line on results file: ", err)
	}

	resultsWriter.Flush()

	e.mutex.Unlock()
}
