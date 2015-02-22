// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"github.com/gonum/matrix/mat64"
	"github.com/nu7hatch/gouuid"
)

// parameters are comprised of fixed input avlues that are
// unique to the problem specification that are referenced
// by the algorithm at various stage of the solution process
type Parameters struct {
	SrcInd  int
	DstInd  int
	PopSize int
}

// new problem parameters function
func NewParameters(sourceIndex, destinationIndex, populationSize int) *Parameters {

	// return output
	return &Parameters{
		SrcInd:  sourceIndex,
		DstInd:  destinationIndex,
		PopSize: populationSize,
	}
}

// domains are comprised of boolean arrays which indicate the
// feasible locations for the search algorithm
type Domain struct {
	Id     int
	Size   int
	Stride int
	Vals   []bool
}

// new domain initialization function
func NewDomain(identifier, domainSize, domainStride int, domainValues []bool) *Domain {

	//return output
	return &Domain{
		Id:     identifier,
		Size:   domainSize,
		Stride: domainStride,
		Vals:   domainValues,
	}
}

// objectives are comprised of maps which use location indices
// to key to floating point fitness values within the search
// domain
type Objective struct {
	Id      int
	Fitness []float64
}

// new objective initialization function
func NewObjective(identifier int, fitnessValues []float64) *Objective {

	// return output
	return &Objective{
		Id:      identifier,
		Fitness: fitnessValues,
	}
}

// individuals are comprised of row column indices to some
// spatially reference search domain.
type Individual struct {
	Id          *uuid.UUID
	Indices     []int
	Fitness     []float64
	MeanFitness float64
}

// new individual initialization function

func NewIndividual(searchDomain *Domain, problemParameters *Parameters) *Individual {

	// initialize iterator and output variables
	i := 1
	maxLen := 100
	ind := make([]int, 1, maxLen)
	ind[0] = problemParameters.SrcInd

	// initialize mu and sigma
	muVec := make([]float64, 2)
	sigmaVec := make([]float64, 4)

	// set mu elements
	muVec[0] = 1
	muVec[1] = 1

	// set sigma elements
	sigmaVec[0] = 1
	sigmaVec[1] = 0
	sigmaVec[2] = 0
	sigmaVec[3] = 1

	// generate dense matrices
	mu := mat64.NewDense(1, 2, muVec)
	sigma := mat64.NewDense(2, 2, sigmaVec)
	var try int

	for {
		try = Newind(ind[len(ind)-1], mu, sigma, searchDomain)
		if i == maxLen-1 {
			break
		} else if try == problemParameters.DstInd {
			ind = append(ind, try)
			break
		} else {
			ind = append(ind, try)
			i += 1
		}
	}

	// FOR NOW I AM JUST WRITING SOME PLACE HOLDER VALUES HERE BUT THESE
	// WILL BE REPLACED BY FITNESS EVALUATIONS IN THE FUTURE
	uuid, _ := uuid.NewV4()
	fit := make([]float64, len(ind))
	var meanfit float64
	meanfit = 0.0

	return &Individual{
		Id:          uuid,
		Indices:     ind,
		Fitness:     fit,
		MeanFitness: meanfit,
	}
}

// populations are comprised of a fixed number of individuals.
// this number corresponds to the populationSize.
type Population struct {
	Id          int
	Individuals *[]Individual
	MeanFitness float64
}

// evolutions are comprised of a stochastic number of populations.
// this number is determined by the convergence rate of the
// algorithm.
type Evolution struct {
	Id          int
	Populations *[]Population
}
