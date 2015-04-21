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
	toyDomain := corridor.CsvToDomain("convexSmall.csv")

	///////////////////////////////////////////////////////////////////////////////////

	// view domain
	corridor.ViewDomain(toyDomain)

	///////////////////////////////////////////////////////////////////////////////////

	// initialize objectives
	objectiveCount := 3
	toyObjectives := corridor.NewToyObjectives(toyDomain.Rows, toyDomain.Cols, objectiveCount)

	///////////////////////////////////////////////////////////////////////////////////

	// initialize parameters
	toyParameters := corridor.NewToyParameters(toyDomain)
	toyParameters.SrcSubs[0] = 8
	toyParameters.SrcSubs[1] = 14
	toyParameters.DstSubs[0] = toyDomain.Rows - 8
	toyParameters.DstSubs[1] = toyDomain.Rows - 14
	toyParameters.PopSize = 100
	toyParameters.EvoSize = 100

	//////////////////////////////////////////////////////////////////////////////////

	// evolve populations
	toyEvolution := corridor.NewEvolution(toyParameters, toyDomain, toyObjectives)

	///////////////////////////////////////////////////////////////////////////////////

	// view output population
	//corridor.ViewPopulation(toyDomain, toyParameters, <-toyEvolution.Populations)

	///////////////////////////////////////////////////////////////////////////////////

	// select elite fraction
	toyElites := corridor.NewEliteFraction(0.5, <-toyEvolution.Populations)

	for i := 0; i < 50; i++ {
		fmt.Println(toyElites[i].AggregateFitness)
	}

	///////////////////////////////////////////////////////////////////////////////////

	// stop clock and print runtime
	fmt.Printf("Elapsed Time: %s\n", time.Since(start))

	///////////////////////////////////////////////////////////////////////////////////
}
