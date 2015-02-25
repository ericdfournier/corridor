// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"fmt"

	"github.com/cheggaaa/pb"
	"github.com/gonum/matrix/mat64"
	"github.com/nu7hatch/gouuid"
)

// new problem parameters function
func NewParameters(sourceSubscripts, destinationSubscripts []int, randomnessCoefficient, populationSize int) *Parameters {

	// return output
	return &Parameters{
		SrcSubs: sourceSubscripts,
		DstSubs: destinationSubscripts,
		RndCoef: randomnessCoefficient,
		PopSize: populationSize,
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
func NewChromosome(searchDomain *Domain, searchParameters *Parameters, searchObjective *Objective, basisSolution *Basis) *Chromosome {

	// generate subscripts from directed walk procedure
	subs := Dirwlk(searchDomain, searchParameters, basisSolution)

	// evaluate fitness for subscripts
	fitVal, totFit := Fitness(subs, searchObjective.Matrix)

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
func NewPopulation(searchDomain *Domain, searchParameters *Parameters, searchObjective *Objective, basisSolution *Basis) *Population {

	// print start
	fmt.Println("Initializing Population")

	// initialize slice of structs
	chroms := make([]Chromosome, searchParameters.PopSize)

	// initialize communication channel
	chr := make(chan Chromosome)

	// initialize new empty chromosome before entering loop
	emptyChrom := NewChromosome(searchDomain, searchParameters, searchObjective, basisSolution)
	var cumFit float64 = 0.0

	//// initialize progress bar
	bar := pb.StartNew(searchParameters.PopSize)

	// generate chromosomes via go routines
	for i := 0; i < searchParameters.PopSize; i++ {

		//increment bar
		bar.Increment()

		// get new emptyChrom
		emptyChrom = &chroms[i]

		// launch go routines
		go func(emptyChrom *Chromosome) {
			emptyChrom = NewChromosome(searchDomain, searchParameters, searchObjective, basisSolution)
			chr <- *emptyChrom
		}(emptyChrom)

		// read from channel
		newChrom := <-chr
		chroms[i] = newChrom
	}

	// close progress bar
	bar.FinishPrint("Finished")

	// generate mean fitness
	meanFit := cumFit / float64(searchParameters.PopSize)

	// generate placeholder variables
	var identifier int = 1

	// return output
	return &Population{
		Id:          identifier,
		Chromosomes: &chroms,
		MeanFitness: meanFit,
	}

}
