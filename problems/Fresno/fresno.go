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

	// import domain
	searchDomain := corridor.CsvToDomain(1,
		"searchDomain.csv")

	///////////////////////////////////////////////////////////////////////////////////

	// initialize objectives
	searchObjectives := corridor.CsvToMultiObjective(1,
		"accessibility.csv",
		"slope.csv",
		"disturbance.csv")

	///////////////////////////////////////////////////////////////////////////////////

	// initialize parameters
	toyParameters := corridor.NewToyParameters(toyDomain)
	toyParameters.SrcSubs[0] = 8
	toyParameters.SrcSubs[1] = 14
	toyParameters.DstSubs[0] = toyDomain.Rows - 8
	toyParameters.DstSubs[1] = toyDomain.Rows - 14

	//////////////////////////////////////////////////////////////////////////////////

	//// evolve populations
	//toyEvolution := corridor.NewEvolution(toyParameters, toyDomain, toyObjectives)

	/////////////////////////////////////////////////////////////////////////////////////

	//// view output population
	//corridor.ViewPopulation(toyDomain, toyParameters, <-toyEvolution.Populations)

	///////////////////////////////////////////////////////////////////////////////////

	// stop clock and print runtime
	fmt.Printf("Elapsed Time: %s\n", time.Since(start))

	///////////////////////////////////////////////////////////////////////////////////
}
