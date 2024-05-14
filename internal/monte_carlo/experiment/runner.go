package experiment

import (
	"math"
	"sync"
)

func Run(k uint16, r uint8, wg *sync.WaitGroup) {
	defer wg.Done()

	numberOfRuns := int(math.Pow10(int(k)))

	wgRuns := sync.WaitGroup{}

	wgRuns.Add(numberOfRuns)

	for i := 0; i < numberOfRuns; i++ {

	}

	wgRuns.Wait()
}
