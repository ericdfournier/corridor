// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

// domains are comprised of maps which use location indices
// to key to boolean values of feasible search domain
type Domain struct {
	id   int
	size int
	vals []bool
}

// new domain initialization function
func NewDomain(identifier, domainSize int, domainValues []bool) *Domain {

	//return output
	return &Domain{
		id:   identifier,
		size: domainSize,
		vals: domainValues,
	}
}

// new test domain initialization function
func NewTestDomain(identifier, domainSize, domainStride int) *Domain {

	// initialize value slice
	domainValues := make([]bool, domainSize)

	// loop through index values define domain
	for i := 0; i < domainSize; i++ {
		if i >= 0 && i < domainStride {
			domainValues[i] = false
		} else if i > (domainSize - domainStride) {
			domainValues[i] = false
		} else {
			domainValues[i] = true
		}
	}
	for i := 0; i < domainSize; i = (i + domainStride) {
		domainValues[i] = false
	}
	for i := 1; i < domainSize; i = (i + domainStride) {
		domainValues[i] = false
	}

	// return output
	return &Domain{
		id:   identifier,
		size: domainSize,
		vals: domainValues,
	}
}

// objectives are comprised of maps which use location indices
// to key to floating point fitness values within the search
// domain
type Objective struct {
	id      int
	size    int
	fitness []float64
}

// new objective initialization function
func NewObjective(identifier, objectiveSize int, fitnessValues []float64) *Objective {

	// return output
	return &Objective{
		id:      identifier,
		size:    objectiveSize,
		fitness: fitnessValues,
	}
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
