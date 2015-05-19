// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/gonum/matrix/mat64"
)

// multivariatenormalrandom generates pairs of bivariate normally distributed
// random numbers given an input mean vector and covariance matrix
func MultiVariateNormalRandom(mu *mat64.Dense, sigma *mat64.SymDense) (rndsmp *mat64.Dense) {

	// initialize vector slices
	o := make([]float64, 2)
	n := make([]float64, 2)

	// generate random numbers from normal distribution, prohibit [0,0]
	// combinations
	rand.Seed(time.Now().UnixNano())

	// enter loop
	for i := 0; i < 2; i++ {
		n[i] = rand.NormFloat64()
	}

	// convert to matrix type
	rnd := mat64.NewDense(2, 1, n)
	output := mat64.NewDense(2, 1, o)

	// perform cholesky decomposition on covariance matrix
	lower := mat64.NewTriDense(2, false, nil)
	lower.Cholesky(sigma, true)

	// compute output
	output.Mul(lower, rnd)
	output.Add(output, mu)

	//return final output
	return output
}

// fixrandom converts an input vector of bivariate normally distributed random
// numbers into a version where the values have been fixed to a [-1, 0 ,1]
// range
func FixRandom(rndsmp *mat64.Dense) (fixsmp *mat64.Dense) {

	// initialize vector slice
	o := make([]float64, 2)

	// write up down movement direction
	if rndsmp.At(0, 0) > 1.0 {
		o[0] = 1
	} else if rndsmp.At(0, 0) >= -1.0 && rndsmp.At(0, 0) <= 1.0 {
		o[0] = 0
	} else if rndsmp.At(0, 0) < -1.0 {
		o[0] = -1
	}

	// write left right movement direction
	if rndsmp.At(1, 0) > 1.0 {
		o[1] = 1
	} else if rndsmp.At(1, 0) >= -1.0 && rndsmp.At(1, 0) <= 1.0 {
		o[1] = 0
	} else if rndsmp.At(1, 0) < -1.0 {
		o[1] = -1
	}

	// convert to matrix type
	output := mat64.NewDense(1, 2, o)

	// return final output
	return output
}

// newrandom repeatedly generates a new random sample from mvrnd and then fixes
// it using fixrandom until the sample is comprised of a non [0, 0] case
func NewRandom(mu *mat64.Dense, sigma *mat64.SymDense) (newRand []int) {

	// initialize rndsmp and fixsmp and output variables
	rndsmp := mat64.NewDense(2, 1, nil)
	fixsmp := mat64.NewDense(1, 2, nil)

	// generate random vectors prohibiting zero-zero cases
	for {
		rndsmp = MultiVariateNormalRandom(mu, sigma)
		fixsmp = FixRandom(rndsmp)
		if fixsmp.At(0, 0) == 0 && fixsmp.At(0, 1) == 0 {
			continue
		} else {
			break
		}
	}

	// initialize output
	output := make([]int, 2)

	// write output values
	output[0] = int(fixsmp.At(0, 0))
	output[1] = int(fixsmp.At(0, 1))

	// return final output
	return output
}

// newmu generates a matrix representation of mu that reflects the
// spatial orentiation between the input current subscript and the
// destination subscript
func NewMu(curSubs, dstSubs []int) (mu *mat64.Dense) {

	// compute mu as the orientation vector
	orientVec := Orientation(curSubs, dstSubs)

	// convert mu to float
	var muVec = []float64{float64(orientVec[0]), float64(orientVec[1])}

	// initialize matrix output
	output := mat64.NewDense(2, 1, muVec)

	// return final output
	return output
}

// newsigma generates a matrix representation of sigma that reflects the
// number of iterations in the sampling process as well as the distance
// from the basis euclidean solution
func NewSigma(iterations int, randomness, distance float64) (sigma *mat64.SymDense) {

	// impose lower bound on distance
	if distance < 1 {
		distance = 1.0
	}

	// set numerator
	var num float64 = 1.0

	// initialize covariance
	var cov float64

	// compute covariance
	if distance == 1.0 {
		cov = 1.0
	} else {
		cov = math.Pow(distance, (num/randomness)) / math.Pow(float64(iterations), (num/randomness))
	}

	// initialize matrix output
	output := mat64.NewSymDense(2, nil)

	// set values
	output.SetSym(0, 0, cov)
	output.SetSym(0, 1, 0.0)
	output.SetSym(1, 0, 0.0)
	output.SetSym(1, 1, cov)

	// return final output
	return output
}

// newsubs generates a feasible new subscript value set within the
// input search domain
func NewSubs(curSubs []int, curDist float64, searchParameters *Parameters, searchDomain *Domain) (newSubscripts []int) {

	// initialize iteration counter
	var iterations int = 1

	// initialize output
	output := make([]int, 2)

	// generate and fix a bivariate normally distributed random vector
	// prohibit all zero cases and validate using the search domain
	for {

		// generate mu and sigma values
		mu := NewMu(curSubs, searchParameters.DstSubs)
		sigma := NewSigma(iterations, searchParameters.RndCoef, curDist)

		// generate fixed random bivariate normally distributed numbers
		try := NewRandom(mu, sigma)

		// write output
		output[0] = curSubs[0] + try[0]
		output[1] = curSubs[1] + try[1]

		// DEBUG
		/*		if searchDomain.Matrix.At(output[0], output[1]) == 0.0 {
				iterations += 1
				continue
			} else*/

		// test if currentIndex inside search domain
		if output[0] > searchDomain.Rows-1 || output[1] > searchDomain.Cols-1 || output[0] < 0 || output[1] < 0 {
			iterations += 1
			continue
		} else {
			break
		}
	}

	// return final output
	return output
}

// directedwalk generates a new directed walk connecting a source subscript to a
// destination subscript within the context of an input search domain
func DirectedWalk(searchDomain *Domain, searchParameters *Parameters, basisSolution *Basis) (subscripts [][]int, destinationTest bool) {

	// initialize chromosomal 2D slice with source subscript as first
	// element
	output := make([][]int, 1, searchDomain.MaxLen)
	output[0] = make([]int, 2)
	output[0][0] = searchParameters.SrcSubs[0]
	output[0][1] = searchParameters.SrcSubs[1]

	// initialize new tabu matrix
	tabu := mat64.NewDense(searchDomain.Rows, searchDomain.Cols, nil)
	for i := 0; i < searchDomain.Rows; i++ {
		for j := 0; j < searchDomain.Cols; j++ {

			if i == 0 || i == searchDomain.Rows-1 || j == 0 || j == searchDomain.Cols-1 {
				tabu.Set(i, j, 0.0)
			} else {
				tabu.Set(i, j, 1.0)
			}
		}
	}

	//tabu.Clone(searchDomain.Matrix)
	tabu.Set(searchParameters.SrcSubs[0], searchParameters.SrcSubs[1], 0.0)

	// initialize current subscripts, distance, try, and iteration counter
	curSubs := make([]int, 2)
	var curDist float64
	var try []int
	var dstTest bool

	// enter bounded for loop
	for i := 0; i < searchDomain.MaxLen; i++ {

		// get current subscripts
		curSubs = output[len(output)-1]

		// validate tabu neighborhood
		if ValidateTabu(curSubs, tabu) == false {
			break
		}

		// compute current distance
		curDist = basisSolution.Matrix.At(curSubs[0], curSubs[1])

		// generate new try
		try = NewSubs(curSubs, curDist, searchParameters, searchDomain)

		// apply control conditions
		if try[0] == searchParameters.DstSubs[0] && try[1] == searchParameters.DstSubs[1] {
			output = append(output, try)
			break
		} else if tabu.At(try[0], try[1]) == 0.0 {
			continue
		} else {
			output = append(output, try)
			tabu.Set(try[0], try[1], 0.0)
		}
	}

	if output[len(output)-1][0] == searchParameters.DstSubs[0] && output[len(output)-1][1] == searchParameters.DstSubs[1] {
		dstTest = true
	} else {
		dstTest = false
	}

	// return final output
	return output, dstTest
}

// mutationwalk generates a new directed walk connecting a source subscript to a
// destination subscript within the context of an input mutation search domain
func MutationWalk(searchDomain *Domain, searchParameters *Parameters, basisSolution *Basis) (subscripts [][]int, tabuTest bool) {

	// initialize chromosomal 2D slice with source subscript as first
	// element
	output := make([][]int, 1, searchDomain.MaxLen)
	output[0] = make([]int, 2)
	output[0][0] = searchParameters.SrcSubs[0]
	output[0][1] = searchParameters.SrcSubs[1]

	// initialize new tabu matrix
	tabu := mat64.NewDense(searchDomain.Rows, searchDomain.Cols, nil)
	tabu.Clone(searchDomain.Matrix)
	tabu.Set(searchParameters.SrcSubs[0], searchParameters.SrcSubs[1], 0.0)

	// initialize current subscripts, distance, try, and iteration counter
	curSubs := make([]int, 2)
	var curDist float64
	var try []int
	var test bool

	// enter bounded for loop
	for {

		// get current subscripts
		curSubs = output[len(output)-1]

		// compute current distance
		curDist = basisSolution.Matrix.At(curSubs[0], curSubs[1])

		// generate new try
		try = NewSubs(curSubs, curDist, searchParameters, searchDomain)

		// apply control conditions
		if try[0] == searchParameters.DstSubs[0] && try[1] == searchParameters.DstSubs[1] {
			output = append(output, try)
			break
		} else if tabu.At(try[0], try[1]) == 0.0 {
			continue
		} else {
			output = append(output, try)
			tabu.Set(try[0], try[1], 0.0)
		}

		// validate tabu matrix
		test = ValidateSubDomain(try, searchParameters.DstSubs, tabu)

		// reset if tabu is invalid
		if test == false {
			break
		}
	}

	// return final output
	return output, test
}

// newnodesubs generates an poutput slice of new intermediate destination nodes
// that are progressively further, in terms of euclidean distance, from
// a given input source location and are orientation towards a given
// destination location
func NewNodeSubs(searchDomain *Domain, searchParameters *Parameters) (nodeSubs [][]int) {

	// check band count against input distance matrix size
	if searchDomain.BndCnt < 3 {
		err := errors.New("Band count must be greater than three \n")
		panic(err)
	}

	// generate distance matrix from source subscripts
	distMat := AllDistance(searchParameters.SrcSubs, searchDomain.Matrix)

	// initialize output
	output := make([][]int, 1)
	output[0] = searchParameters.SrcSubs

	// encode distance bands
	bandMat := DistanceBands(searchDomain.BndCnt, distMat)

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// loop through band vector and generate band value subscripts
	for i := 1; i < searchDomain.BndCnt-1; i++ {

		// generate band mask
		bandMaskMat := BandMask(float64(i), bandMat)

		// generate orientation mask
		orientMaskMat := OrientationMask(output[i-1], searchParameters.DstSubs, searchDomain.Matrix)

		// initialize final mask
		finalMaskMat := mat64.NewDense(searchDomain.Rows, searchDomain.Cols, nil)

		// compute final mask through elementwise multiplication
		finalMaskMat.MulElem(bandMaskMat, orientMaskMat)

		// generate subs from final mask
		finalSubs := NonZeroSubs(finalMaskMat)

		// generate random number of length interval
		randInd := finalSubs[rand.Intn(len(finalSubs))]

		// extract randomly selected value and write to output
		output = append(output, randInd)
	}

	// set the final subscript to the destination
	output = append(output, searchParameters.DstSubs)

	// return output
	return output
}

// MultiDirectedWalk generates a new multipart directed walk from a given set
// of input problem parameters
//func MultiDirectedWalk(searchDomain *Domain, searchParameters *Parameters) (subscripts [][]int) {

//}
