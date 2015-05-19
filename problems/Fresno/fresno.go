// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"runtime"
	"time"

	"github.com/ericdfournier/corridor"
)

func main() {
	///////////////////////////////////////////////////////////////////////////////////

	// set max processing units
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

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

	// initialize parameters
	populationSize := 100
	evolutionSize := 1000
	randomness := 1.0

	// generate parameter structure
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

	// initialize elite count
	eliteCount := 10

	// extract elite set
	eliteSet := corridor.NewEliteSet(eliteCount, <-searchEvolution.Populations, searchParameters)

	///////////////////////////////////////////////////////////////////////////////////

	// write elite set to file
	corridor.EliteSetToCsv(eliteSet, "fresno_p-100000_e-1000_eliteSet.csv")

	///////////////////////////////////////////////////////////////////////////////////

	// write log data to file
	corridor.RuntimeLogToCsv(searchEvolution, time.Since(start), "log.csv")

	///////////////////////////////////////////////////////////////////////////////////
}
