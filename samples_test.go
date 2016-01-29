/* Copyright ©2015 The corridor Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file. */

package corridor

import (
	"fmt"
	"runtime"
	"testing"
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

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

	// view sample chromosome
	ViewChromosome(sampleDomain, sampleParameters, <-finalPop.Chromosomes)

	// print top individual fitness
	fmt.Println("Population Mean Fitness =")
	fmt.Println(finalPop.MeanFitness)
}

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
