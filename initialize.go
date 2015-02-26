// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"github.com/gonum/matrix/mat64"
	"github.com/nu7hatch/gouuid"
)

// new problem parameters function
func NewParameters(sourceSubscripts, destinationSubscripts []int, randomnessCoefficient float64, populationSize int, selectionFraction, selectionProbability float64) *Parameters {

	// return output
	return &Parameters{
		SrcSubs: sourceSubscripts,
		DstSubs: destinationSubscripts,
		RndCoef: randomnessCoefficient,
		PopSize: populationSize,
		SelProb: selectionProbability,
		SelFrac: selectionFraction,
	}
}

// new domain initialization function
func NewDomain(identifier int, domainMatrix *mat64.Dense) *Domain {

	//return output
	return &Domain{
		Id:     identifier,
		Matrix: domainMatrix,
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

	// generate subscripts from directed walk procedure
	subs := Dirwlk(searchDomain, searchParameters, basisSolution)

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

// new population initialization function
func NewPopulation(identifier int, searchDomain *Domain, searchParameters *Parameters, searchObjective *Objective, basisSolution *Basis) *Population {

	// initialize communication channel
	chr := make(chan *Chromosome, searchParameters.PopSize)

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
