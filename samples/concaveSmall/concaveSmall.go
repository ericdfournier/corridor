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

	////////////////////////////////////////////////////////////////////////////////////

	// import domain
	domainID := 1
	toyDomain := corridor.CsvToDomain(domainID, "concaveSmall.csv")
	rows, cols := toyDomain.Matrix.Dims()

	///////////////////////////////////////////////////////////////////////////////////

	// initialize objectives
	objectiveID := 1
	objectiveCount := 3
	toyObjectives := corridor.NewToyObjectives(objectiveID, rows, cols, objectiveCount)

	///////////////////////////////////////////////////////////////////////////////////

	// initialize parameters
	toyParameters := corridor.NewToyParameters(rows, cols)
	toyParameters.SrcSubs[0] = 40
	toyParameters.SrcSubs[1] = 40
	toyParameters.DstSubs[0] = 105
	toyParameters.DstSubs[1] = 95
	toyParameters.PopSize = 10
	toyParameters.EvoSize = 1

	///////////////////////////////////////////////////////////////////////////////////

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
