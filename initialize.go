/* Copyright ©2015 The corridor Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file. */

package corridor

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"sort"
	"sync"

	"github.com/gonum/diff/fd"
	"github.com/gonum/matrix/mat64"
	"github.com/satori/go.uuid"
)

// new problem parameters function
func NewParameters(sourceSubscripts, destinationSubscripts []int, populationSize, evolutionSize int, randomnessCoefficient float64) *Parameters {

	// set default integer parameter values
	var (
		mutationCount  int = 1
		maxConcurrency int = runtime.NumCPU()
	)

	// set default floating point parameter values
	var (
		mutationFraction     float64 = 0.2
		selectionFraction    float64 = 0.5
		selectionProbability float64 = 0.8
	)

	// return output
	return &Parameters{
		SrcSubs: sourceSubscripts,
		DstSubs: destinationSubscripts,
		RndCoef: randomnessCoefficient,
		PopSize: populationSize,
		SelFrac: selectionFraction,
		SelProb: selectionProbability,
		MutaCnt: mutationCount,
		MutaFrc: mutationFraction,
		EvoSize: evolutionSize,
		ConSize: maxConcurrency,
	}
}

// new domain initialization function
func NewDomain(domainMatrix *mat64.Dense) *Domain {

	// get domain size
	rows, cols := domainMatrix.Dims()

	// compute band count
	bandCount := 2 + (int(math.Floor(math.Sqrt(math.Pow(float64(rows), 2.0)+math.Pow(float64(cols), 2.0)))) / 142)

	//return output
	return &Domain{
		Rows:   rows,
		Cols:   cols,
		Matrix: domainMatrix,
		BndCnt: bandCount,
	}
}

// new objective initialization function
func NewObjective(identifier int, fitnessMatrix *mat64.Dense) *Objective {

	// return output
	return &Objective{
		Id:     identifier,
		Matrix: fitnessMatrix,
	}
}

// new basis solution initialization function
func NewBasis(sourceSubs, destinationSubs []int, searchDomain *Domain) *Basis {

	// set default length factor
	var lenFactor int = 10

	// compute all minimum euclidean distances for search domain
	allMinimumDistances := AllMinDistance(sourceSubs, destinationSubs, searchDomain.Matrix)

	// generate subscripts from bresenham's algorithm
	subs := Bresenham(sourceSubs, destinationSubs)

	// compute maximum permitted chromosome length
	maxLength := len(subs) * lenFactor

	// return output
	return &Basis{
		Matrix: allMinimumDistances,
		Subs:   subs,
		MaxLen: maxLength,
	}
}

// new chromosome initialization function
func NewChromosome(searchDomain *Domain, searchParameters *Parameters, searchObjectives *MultiObjective) *Chromosome {

	// initialize floating point parameter values
	var aggFit float64 = 0.0

	// generate node subscripts
	nodeSubs := NewNodeSubs(searchDomain, searchParameters)

	// generate subscripts from directed walk procedure
	subs := MultiPartDirectedWalk(nodeSubs, searchDomain, searchParameters)

	// initialize empty fitness place holders
	fitVal := make([][]float64, searchObjectives.ObjectiveCount)
	for i := 0; i < searchObjectives.ObjectiveCount; i++ {
		fitVal[i] = make([]float64, len(subs))
	}
	totFit := make([]float64, searchObjectives.ObjectiveCount)

	// generate placeholder variables
	uuid := uuid.NewV4()

	// return output
	return &Chromosome{
		Id:               uuid,
		Subs:             subs,
		Fitness:          fitVal,
		TotalFitness:     totFit,
		AggregateFitness: aggFit,
	}
}

// new empty chromosome initialization function
func NewEmptyChromosome(searchDomain *Domain, searchObjectives *MultiObjective) *Chromosome {

	// initialize floating point parameter values
	var aggFit float64 = 0.0

	// initialize subscripts
	subs := make([][]int, 0)

	// generate placeholder id
	uuid := uuid.NewV4()

	// initialize empty fitness place holders
	fitVal := make([][]float64, searchObjectives.ObjectiveCount)
	for i := 0; i < searchObjectives.ObjectiveCount; i++ {
		fitVal[i] = make([]float64, len(subs))
	}
	totFit := make([]float64, searchObjectives.ObjectiveCount)

	// return output
	return &Chromosome{
		Id:               uuid,
		Subs:             subs,
		Fitness:          fitVal,
		TotalFitness:     totFit,
		AggregateFitness: aggFit,
	}
}

// new walker initialization function
func NewWalker(searchDomain *Domain, searchParameters *Parameters, searchObjectives *MultiObjective) Walker {

	// generate uuid
	uuid := uuid.NewV4()

	// create, and return the walker
	walker := Walker{
		Id:               uuid,
		SearchDomain:     searchDomain,
		SearchParameters: searchParameters,
		SearchObjectives: searchObjectives,
	}

	return walker
}

// new population initialization function
func NewPopulation(identifier int, searchDomain *Domain, searchParameters *Parameters, searchObjectives *MultiObjective) *Population {

	// initialize floating point parameter values
	var aggMeanFit float64 = 0.0

	// initialize communication channel
	chr := make(chan *Chromosome, searchParameters.PopSize)

	// initialize walk request channel
	var walkQueue = make(chan bool, searchParameters.PopSize)

	// populate walkqueue channel
	for j := 0; j < searchParameters.PopSize; j++ {
		walkQueue <- true
	}

	// initialize wait group
	var wg sync.WaitGroup

	// generate chromosomes via go routines
	for i := 0; i < searchParameters.ConSize; i++ {
		walker := NewWalker(searchDomain, searchParameters, searchObjectives)
		walker.Start(chr, walkQueue, wg)
	}

	// wait for walkers to finish
	wg.Wait()

	// initialize fitness placeholder
	meanFit := make([]float64, searchObjectives.ObjectiveCount)

	// return output
	return &Population{
		Id:                   identifier,
		Chromosomes:          chr,
		MeanFitness:          meanFit,
		AggregateMeanFitness: aggMeanFit,
	}

}

// new empty population initialization function
func NewEmptyPopulation(identifier int, searchObjectives *MultiObjective) *Population {

	// initialize floating point parameter values
	var aggMeanFit float64 = 0.0

	// initialize empty chromosomes channel
	chr := make(chan *Chromosome)

	// initialize fitness placeholder
	meanFit := make([]float64, searchObjectives.ObjectiveCount)

	// return output
	return &Population{
		Id:                   identifier,
		Chromosomes:          chr,
		MeanFitness:          meanFit,
		AggregateMeanFitness: aggMeanFit,
	}
}

// new empty evolution initialization function
func NewEmptyEvolution(searchParameters *Parameters) *Evolution {

	// initialize empty population channel
	popChan := make(chan *Population, searchParameters.EvoSize)

	// initialize empty fitness gradient
	gradFit := make([]float64, searchParameters.EvoSize)

	// return output
	return &Evolution{
		Populations:     popChan,
		FitnessGradient: gradFit,
	}
}

// new mutator initialization function
func NewMutator(searchDomain *Domain, searchParameters *Parameters, searchObjectives *MultiObjective) Mutator {

	// create, and return the walker
	mutator := Mutator{
		SearchDomain:     searchDomain,
		SearchParameters: searchParameters,
		SearchObjectives: searchObjectives,
	}

	// return output
	return mutator
}

// new evolution initialization function
func NewEvolution(searchParameters *Parameters, searchDomain *Domain, searchObjectives *MultiObjective) *Evolution {

	// initialize seed population identifier
	var popID int = 0

	// initialize population channel
	popChan := make(chan *Population, 1)

	// print initialization status message
	fmt.Println("Initializing Seed Population...")

	// initialize seed population
	seedPop := NewPopulation(popID, searchDomain, searchParameters, searchObjectives)
	popChan <- PopulationFitness(seedPop, searchParameters, searchObjectives)

	// initialize raw fitness data slice
	rawAggMeanFit := make([]float64, searchParameters.EvoSize)

	// initialize fitness gradient variable
	gradFit := make([]float64, searchParameters.EvoSize)

	// enter loop
	for i := 0; i < searchParameters.EvoSize; i++ {

		// perform population evolution
		newPop := PopulationEvolution(<-popChan, searchDomain, searchParameters, searchObjectives)

		// compute population fitness
		newPop = PopulationFitness(newPop, searchParameters, searchObjectives)

		// write aggregate mean fitness value to vector
		rawAggMeanFit[i] = newPop.AggregateMeanFitness

		// generate inline fitness gradient function
		var fitnessGradFnc = func(n float64) float64 { return rawAggMeanFit[int(n)] }

		// compute fitness gradient
		gradFit[i] = fd.Derivative(fitnessGradFnc, float64(i), nil)

		// skip gradient check on first iteration
		if i < 1 {

			// return new population to channel
			popChan <- newPop

			// increment progress
			fmt.Println("Evolution: ", i+1)
			fmt.Printf("Gradient: %f \n", math.Log10(math.Abs(gradFit[i])))
			fmt.Printf("Average Fitness: %f \n", newPop.AggregateMeanFitness)

		} else if i >= 1 && i < (searchParameters.EvoSize-1) {

			if gradFit[i] > 0 {

				// return current population to channel
				popChan <- newPop

				// close population channel
				close(popChan)

				// print success message
				fmt.Println("Convergence Achieved, Evolution Commplete!")

				// break loop
				break

			} else if gradFit[i] <= 0 {

				// return new population to channel
				popChan <- newPop

				// increment progress bar
				fmt.Println("Evolution: ", i+1)
				fmt.Printf("Gradient: %f \n", math.Log10(math.Abs(gradFit[i])))
				fmt.Printf("Average Fitness: %f \n", newPop.AggregateMeanFitness)

			}

		} else if i == searchParameters.EvoSize-1 {

			// return new population to channel
			popChan <- newPop

			// close population channel
			close(popChan)

			// print termination message
			fmt.Println("Convergence Not Achieved, Maximum Number of Evolutions Reached...")
			fmt.Printf("Gradient: %f \n", math.Log10(math.Abs(gradFit[i])))
			fmt.Printf("Average Fitness: %f \n", newPop.AggregateMeanFitness)

			// break loop
			break
		}
	}

	// return output
	return &Evolution{
		Populations:     popChan,
		FitnessGradient: gradFit,
	}
}

/* function to return copies of a user specified fraction of
the individual chromosomes within a population ranked in terms
of individual aggregate fitness */
func NewEliteFraction(inputFraction float64, inputPopulation *Population) (outputChromosomes []*Chromosome) {

	// count input chromosomes
	chromCount := cap(inputPopulation.Chromosomes)

	// initialize aggregate score slice
	chromFrac := int(math.Ceil(inputFraction * float64(chromCount)))

	// initialize chromosome map
	chromMap := make(map[float64]*Chromosome)

	// initialize chromosome map key slice
	chromKey := make([]float64, chromCount)

	// initialize output slice
	output := make([]*Chromosome, chromFrac)

	// loop through channel to populate slice
	for i := 0; i < chromCount; i++ {
		curChrom := <-inputPopulation.Chromosomes
		chromMap[curChrom.AggregateFitness] = curChrom
		chromKey[i] = curChrom.AggregateFitness
	}

	// sort on aggregate fitness keys
	sort.Float64s(chromKey)

	// loop through and generate output slice faction
	for j := 0; j < chromFrac; j++ {
		output[j] = chromMap[chromKey[j]]
	}

	// return output
	return output
}

/* function to return copies of a user specified number of
unique individual chromosomes from within a population
with each chromosome being ranked in terms of its
individual aggregate fitness */
func NewEliteSet(inputCount int, inputPopulation *Population, inputParameters *Parameters) (outputChromosomes []*Chromosome) {

	// check band count against population size
	if inputCount >= int(math.Floor((0.5 * float64(inputParameters.PopSize)))) {
		err := errors.New("Input elite set count must be less than 1/2 the input population size \n")
		panic(err)
	}

	// count input chromosomes
	chromCount := cap(inputPopulation.Chromosomes)

	// initialize chromosome map
	chromMap := make(map[float64]*Chromosome)

	// initialize chromosome map key slice
	chromKey := make([]float64, chromCount)

	// initialize output slice
	output := make([]*Chromosome, inputCount)

	// loop through channel to populate slice from channel
	for i := 0; i < chromCount; i++ {
		curChrom := <-inputPopulation.Chromosomes
		chromMap[curChrom.AggregateFitness] = curChrom
		chromKey[i] = curChrom.AggregateFitness
		inputPopulation.Chromosomes <- curChrom
	}

	// sort on aggregate fitness keys
	sort.Float64s(chromKey)

	// initalize iteration counter
	var iter int = 0

	// loop through and generate output slice set
	for j := 0; j < chromCount; j++ {

		// deal with initial state
		if j == 0 {
			output[iter] = chromMap[chromKey[j]]
			iter += 1
			continue
		}

		// get uuids
		prevUuid := chromMap[chromKey[j-1]].Id.String()
		curUuid := chromMap[chromKey[j]].Id.String()

		// impose uniqueness constraint
		if prevUuid != curUuid {
			output[iter] = chromMap[chromKey[j]]
			iter += 1
		} else {
			continue
		}

		// stop if inputCount reached
		if iter == inputCount {
			break
		}
	}

	// return output
	return output
}
