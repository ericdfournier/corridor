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

// single small problem benchmark
func BenchmarkSingleSmall(b *testing.B) {

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

// medium single problem benchmark
func BenchmarkSingleMedium(b *testing.B) {

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

// large single problem benchmark
func BenchmarkSingleLarge(b *testing.B) {

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

// small problem monte carlo simulation benchmark
func BenchmarkMonteCarloSmall(b *testing.B) {

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

	// print simulation start message
	fmt.Printf("Simulation 1 of %v \n", sampleCount)

	// start sample runs
	for i := 0; i < sampleCount; i++ {

		fmt.Printf("Simulation %v of %v \n", i+1, sampleCount)

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

	// print sample size
	fmt.Println("Sample Size (N) = 100")

	// print mean aggregate fitness
	fmt.Printf("Mean Population Aggregate Fitnesses = %v \n", stat.Mean(aggMeanFitnesses, nil))

	// pritn standard deviation of aggregate fitness
	fmt.Printf("Standard Deviation of Aggregate Fitnesses = %v \n", stat.StdDev(aggMeanFitnesses, nil))

	// print mean runtimes
	fmt.Printf("Mean Runtime in Seconds = %v \n", stat.Mean(runtimes, nil))

	// print standard deviation of mean runtimes
	fmt.Printf("Standard Deviationof Runtimes in Seconds = %v \n", stat.StdDev(runtimes, nil))
}

// medium problem monte carlo simulation benchmark
func BenchmarkMonteCarloMedium(b *testing.B) {

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

	// print simulation start message
	fmt.Printf("Simulation 1 of %v \n", sampleCount)

	// start sample runs
	for i := 0; i < sampleCount; i++ {

		fmt.Printf("Simulation %v of %v \n", i+1, sampleCount)

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

	// print sample size
	fmt.Println("Sample Size (N) = 100")

	// print mean aggregate fitness
	fmt.Printf("Mean Population Aggregate Fitnesses = %v \n", stat.Mean(aggMeanFitnesses, nil))

	// pritn standard deviation of aggregate fitness
	fmt.Printf("Standard Deviation of Aggregate Fitnesses = %v \n", stat.StdDev(aggMeanFitnesses, nil))

	// print mean runtimes
	fmt.Printf("Mean Runtime in Seconds = %v \n", stat.Mean(runtimes, nil))

	// print standard deviation of mean runtimes
	fmt.Printf("Standard Deviationof Runtimes in Seconds = %v \n", stat.StdDev(runtimes, nil))
}

// large problem monte carlo simulation benchmark
func BenchmarkMonteCarloLarge(b *testing.B) {

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

	// print simulation start message
	fmt.Printf("Simulation 1 of %v \n", sampleCount)

	// start sample runs
	for i := 0; i < sampleCount; i++ {

		fmt.Printf("Simulation %v of %v \n", i+1, sampleCount)

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

	// print sample size
	fmt.Println("Sample Size (N) = 100")

	// print mean aggregate fitness
	fmt.Printf("Mean Population Aggregate Fitnesses = %v \n", stat.Mean(aggMeanFitnesses, nil))

	// pritn standard deviation of aggregate fitness
	fmt.Printf("Standard Deviation of Aggregate Fitnesses = %v \n", stat.StdDev(aggMeanFitnesses, nil))

	// print mean runtimes
	fmt.Printf("Mean Runtime in Seconds = %v \n", stat.Mean(runtimes, nil))

	// print standard deviation of mean runtimes
	fmt.Printf("Standard Deviationof Runtimes in Seconds = %v \n", stat.StdDev(runtimes, nil))

}
