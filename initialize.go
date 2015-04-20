// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"fmt"
	"math"

	"github.com/gonum/diff/fd"
	"github.com/gonum/matrix/mat64"
	"github.com/satori/go.uuid"
)

// new problem parameters function
func NewParameters(sourceSubscripts, destinationSubscripts []int, populationSize, evolutionSize int, randomnessCoefficient float64) *Parameters {

	// set default mutation count
	mutationCount := 1

	// set defacult mutation fraction
	mutationFraction := 0.2

	// set selection fraction
	selectionFraction := 0.5

	// set selection probability
	selectionProbability := 0.8

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
	}
}

// new domain initialization function
func NewDomain(domainMatrix *mat64.Dense) *Domain {

	// get domain size
	rows, cols := domainMatrix.Dims()

	// compute maximum permitted chromosome length
	maximumLength := 10 * int(math.Floor(math.Sqrt(domainMatrix.Sum())))

	//return output
	return &Domain{
		Rows:   rows,
		Cols:   cols,
		Matrix: domainMatrix,
		MaxLen: maximumLength,
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
func NewBasis(searchDomain *Domain, searchParameters *Parameters) *Basis {

	// compute all minimum euclidean distances for search domain
	allMinimumDistances := AllMinDistance(searchParameters.SrcSubs, searchParameters.DstSubs, searchDomain.Matrix)

	// generate subscripts from bresenham's algorithm
	subs := Bresenham(searchParameters.SrcSubs, searchParameters.DstSubs)

	// initialize convexity count at zero
	var convexSum int = 0

	// intialize convexity boolean as true
	var convexBool bool = true

	// loop through subs to evaluate if the search domain is exited
	for i := 0; i < len(subs); i++ {
		if searchDomain.Matrix.At(subs[i][0], subs[0][1]) == 0 {
			convexSum += 1
		}
	}

	// if domain is exited flip boolean
	if convexSum > 0 {
		convexBool = false
	}

	// return output
	return &Basis{
		Matrix: allMinimumDistances,
		Subs:   subs,
		Convex: convexBool,
	}
}

// new chromosome initialization function
func NewChromosome(searchDomain *Domain, searchParameters *Parameters, searchObjectives *MultiObjective, basisSolution *Basis) *Chromosome {

	// initialize variables
	var subs [][]int
	var dstTest bool

	// enter unbounded for loop
	for {
		// generate subscripts from directed walk procedure
		subs, dstTest = DirWlk(searchDomain, searchParameters, basisSolution)

		// regenerate walk if destination not met within maximum chromosome length
		if dstTest == false {
			continue
		} else {
			break
		}
	}

	// initialize empty fitness place holders
	fitVal := make([][]float64, searchObjectives.ObjectiveCount)
	for i := 0; i < searchObjectives.ObjectiveCount; i++ {
		fitVal[i] = make([]float64, len(subs))
	}
	totFit := make([]float64, searchObjectives.ObjectiveCount)
	var aggFit float64 = 0.0

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

	// initialize subscripts
	subs := make([][]int, searchDomain.MaxLen)

	// generate placeholder id
	uuid := uuid.NewV4()

	// initialize empty fitness place holders
	fitVal := make([][]float64, searchObjectives.ObjectiveCount)
	for i := 0; i < searchObjectives.ObjectiveCount; i++ {
		fitVal[i] = make([]float64, len(subs))
	}
	totFit := make([]float64, searchObjectives.ObjectiveCount)
	var aggFit float64 = 0.0

	// return output
	return &Chromosome{
		Id:               uuid,
		Subs:             subs,
		Fitness:          fitVal,
		TotalFitness:     totFit,
		AggregateFitness: aggFit,
	}
}

// new population initialization function
func NewPopulation(identifier int, searchDomain *Domain, searchParameters *Parameters, searchObjectives *MultiObjective) *Population {

	// initialize communication channel
	chr := make(chan *Chromosome, searchParameters.PopSize)

	// generate basis solution
	basisSolution := NewBasis(searchDomain, searchParameters)

	// initialize new empty chromosome before entering loop
	emptyChrom := NewChromosome(searchDomain, searchParameters, searchObjectives, basisSolution)

	// generate chromosomes via go routines
	for i := 0; i < searchParameters.PopSize; i++ {

		// launch chromosome initialization go routines
		go func(searchDomain *Domain, searchParameters *Parameters, searchObjectives *MultiObjective, basisSolution *Basis) {
			emptyChrom = NewChromosome(searchDomain, searchParameters, searchObjectives, basisSolution)
			chr <- ChromosomeFitness(emptyChrom, searchObjectives)
		}(searchDomain, searchParameters, searchObjectives, basisSolution)

	}

	// initialize fitness placeholder
	meanFit := make([]float64, searchObjectives.ObjectiveCount)
	var aggMeanFit float64 = 0.0

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

	// initialize empty chromosomes channel
	chr := make(chan *Chromosome)

	// initialize fitness placeholder
	meanFit := make([]float64, searchObjectives.ObjectiveCount)
	var aggMeanFit float64 = 0.0

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

			// increment progress bar
			fmt.Println("Evolution: ", i+1)

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

			}

		} else if i == searchParameters.EvoSize-1 {

			// return new population to channel
			popChan <- newPop

			// close population channel
			close(popChan)

			// print termination message
			fmt.Println("Convergence Not Achieved, Maximum Number of Evolutions Reached...")

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

//// THIS ROUTINE IS STILL UNDER DEVELOPMENT...

//// function to return copies of a user specified fraction of
//// the individual chromosomes within a population ranked in terms
//// of individual aggregate fitness
//func NewElites(inputFraction float64, inputPopulation *Population) (outputChromosomes []*Chromosomes) {

//	// initialize aggregate score slice
//	count := math.Floor(inputFraction * cap(inputPopulation.Chromosomes))

//	// initialize map
//	chromMap := make(map[uuid.UUID]float64)

//	// initialize

//}
