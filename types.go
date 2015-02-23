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
	RndCoef int
	PopSize int
}

// new problem parameters function
func NewParameters(sourceSubscripts, destinationSubscripts []int, randomnessCoefficient, populationSize int) *Parameters {

	// return output
	return &Parameters{
		SrcSub:  sourceSubscripts,
		DstSub:  destinationSubscripts,
		RndCoef: randomnessCoefficient,
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
	Id     int
	Matrix *mat64.Dense
}

// new objective initialization function
func NewObjective(identifier int, fitnessMatrix *mat64.Dense) *Objective {

	// return output
	return &Objective{
		Id:     identifier,
		Matrix: fitnessMatrix,
	}
}

// individuals are comprised of row column indices to some
// spatially reference search domain.
type Individual struct {
	Id           *uuid.UUID
	Subs         [][]int
	Fitness      []float64
	TotalFitness float64
}

// new individual initialization function
func NewIndividual(searchDomain *Domain, searchParameters *Parameters, searchObjective *Objective) *Individual {

	// generate subscripts from directed walk procedure
	sub := Dirwlk(searchParameters, searchDomain)

	// evaluate fitness for subscripts
	fitVal, totFit := Fitness(sub, searchObjective.Matrix)

	// generate placeholder variables
	uuid, _ := uuid.NewV4()

	// return output
	return &Individual{
		Id:           uuid,
		Subs:         sub,
		Fitness:      fitVal,
		TotalFitness: totFit,
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
