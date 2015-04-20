// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/ericdfournier/corridor"
)

func main() {
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
	populationSize := 10
	evolutionSize := 1
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

	// output final population
	finalPopulation := <-searchEvolution.Populations

	///////////////////////////////////////////////////////////////////////////////////

	// output chromsome
	testChrom := <-finalPopulation.Chromosomes

	///////////////////////////////////////////////////////////////////////////////////

	// write chromosome to file
	corridor.ChromosomeToCsv(testChrom, "testChrom.csv")

	///////////////////////////////////////////////////////////////////////////////////

	// stop clock and print runtime
	fmt.Printf("Elapsed Time: %s\n", time.Since(start))

	///////////////////////////////////////////////////////////////////////////////////
}
