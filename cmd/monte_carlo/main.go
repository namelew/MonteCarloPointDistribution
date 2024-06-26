package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/namelew/MonteCarloPointDistribution/internal/monte_carlo/experiment"
	"github.com/namelew/MonteCarloPointDistribution/package/distance"
)

func main() {
	var (
		numberOfPoints       = flag.Uint("k", 300, "Number of points that will be generate on each experiment")
		coordinationsType    = flag.Uint("t", uint(distance.EUCLIDIAN), "Type of coordinations that will be used to generate points in the circle")
		powOfExperiments     = flag.Uint("r", 6, "Greater exponent of the base 10 potency that set the number of runs on each experiment")
		radius               = flag.Float64("radius", 0.5, "Radius of the ring that will contain all points")
		seed                 = flag.Uint("seed", uint(time.Now().Nanosecond()), "Seed value to generate the experiments")
		pointsBufferMaxSize  = flag.Float64("points-buffer", 0.1, "Max number of points stored on main memory")
		resultsBufferMaxSize = flag.Float64("results-buffer", 0.1, "Max number of points stored on main memory")
	)

	flag.Parse()

	wg := sync.WaitGroup{}

	wg.Add(int(*powOfExperiments))

	result_filename := "monte-carlo-simulation-results-%d_%d_%d_%d_%2f.csv"

	results_file, err := os.Create(fmt.Sprintf(result_filename, *coordinationsType, *seed, *numberOfPoints, *powOfExperiments, *radius))

	if err != nil {
		log.Fatal("Unable to create simulation results file: ", err)
	}

	points_filename := "monte-carlo-simulation-points-%d_%d_%d_%d_%2f.csv"

	points_file, err := os.Create(fmt.Sprintf(points_filename, *coordinationsType, *seed, *numberOfPoints, *powOfExperiments, *radius))

	if err != nil {
		log.Fatal("Unable to create simulation points file: ", err)
	}

	pointsMemThreshould := 0

	for i := 0; i < int(*powOfExperiments); i++ {
		pointsMemThreshould += int(math.Pow10(i+1)) * int(*numberOfPoints)
	}

	pointsMemThreshould = int(math.Floor(float64(pointsMemThreshould) * *pointsBufferMaxSize))

	resultsMemThreshould := 0

	for i := 0; i < int(*powOfExperiments); i++ {
		resultsMemThreshould += int(math.Pow10(i + 1))
	}

	resultsMemThreshould = int(math.Floor(float64(resultsMemThreshould) * *resultsBufferMaxSize))

	pointsWriter := csv.NewWriter(points_file)
	resultsWriter := csv.NewWriter(results_file)

	err = pointsWriter.Write([]string{"scenario", "number-of-points(k)", "number-of-runs(r)", "radius", "seed", "x", "y", "distance"})

	if err != nil {
		log.Fatal("Unable to write point line on points file: ", err)
	}

	err = resultsWriter.Write([]string{"scenario", "number-of-points(k)", "number-of-runs(r)", "radius", "seed", "mean", "variance", "stdDeviation"})

	if err != nil {
		log.Fatal("Unable to write point line on points file: ", err)
	}

	pointsRegisters := make([][]string, 0)
	resultsRegisters := make([][]string, 0)

	random := rand.New(rand.NewSource(int64(*seed)))
	PRChannel := make(chan []string)
	RRChannel := make(chan []string)

	for i := 0; i < int(*powOfExperiments); i++ {
		go experiment.Run(
			&experiment.GlobalOptions{
				NumberOfPoints: uint16(*numberOfPoints),
				Seed:           uint32(*seed),
				CType:          distance.CoordinationType(*coordinationsType),
				Radius:         float32(*radius),
				RNG:            random,
			},
			uint8(i+1),
			&wg,
			PRChannel,
			RRChannel,
		)
	}

	go func() {
		wg.Wait()
		close(PRChannel)
		close(RRChannel)
	}()

	wgCollectors := sync.WaitGroup{}

	wgCollectors.Add(2)

	go func() {
		for point := range PRChannel {
			pointsRegisters = append(pointsRegisters, point)

			if len(pointsRegisters) >= pointsMemThreshould {
				err = pointsWriter.WriteAll(pointsRegisters)

				if err != nil {
					log.Fatal("Unable to write point line on points file: ", err)
					return
				}

				pointsWriter.Flush()

				pointsRegisters = make([][]string, 0)
			}
		}
		wgCollectors.Done()
	}()

	go func() {
		for result := range RRChannel {
			resultsRegisters = append(resultsRegisters, result)

			if len(resultsRegisters) >= resultsMemThreshould {
				err = resultsWriter.WriteAll(resultsRegisters)

				if err != nil {
					log.Fatal("Unable to write point line on points file: ", err)
					return
				}

				resultsWriter.Flush()

				resultsRegisters = make([][]string, 0)
			}
		}

		wgCollectors.Done()
	}()

	wgCollectors.Wait()

	err = pointsWriter.WriteAll(pointsRegisters)

	if err != nil {
		log.Fatal("Unable to write point line on points file: ", err)
		return
	}

	pointsWriter.Flush()

	err = resultsWriter.WriteAll(resultsRegisters)

	if err != nil {
		log.Fatal("Unable to write point line on points file: ", err)
		return
	}

	resultsWriter.Flush()
}
