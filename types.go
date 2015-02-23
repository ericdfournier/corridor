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
	SrcSub  []int
	DstSub  []int
	PopSize int
}

// new problem parameters function
func NewParameters(sourceSubscripts, destinationSubscripts []int, populationSize int) *Parameters {

	// return output
	return &Parameters{
		SrcSub:  sourceSubscripts,
		DstSub:  destinationSubscripts,
		PopSize: populationSize,
	}
}

// domains are comprised of boolean arrays which indicate the
// feasible locations for the search algorithm
type Domain struct {
	Id     int
	Matrix *mat64.Dense
}

// new domain initialization function
func NewDomain(identifier int, domainMatrix *mat64.Dense) *Domain {

	//return output
	return &Domain{
		Id:     identifier,
		Matrix: domainMatrix,
	}
}

// objectives are comprised of maps which use location indices
// to key to floating point fitness values within the search
// domain
type Objective struct {
	Id      int
	Fitness *mat64.Dense
}

// new objective initialization function
func NewObjective(identifier int, fitnessMatrix *mat64.Dense) *Objective {

	// return output
	return &Objective{
		Id:      identifier,
		Fitness: fitnessMatrix,
	}
}

// individuals are comprised of row column indices to some
// spatially reference search domain.
type Individual struct {
	Id          *uuid.UUID
	Subs        [][]int
	Fitness     []float64
	MeanFitness float64
}

// new individual initialization function
func NewIndividual(searchDomain *Domain, searchParameters *Parameters) *Individual {

	// initialize iterator and output variables
	i := 1
	maxLen := 100
	sub := make([][]int, 1, maxLen)
	sub[0] = make([]int, 2)
	sub[0][0] = searchParameters.SrcSub[0]
	sub[0][1] = searchParameters.SrcSub[1]

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

	// NEED TO MAKE MU VECTOR AND SIGMA RESPOND TO RELATIVE ORIENTATION OF
	// THE DESTINATION NODE

	// generate dense matrices
	mu := mat64.NewDense(1, 2, muVec)
	sigma := mat64.NewDense(2, 2, sigmaVec)
	var try []int

	// enter unbounded for loop
	for {
		cS := sub[len(sub)-1]
		try = Newind(cS, mu, sigma, searchDomain)
		if i == maxLen-1 {
			break
		} else if try[0] == searchParameters.DstSub[0] && try[1] == searchParameters.DstSub[1] {
			sub = append(sub, try)
			break
		} else {
			sub = append(sub, try)
			i += 1
		}
	}

	uuid, _ := uuid.NewV4()
	fit := make([]float64, len(sub))
	var meanfit float64
	meanfit = 0.0

	return &Individual{
		Id:          uuid,
		Subs:        sub,
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
