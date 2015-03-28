// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"math"

	"github.com/gonum/matrix/mat64"
	"github.com/nu7hatch/gouuid"
)

// new problem parameters function
func NewParameters(sourceSubscripts, destinationSubscripts []int, randomnessCoefficient float64, populationSize, evolutionSize int, selectionFraction, selectionProbability float64) *Parameters {

	// return output
	return &Parameters{
		SrcSubs: sourceSubscripts,
		DstSubs: destinationSubscripts,
		RndCoef: randomnessCoefficient,
		PopSize: populationSize,
		SelProb: selectionProbability,
		SelFrac: selectionFraction,
		EvoSize: evolutionSize,
	}
}

// new domain initialization function
func NewDomain(identifier int, domainMatrix *mat64.Dense) *Domain {

	// get domain size
	rows, cols := domainMatrix.Dims()

	// initialize fixed parameters
	var p float64 = 2
	var s float64 = 5

	// compute maximum permitted chromosome length
	maximumLength := int(math.Ceil(s * math.Sqrt(math.Pow(float64(rows), p)+math.Pow(float64(cols), p))))

	//return output
	return &Domain{
		Id:     identifier,
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

func NewBasis(searchDomain *Domain, searchParameters *Parameters) *Basis {

	// compute all minimum euclidean distances for search domain
	allMinimumDistances := AllMinDistance(searchParameters.SrcSubs, searchParameters.DstSubs, searchDomain.Matrix)

	// generate subscripts from bresenham's algorithm
	subs := Bresenham(searchParameters.SrcSubs, searchParameters.DstSubs)

	// return output
	return &Basis{
		Id:     searchDomain.Id,
		Matrix: allMinimumDistances,
		Subs:   subs,
	}
}

// new chromosome initialization function
func NewChromosome(searchDomain *Domain, searchParameters *Parameters, basisSolution *Basis) *Chromosome {

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
	fitVal := make([]float64, len(subs))
	var totFit float64 = 0.0

	// generate placeholder variables
	uuid, _ := uuid.NewV4()

	// return output
	return &Chromosome{
		Id:           uuid,
		Subs:         subs,
		Fitness:      fitVal,
		TotalFitness: totFit,
	}
}

// new empty chromosome initialization function
func NewEmptyChromosome(searchDomain *Domain) *Chromosome {

	// initialize subscripts
	subs := make([][]int, searchDomain.MaxLen)

	// generate placeholder id
	uuid, _ := uuid.NewV4()

	// initialize empty fitness place holders
	fitVal := make([]float64, searchDomain.MaxLen)
	var totFit float64 = 0.0

	// return output
	return &Chromosome{
		Id:           uuid,
		Subs:         subs,
		Fitness:      fitVal,
		TotalFitness: totFit,
	}
}

// new population initialization function
func NewPopulation(identifier int, searchDomain *Domain, searchParameters *Parameters, searchObjective *Objective) *Population {

	// initialize communication channel
	chr := make(chan *Chromosome, searchParameters.PopSize)

	// generate basis solution
	basisSolution := NewBasis(searchDomain, searchParameters)

	// initialize new empty chromosome before entering loop
	emptyChrom := NewChromosome(searchDomain, searchParameters, basisSolution)

	// generate chromosomes via go routines
	for i := 0; i < searchParameters.PopSize; i++ {

		// launch chromosome initialization go routines
		go func(searchDomain *Domain, searchParameters *Parameters, searchObjective *Objective, basisSolution *Basis) {
			emptyChrom = NewChromosome(searchDomain, searchParameters, basisSolution)
			chr <- ChromosomeFitness(emptyChrom, searchObjective)
		}(searchDomain, searchParameters, searchObjective, basisSolution)

	}

	// initialize mean fitness
	var cumFit float64 = 0.0

	// generate mean fitness
	meanFit := cumFit / float64(searchParameters.PopSize)

	// return output
	return &Population{
		Id:          identifier,
		Chromosomes: chr,
		MeanFitness: meanFit,
	}

}

// new empty population initialization function
func NewEmptyPopulation(identifier int) *Population {

	// initialize empty chromosomes channel
	chr := make(chan *Chromosome)

	// initialize mean fitness value
	var meanFit float64 = 0.0

	// return output
	return &Population{
		Id:          identifier,
		Chromosomes: chr,
		MeanFitness: meanFit,
	}
}

// new empty evolution initialization function
func NewEmptyEvolution(searchParameters *Parameters) *Evolution {

	// generate evolution id
	uuid, _ := uuid.NewV4()

	// initialize empty population channel
	popChan := make(chan *Population, searchParameters.EvoSize)

	// initialize empty fitness gradient
	gradFit := make([]float64, searchParameters.EvoSize)

	// return output
	return &Evolution{
		Id:              uuid,
		Populations:     popChan,
		FitnessGradient: gradFit,
	}
}

// new evolution initialization function
func NewEvolution(searchParameters *Parameters, searchDomain *Domain, searchObjective *Objective) *Evolution {

	// generate evolution id
	uuid, _ := uuid.NewV4()

	// initialize seed population identifier
	var popID int = 0

	// initialize seed population
	popChan := make(chan *Population, searchParameters.EvoSize)

	// initiali
	curPop := NewPopulation(popID, searchDomain, searchParameters, searchObjective)
	curPop = PopulationFitness(curPop, searchParameters, searchObjective)

	// enter loop
	for i := 0; i < searchParameters.EvoSize; i++ {
		newPop := PopulationEvolution(curPop, searchDomain, searchParameters, searchObjective)
		newPop = PopulationFitness(newPop, searchParameters, searchObjective)
		curPop = newPop
		popChan <- newPop
	}

	// PLACEHOLDER FITNESS GRADIENT
	gradFit := make([]float64, searchParameters.EvoSize)

	// evaluate seed population
	// return output
	return &Evolution{
		Id:              uuid,
		Populations:     popChan,
		FitnessGradient: gradFit,
	}
}
