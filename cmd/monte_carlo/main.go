package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/namelew/MonteCarloPointDistribution/internal/monte_carlo/experiment"
)

func main() {
	var (
		numberOfPoints   = flag.Uint("k", 300, "Number of points that will be generate on each experiment")
		powOfExperiments = flag.Uint("r", 6, "Greater exponent of the base 10 potency that set the number of runs on each experiment")
		radius           = flag.Float64("radius", 0.5, "Radius of the ring that will contain all points")
		seed             = flag.Uint("seed", uint(time.Now().Nanosecond()), "Seed value to generate the experiments")
	)

	flag.Parse()

	wg := sync.WaitGroup{}

	wg.Add(int(*powOfExperiments))

	result_filename := "monte-carlo-simulation-results-%d_%d_%d_%2f.csv"

	results_file, err := os.Create(fmt.Sprintf(result_filename, *seed, *numberOfPoints, *powOfExperiments, *radius))

	if err != nil {
		log.Fatal("Unable to create simulation results file: ", err)
	}

	points_filename := "monte-carlo-simulation-points-%d_%d_%d_%2f.csv"

	points_file, err := os.Create(fmt.Sprintf(points_filename, *seed, *numberOfPoints, *powOfExperiments, *radius))

	if err != nil {
		log.Fatal("Unable to create simulation points file: ", err)
	}

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

	for i := 0; i < int(*powOfExperiments); i++ {
		go experiment.Run(
			uint16(*numberOfPoints),
			uint8(i+1),
			uint32(*seed),
			*radius,
			&wg,
			points_file,
			results_file,
		)
	}

	wg.Wait()
}
