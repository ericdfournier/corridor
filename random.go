// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"math"
	"math/rand"
	"time"

	"github.com/gonum/matrix/mat64"
)

// mvnrnd generates pairs of bivariate normally distributed random numbers
// given an input mean vector and covariance matrix
func Mvnrnd(mu, sigma *mat64.Dense) (rndsmp *mat64.Dense) {

	// initialize vector slices
	o := make([]float64, 2)
	n := make([]float64, 2)

	// generate random numbers from normal distribution, prohibit [0,0]
	// combinations
	rand.Seed(time.Now().UnixNano())
	for i := 0; i <= 1; i++ {
		n[i] = rand.NormFloat64()
	}

	// convert to matrix type
	rnd := mat64.NewDense(2, 1, n)
	output := mat64.NewDense(2, 1, o)

	// perform cholesky decomposition on covariance matrix
	cholFactor := mat64.Cholesky(sigma)
	lower := cholFactor.L

	// compute output
	output.Mul(lower, rnd)
	output.Add(output, mu)

	//return final output
	return output
}

// fixrnd converts an input vector of bivariate normally distributed random
// numbers into a version where the values have been fixed to a [-1, 0 ,1]
// range
func Fixrnd(rndsmp *mat64.Dense) (fixsmp *mat64.Dense) {

	// initialize vector slice
	o := make([]float64, 2)

	// write up down movement direction
	if rndsmp.At(0, 0) > 0.5 {
		o[0] = 1
	} else if rndsmp.At(0, 0) >= -0.5 && rndsmp.At(0, 0) <= 0.5 {
		o[0] = 0
	} else if rndsmp.At(0, 0) < -0.5 {
		o[0] = -1
	}

	// write left right movement direction
	if rndsmp.At(1, 0) > 0.5 {
		o[1] = 1
	} else if rndsmp.At(1, 0) >= -0.5 && rndsmp.At(1, 0) <= 0.5 {
		o[1] = 0
	} else if rndsmp.At(1, 0) < -0.5 {
		o[1] = -1
	}

	// convert to matrix type
	output := mat64.NewDense(1, 2, o)

	// return final output
	return output
}

// newrnd repeatedly generates a new random sample from mvrnd and then fixes
// it using fixrnd until the sample is comprised of a non [0, 0] case
func Newrnd(mu, sigma *mat64.Dense) (newRand []int) {

	// initialize rndsmp and fixsmp and output variables
	rndsmp := mat64.NewDense(2, 1, nil)
	fixsmp := mat64.NewDense(1, 2, nil)

	// generate random vectors prohibiting zero-zero cases
	for {
		rndsmp = Mvnrnd(mu, sigma)
		fixsmp = Fixrnd(rndsmp)
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
func Newmu(curSubs, dstSubs []int) (mu *mat64.Dense) {

	// initialize vector slice
	muVec := make([]float64, 2)

	// assign row based parameter
	if curSubs[0]-dstSubs[0] < 0 {
		muVec[0] = 1
	} else if curSubs[0]-dstSubs[0] == 0 {
		muVec[0] = 0
	} else if curSubs[0]-dstSubs[0] > 0 {
		muVec[0] = -1
	}

	// assign column based parameter
	if curSubs[1]-dstSubs[1] < 0 {
		muVec[1] = 1
	} else if curSubs[1]-dstSubs[1] == 0 {
		muVec[1] = 0
	} else if curSubs[1]-dstSubs[1] > 0 {
		muVec[1] = -1
	}

	// initialize matrix output
	output := mat64.NewDense(2, 1, muVec)

	// return final output
	return output
}

// newsig generates a matrix representation of sigma that reflects the
// number of iterations in the sampling process as well as the distance
// from the basis euclidean solution
func Newsig(iterations int, randomness, distance float64) (sigma *mat64.Dense) {

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
	output := mat64.NewDense(2, 2, nil)

	// set values
	output.Set(0, 0, cov)
	output.Set(0, 1, 0.0)
	output.Set(1, 0, 0.0)
	output.Set(1, 1, cov)

	// return final output
	return output
}

// newind generates a feasible new index value within the input search
// domain
func Newind(curSubs []int, curDist float64, searchParameters *Parameters, searchDomain *Domain) (newSubscripts []int) {

	// initialize iteration counter
	var iterations int = 1

	// initialize output
	output := make([]int, 2)

	// get search domain matrix dimensions
	maxRows, maxCols := searchDomain.Matrix.Dims()

	// generate and fix a bivariate normally distributed random vector
	// prohibit all zero cases and validate using the search domain
	for {

		// generate mu and sigma values
		mu := Newmu(curSubs, searchParameters.DstSubs)
		sigma := Newsig(iterations, searchParameters.RndCoef, curDist)

		// generate fixed random bivariate normally distributed numbers
		try := Newrnd(mu, sigma)

		// write output
		output[0] = curSubs[0] + try[0]
		output[1] = curSubs[1] + try[1]

		// test if currentIndex inside search domain
		if searchDomain.Matrix.At(output[0], output[1]) == 0.0 {
			iterations += 1
			continue
		} else if output[0] > maxRows-1 || output[1] > maxCols-1 || output[0] < 0 || output[1] < 0 {
			iterations += 1
			continue
		} else {
			break
		}
	}

	// return final output
	return output
}

// dirwlk generates a new directed walk connecting a source subscript to a
// destination subscript within the context of an input search domain
func Dirwlk(searchDomain *Domain, searchParameters *Parameters, basisSolution *Basis) (subscripts [][]int) {

	// initialize iterator and output variables
	rows, cols := searchDomain.Matrix.Dims()

	// set maximum chromosome length based upon search domain dimensions
	//var p float64 = 2
	//var s float64 = 5
	//maxLen := int(math.Ceil(s * math.Sqrt(math.Pow(float64(rows), p)+math.Pow(float64(cols), p))))

	// initialize chromosomal 2D slice with source subscript as first
	// element
	output := make([][]int, 1, searchDomain.MaxLen)
	output[0] = make([]int, 2)
	output[0][0] = searchParameters.SrcSubs[0]
	output[0][1] = searchParameters.SrcSubs[1]

	// initialize new tabu matrix
	tabu := mat64.NewDense(rows, cols, nil)
	tabu.Clone(searchDomain.Matrix)
	tabu.Set(searchParameters.SrcSubs[0], searchParameters.SrcSubs[1], 0.0)

	// initialize current subscripts, distance, try, and iteration counter
	curSubs := make([]int, 2)
	var curDist float64
	var try []int

	// enter unbounded for loop
	for i := 0; i < searchDomain.MaxLen; i++ {

		curSubs = output[len(output)-1]
		curDist = basisSolution.Matrix.At(curSubs[0], curSubs[1])
		try = Newind(curSubs, curDist, searchParameters, searchDomain)

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

	// return final output
	return output
}
