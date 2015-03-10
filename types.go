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
	SrcSubs []int
	DstSubs []int
	RndCoef float64
	PopSize int
	SelFrac float64
	SelProb float64
	EvoSize int
}

// domains are comprised of boolean arrays which indicate the
// feasible locations for the search algorithm
type Domain struct {
	Id     int
	Matrix *mat64.Dense
	MaxLen int
}

// objectives are comprised of maps which use location indices
// to key to floating point fitness values within the search
// domain
type Objective struct {
	Id     int
	Matrix *mat64.Dense
}

// a basis solution is comprised of the subscript indices forming
// the euclidean shortest path connecting the source to the dest
type Basis struct {
	Id     int
	Matrix *mat64.Dense
	Subs   [][]int
}

// chromosomess are comprised of genes which are distinct row column
// indices to some spatially reference search domain.
type Chromosome struct {
	Id           *uuid.UUID
	Subs         [][]int
	Fitness      []float64
	TotalFitness float64
}

// populations are comprised of a fixed number of chromosomes.
// this number corresponds to the populationSize.
type Population struct {
	Id          int
	Chromosomes chan *Chromosome
	MeanFitness float64
}

// evolutions are comprised of a stochastic number of populations.
// this number is determined by the convergence rate of the
// algorithm.
type Evolution struct {
	Id              *uuid.UUID
	Populations     chan *Population
	FitnessGradient float64
}
