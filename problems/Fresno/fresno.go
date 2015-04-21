// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/ericdfournier/corridor"
)

func main() {
	///////////////////////////////////////////////////////////////////////////////////

	// set max processing units
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	fmt.Println("CPU Count:")
	fmt.Println(cpuCount)

	///////////////////////////////////////////////////////////////////////////////////

	// start clock
	start := time.Now()

	///////////////////////////////////////////////////////////////////////////////////

	// import source subscripts
	source := corridor.CsvToSubs("sourceSubs.csv")

	// import destination subscripts
	destination := corridor.CsvToSubs("destinationSubs.csv")

	///////////////////////////////////////////////////////////////////////////////////

	// import domain
	searchDomain := corridor.CsvToDomain("searchDomain.csv")

	///////////////////////////////////////////////////////////////////////////////////

	// initialize objectives
	searchObjectives := corridor.CsvToMultiObjective(
		"accessibility.csv",
		"slope.csv",
		"disturbance.csv")

	///////////////////////////////////////////////////////////////////////////////////

	//// initialize parameters
	populationSize := 10000
	evolutionSize := 1000
	randomness := 1.0

	searchParameters := corridor.NewParameters(
		source,
		destination,
		populationSize,
		evolutionSize,
		randomness)

	//////////////////////////////////////////////////////////////////////////////////

	// evolve populations
	searchEvolution := corridor.NewEvolution(
		searchParameters,
		searchDomain,
		searchObjectives)

	///////////////////////////////////////////////////////////////////////////////////

	// generate elite set
	eliteCount := 100
	eliteSet := corridor.NewEliteSet(eliteCount, <-searchEvolution.Populations)

	///////////////////////////////////////////////////////////////////////////////////

	// write chromosome to file
	corridor.EliteSetToCsv(eliteSet, "eliteSet.csv")

	///////////////////////////////////////////////////////////////////////////////////

	// stop clock and print runtime
	fmt.Printf("Elapsed Time: %s\n", time.Since(start))

	///////////////////////////////////////////////////////////////////////////////////
}
