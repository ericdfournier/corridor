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

	fmt.Println(source)
	fmt.Println(destination)

	///////////////////////////////////////////////////////////////////////////////////

	// import domain
	searchDomain := corridor.CsvToDomain("searchDomain.csv")

	fmt.Println(searchDomain)

	///////////////////////////////////////////////////////////////////////////////////

	// initialize objectives
	searchObjectives := corridor.CsvToMultiObjective(
		"accessibility.csv",
		"slope.csv",
		"disturbance.csv")

	fmt.Println(searchObjectives)

	///////////////////////////////////////////////////////////////////////////////////

	//// initialize parameters
	populationSize := 100
	evolutionSize := 100
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

	fmt.Println(searchEvolution)

	///////////////////////////////////////////////////////////////////////////////////

	// view output population
	corridor.ViewPopulation(
		searchDomain,
		searchParameters,
		<-searchEvolution.Populations)

	///////////////////////////////////////////////////////////////////////////////////

	// stop clock and print runtime
	fmt.Printf("Elapsed Time: %s\n", time.Since(start))

	///////////////////////////////////////////////////////////////////////////////////
}
