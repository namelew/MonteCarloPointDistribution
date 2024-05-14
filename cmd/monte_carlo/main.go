package main

import (
	"flag"
	"sync"

	"github.com/namelew/MonteCarloPointDistribution/internal/monte_carlo/experiment"
)

func main() {
	var (
		numberOfPoints   = flag.Uint("k", 300, "Number of points that will be generate on each experiment")
		powOfExperiments = flag.Uint("r", 6, "Greater exponent of the base 10 potency that set the number of runs on each experiment")
	)

	flag.Parse()

	wg := sync.WaitGroup{}

	wg.Add(int(*powOfExperiments))

	for i := 0; i < int(*powOfExperiments); i++ {
		go experiment.Run(uint16(*numberOfPoints), uint8(i+1), &wg)
	}

	wg.Wait()
}
