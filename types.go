// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"github.com/cheggaaa/pb"
	"github.com/gonum/matrix/mat64"
	"github.com/nu7hatch/gouuid"
)

// parameters are comprised of fixed input avlues that are
// unique to the problem specification that are referenced
// by the algorithm at various stage of the solution process
type Parameters struct {
	SrcSubs []int
	DstSubs []int
	RndCoef int
	PopSize int
}

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

// a basis solution is comprised of the subscript indices forming
// the euclidean shortest path connecting the source subscript
// to the destination subscript as well as information regarding
// the minimum euclidean distances of all locations within the
// search domain to the nearest point on this euclidean shortest
// path
type Basis struct {
	Id     int
	Matrix *mat64.Dense
	Subs   [][]int
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

// individuals are comprised of row column indices to some
// spatially reference search domain.
type Individual struct {
	Id           *uuid.UUID
	Subs         [][]int
	Fitness      []float64
	TotalFitness float64
}

// new individual initialization function
func NewIndividual(searchDomain *Domain, searchParameters *Parameters, searchObjective *Objective, basisSolution *Basis) *Individual {

	// generate subscripts from directed walk procedure
	subs := Dirwlk(searchParameters, searchDomain, basisSolution)

	// evaluate fitness for subscripts
	fitVal, totFit := Fitness(subs, searchObjective.Matrix)

	// generate placeholder variables
	uuid, _ := uuid.NewV4()

	// return output
	return &Individual{
		Id:           uuid,
		Subs:         subs,
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

// new population initialization function
func NewPopulation(searchDomain *Domain, searchParameters *Parameters, searchObjective *Objective, basisSolution *Basis) *Population {

	// initialize slice of structs
	indiv := make([]Individual, searchParameters.PopSize)
	var cumFit float64 = 0.0

	// initialize progress bar
	bar := pb.StartNew(searchParameters.PopSize)

	// generate individuals
	for i := 0; i < searchParameters.PopSize; i++ {
		bar.Increment()
		ind := NewIndividual(searchDomain, searchParameters, searchObjective, basisSolution)
		indiv[i] = *ind
		cumFit = cumFit + indiv[i].TotalFitness
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
		Individuals: &indiv,
		MeanFitness: meanFit,
	}

}

// evolutions are comprised of a stochastic number of populations.
// this number is determined by the convergence rate of the
// algorithm.
type Evolution struct {
	Id          int
	Populations *[]Population
}
