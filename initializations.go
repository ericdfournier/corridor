// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

// new problem parameters function
func NewParameters() {
	
	return &Parameters{
		srcInd:
		dstInd:
		srcSub:
		dstSub:
		popSize:
	}
}

// new domain initialization function
func NewDomain(identifier, domainSize, domainStride int, domainValues []bool) 
	*Domain {

	//return output
	return &Domain{
		id:     identifier,
		size:   domainSize,
		stride: domainStride,
		vals:   domainValues,
	}
}

// new objective initialization function
func NewObjective(identifier int, fitnessValues []float64) *Objective {

	// return output
	return &Objective{
		id:      identifier,
		fitness: fitnessValues,
	}
}

/////////////////////////////////////////////////////////////////////////////
// new individual initialization function
//func NewIndividual(problemParameters Parameters, problemDomain Domain) 
//	*Individual {
//
//}
