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
	sourceSubs := corridor.CsvToSubs("sourceSubs.csv")

	// import destination subscripts
	destinationSubs := corridor.CsvToSubs("destinationSubs.csv")

	fmt.Println(sourceSubs)
	fmt.Println(destinationSubs)

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
	//searchParameters := corridor.NewParameters(

	//	)

	//////////////////////////////////////////////////////////////////////////////////

	////// evolve populations
	//toyEvolution := corridor.NewEvolution(
	//	searchParameters,
	//	searchDomain,
	//	searchObjectives
	//	)

	/////////////////////////////////////////////////////////////////////////////////////

	//// view output population
	//corridor.ViewPopulation(toyDomain, toyParameters, <-toyEvolution.Populations)

	///////////////////////////////////////////////////////////////////////////////////

	// stop clock and print runtime
	fmt.Printf("Elapsed Time: %s\n", time.Since(start))

	///////////////////////////////////////////////////////////////////////////////////
}
