// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/ericdfournier/corridor"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	p := runtime.NumCPU()
	runtime.GOMAXPROCS(p)

	// start clock
	start := time.Now()

	// import domain from csv
	domainID := 1
	toyDomain := corridor.CsvToDomain(domainID, "/Users/ericfournier/go/src/github.com/ericdfournier/corridor/samples/convexSmall.csv")
	rows, cols := toyDomain.Matrix.Dims()

	// initialize objective variables
	objectiveID := 1
	objectiveCount := 3

	// initialize test objective
	toyObjectives := corridor.NewToyObjectives(objectiveID, rows, cols, objectiveCount)

	// initialize test parameters
	toyParameters := corridor.NewToyParameters(rows, cols)
	toyParameters.SrcSubs[0] = 180
	toyParameters.SrcSubs[1] = 100
	toyParameters.DstSubs[0] = rows - 180
	toyParameters.DstSubs[1] = rows - 100
	toyParameters.PopSize = 10
	toyParameters.EvoSize = 1

	/////////////////////////////////////////////////////////////////////

	// evolve populations
	toyEvolution := corridor.NewEvolution(toyParameters, toyDomain, toyObjectives)

	// print toy evolution to command line
	finalPop := <-toyEvolution.Populations
	fmt.Println(finalPop.AggregateMeanFitness)

	// view output population
	corridor.ViewPopulation(toyDomain, toyParameters, finalPop)

	// stop clock
	elapsed := time.Since(start)

	// print runtime
	fmt.Printf("Elapsed Time: %s\n", elapsed)
}
