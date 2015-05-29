// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"fmt"
	"runtime"
	"testing"
)

// small problem benchmark
func BenchmarkSmall(b *testing.B) {

	// set max processing units
	runtime.GOMAXPROCS(1)

	// initialize domain
	sampleDomain := NewSampleDomain(20, 20)
	sampleDomain.BndCnt = 3

	// initialize objectives
	objectiveCount := 3
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

	// print top individual fitness
	fmt.Println("Population Mean Fitness =")
	fmt.Println(finalPop.MeanFitness)
}

// parallel small problem benchmark
func BenchmarkParallelSmall(b *testing.B) {

	// set max processing units
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	// initialize domain
	sampleDomain := NewSampleDomain(20, 20)
	sampleDomain.BndCnt = 3

	// initialize objectives
	objectiveCount := 3
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

	// print top individual fitness
	fmt.Println("Population Mean Fitness =")
	fmt.Println(finalPop.MeanFitness)
}

// medium problem benchmark
func BenchmarkMedium(b *testing.B) {

	// set max processing units
	runtime.GOMAXPROCS(1)

	// initialize domain
	sampleDomain := NewSampleDomain(20, 20)
	sampleDomain.BndCnt = 3

	// initialize objectives
	objectiveCount := 3
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)
	sampleParameters.PopSize = 10000

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

	// print top individual fitness
	fmt.Println("Population Mean Fitness =")
	fmt.Println(finalPop.MeanFitness)
}

// parallel medium problem benchmark
func BenchmarkParallelMedium(b *testing.B) {

	// set max processing units
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	// initialize domain
	sampleDomain := NewSampleDomain(20, 20)
	sampleDomain.BndCnt = 3

	// initialize objectives
	objectiveCount := 3
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)
	sampleParameters.PopSize = 10000

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

	// print top individual fitness
	fmt.Println("Population Mean Fitness =")
	fmt.Println(finalPop.MeanFitness)
}

// large problem benchmark
func BenchmarkLarge(b *testing.B) {

	// set max processing units
	runtime.GOMAXPROCS(1)

	// initialize domain
	sampleDomain := NewSampleDomain(20, 20)
	sampleDomain.BndCnt = 3

	// initialize objectives
	objectiveCount := 3
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)
	sampleParameters.PopSize = 100000

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

	// print top individual fitness
	fmt.Println("Population Mean Fitness =")
	fmt.Println(finalPop.MeanFitness)
}

// parallel large problem benchmark
func BenchmarkParallelLarge(b *testing.B) {

	// set max processing units
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	// initialize domain
	sampleDomain := NewSampleDomain(20, 20)
	sampleDomain.BndCnt = 3

	// initialize objectives
	objectiveCount := 3
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)
	sampleParameters.PopSize = 100000

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

	// print top individual fitness
	fmt.Println("Population Mean Fitness =")
	fmt.Println(finalPop.MeanFitness)
}
