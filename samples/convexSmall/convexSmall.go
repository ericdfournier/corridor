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

	///////////////////////////////////////////////////////////////////////////////////

	// start clock
	start := time.Now()

	///////////////////////////////////////////////////////////////////////////////////

	// create domain
	toyDomain := corridor.NewToyDomain(150, 150)

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
	toyParameters.PopSize = 1000
	toyParameters.EvoSize = 1000

	//////////////////////////////////////////////////////////////////////////////////

	// evolve populations
	toyEvolution := corridor.NewEvolution(toyParameters, toyDomain, toyObjectives)

	///////////////////////////////////////////////////////////////////////////////////

	// view output population
	finalPop := <-toyEvolution.Populations
	corridor.ViewPopulation(toyDomain, toyParameters, finalPop)

	///////////////////////////////////////////////////////////////////////////////////

	// select elite fraction
	eliteCount := 10
	toyElites := corridor.NewEliteSet(eliteCount, finalPop)

	///////////////////////////////////////////////////////////////////////////////////

	// write elite set to file
	corridor.EliteSetToCsv(toyElites, "output.csv")

	///////////////////////////////////////////////////////////////////////////////////

	// stop clock and print runtime
	fmt.Printf("Elapsed Time: %s\n", time.Since(start))

	///////////////////////////////////////////////////////////////////////////////////
}
