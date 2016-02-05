/* Copyright ©2015 The corridor Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file. */

package corridor

import (
	"fmt"
	"github.com/gonum/stat"
	"runtime"
	"testing"
	"time"
)

// parallel small problem benchmark
func BenchmarkSmall(b *testing.B) {

	// set max processing units
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	// initialize integer constants
	const (
		xDim           int = 20
		yDim           int = 20
		bandCount      int = 3
		objectiveCount int = 3
		populationSize int = 1000
		sampleCount    int = 100
	)

	// initialize domain
	sampleDomain := NewSampleDomain(xDim, yDim)
	sampleDomain.BndCnt = bandCount

	// initialize objectives
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)
	sampleParameters.PopSize = populationSize

	// initialize results slices
	aggMeanFitnesses := make([]float64, sampleCount)
	runtimes := make([]float64, sampleCount)

	// start sample runs
	for i := 0; i < sampleCount; i++ {

		// start clock
		start := time.Now()

		// generate evolution
		toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

		// write runtime
		runtimes[i] = time.Since(start).Seconds()

		// extract final output population
		finalPop := <-toyEvolution.Populations

		// write aggregate mean fitness
		aggMeanFitnesses[i] = finalPop.AggregateMeanFitness

	}

	// print mean aggregate fitness
	fmt.Println("Mean Population Aggregate Fitnesses =")
	fmt.Println(stat.Mean(aggMeanFitnesses, nil))

	// print mean runtimes
	fmt.Println("Mean Runtime in Seconds")
	fmt.Println(stat.Mean(runtimes, nil))
}

// TODO UPDATE BENCHMARKS TO MONTE CARLO SAMPLING FOR ALL!!!

// medium problem benchmark
func BenchmarkMedium(b *testing.B) {

	// set max processing units
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	// initialize integer constants
	const (
		xDim           int = 20
		yDim           int = 20
		bandCount      int = 3
		objectiveCount int = 3
		populationSize int = 10000
	)

	// initialize domain
	sampleDomain := NewSampleDomain(xDim, yDim)
	sampleDomain.BndCnt = bandCount

	// initialize objectives
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)
	sampleParameters.PopSize = populationSize

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view sample chromosome
	ViewChromosome(sampleDomain, sampleParameters, <-finalPop.Chromosomes)

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

	// print top individual fitness
	fmt.Println("Population Mean Fitness =")
	fmt.Println(finalPop.MeanFitness)
}

// large problem benchmark
func BenchmarkLarge(b *testing.B) {

	// set max processing units
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	// initialize integer constants
	const (
		xDim           int = 20
		yDim           int = 20
		bandCount      int = 3
		objectiveCount int = 3
		populationSize int = 100000
	)

	// initialize domain
	sampleDomain := NewSampleDomain(xDim, yDim)
	sampleDomain.BndCnt = bandCount

	// initialize objectives
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)
	sampleParameters.PopSize = populationSize

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view sample chromosome
	ViewChromosome(sampleDomain, sampleParameters, <-finalPop.Chromosomes)

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

	// print top individual fitness
	fmt.Println("Population Mean Fitness =")
	fmt.Println(finalPop.MeanFitness)
}
