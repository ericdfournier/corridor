/* Copyright ©2015 The corridor Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file. */

package corridor

import (
	"github.com/gonum/matrix/mat64"
	"github.com/satori/go.uuid"
)

/* parameters are comprised of fixed input avlues that are
unique to the problem specification that are referenced
by the algorithm at various stage of the solution process */
type Parameters struct {
	SrcSubs []int   // source subscripts
	DstSubs []int   // destination subscripts
	RndCoef float64 // randomness coefficient
	PopSize int     // population size
	SelFrac float64 // selection fraction
	SelProb float64 // selection probability
	MutaCnt int     // muation count
	MutaFrc float64 // muation fraction
	EvoSize int     // evolution size
	ConSize int     // concurrency limit
}

/* domains are comprised of boolean arrays which indicate the
feasible locations for the search algorithm */
type Domain struct {
	Rows   int          // row count
	Cols   int          // column count
	Matrix *mat64.Dense // domain matrix values
	BndCnt int          // distance band count
}

/* objectives are comprised of matrices which use location
indices to key to floating point fitness values within the
search domain */
type Objective struct {
	Id     int          // objective identification number
	Matrix *mat64.Dense // objective matrix values
}

/* multiObjective objects are comprised of a channel of individual
independent objectives that are used for the evaluation of
chromosome and population level fitness values */
type MultiObjective struct {
	ObjectiveCount int          // objective count
	Objectives     []*Objective // individual objective objects
}

/* a basis solution is comprised of the subscript indices forming
the euclidean shortest path connecting the source to the dest */
type Basis struct {
	Matrix *mat64.Dense // basis matrix values
	Subs   [][]int      // basis row column subscripts
	MaxLen int          // maximum length
}

/* chromosomess are comprised of genes which are distinct row column
indices to some spatially reference search domain */
type Chromosome struct {
	Id               uuid.UUID   // globally unique chromosome identification number
	Subs             [][]int     // chromosome row column subscripts
	Fitness          [][]float64 // objective function values
	TotalFitness     []float64   // total fitness values for each objective
	AggregateFitness float64     // total aggregate fitness value for all objectives
}

/* populations are comprised of a fixed number of chromosomes.
this number corresponds to the populationSize. */
type Population struct {
	Id                   int              // population ordinal identification number
	Chromosomes          chan *Chromosome // chromosome channel
	MeanFitness          []float64        // chromosome mean fitnesses channel
	AggregateMeanFitness float64          // chromosome aggregate fitness channel
}

/* evolutions are comprised of a stochastic number of populations.
this number is determined by the convergence rate of the
algorithm */
type Evolution struct {
	Populations     chan *Population // population channel
	FitnessGradient []float64        // fitness gradient values
}

/*  walkers are used in the concurrency model which facilitates
the parallel problem initializations */
type Walker struct {
	Id               uuid.UUID       // globally unique walker identification number
	SearchDomain     *Domain         // local search domain copy
	SearchParameters *Parameters     // local search parameters copy
	SearchObjectives *MultiObjective // local search objectives copy
}

/* mutators are used to generate point location mutations in
parallel for a subset of chromosomes within a population */
type Mutator struct {
	SearchDomain     *Domain         // local search domain copy
	SearchParameters *Parameters     // local search parameters copy
	SearchObjectives *MultiObjective // local search objectives copy
}
