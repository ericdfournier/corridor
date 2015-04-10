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
	domainID := 1
	toyDomain := corridor.CsvToDomain(domainID, "convexSmall.csv")
	rows, cols := toyDomain.Matrix.Dims()

	///////////////////////////////////////////////////////////////////////////////////

	// view domain
	corridor.ViewDomain(toyDomain)

	///////////////////////////////////////////////////////////////////////////////////

	// initialize objectives
	objectiveID := 1
	objectiveCount := 3
	toyObjectives := corridor.NewToyObjectives(objectiveID, rows, cols, objectiveCount)

	///////////////////////////////////////////////////////////////////////////////////

	// initialize parameters
	toyParameters := corridor.NewToyParameters(toyDomain)
	toyParameters.SrcSubs[0] = 8
	toyParameters.SrcSubs[1] = 14
	toyParameters.DstSubs[0] = rows - 8
	toyParameters.DstSubs[1] = rows - 14

	//////////////////////////////////////////////////////////////////////////////////

	// evolve populations
	toyEvolution := corridor.NewEvolution(toyParameters, toyDomain, toyObjectives)

	///////////////////////////////////////////////////////////////////////////////////

	// view output population
	corridor.ViewPopulation(toyDomain, toyParameters, <-toyEvolution.Populations)

	///////////////////////////////////////////////////////////////////////////////////

	// stop clock and print runtime
	fmt.Printf("Elapsed Time: %s\n", time.Since(start))

	///////////////////////////////////////////////////////////////////////////////////
}
