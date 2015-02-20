// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

// parameters are comprised of fixed input avlues that are
// unique to the problem specification that are referenced
// by the algorithm at various stage of the solution process
type Parameters struct {
	srcInd  int
	dstInd  int
	srcSub  [2]int
	dstSub  [2]int
	popSize int
}

// domains are comprised of boolean arrays which indicate the
// feasible locations for the search algorithm
type Domain struct {
	id     int
	size   int
	stride int
	vals   []bool
}

// objectives are comprised of maps which use location indices
// to key to floating point fitness values within the search
// domain
type Objective struct {
	id      int
	size    int
	fitness []float64
}

// individuals are comprised of row column indices to some
// spatially reference search domain.
type Individual struct {
	id          int
	subscripts  [][]int
	indices     []int
	fitness     []float64
	meanFitness float64
}

// populations are comprised of a fixed number of individuals.
// this number corresponds to the populationSize.
type Population struct {
	id          int
	individuals *[]Individual
	meanFitness float64
}

// evolutions are comprised of a stochastic number of populations.
// this number is determined by the convergence rate of the
// algorithm.
type Evolution struct {
	id          int
	populations *[]Population
}
