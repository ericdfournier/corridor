// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
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
	rnd := mat64.NewDense(1, 2, n)
	output := mat64.NewDense(1, 2, o)

	// perform cholesky decomposition on covariance matrix
	cholFactor := mat64.Cholesky(sigma)
	lower := cholFactor.L

	// compute output
	output.Mul(rnd, lower)
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
	if rndsmp.At(0, 0) > 1 {
		o[0] = 1
	} else if rndsmp.At(0, 0) >= -1 && rndsmp.At(0, 0) <= 1 {
		o[0] = 0
	} else if rndsmp.At(0, 0) < -1 {
		o[0] = -1
	}

	// write left right movement direction
	if rndsmp.At(0, 1) > 1 {
		o[1] = 1
	} else if rndsmp.At(0, 1) >= -1 && rndsmp.At(0, 0) <= 1 {
		o[1] = 0
	} else if rndsmp.At(0, 1) < -1 {
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
	rndsmp := mat64.NewDense(1, 2, nil)
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

// newind generates a feasible new index value within the input search
// domain
func Newind(currentSubscripts []int, mu, sigma *mat64.Dense, searchDomain *Domain) (newSubscripts []int) {

	// initialize output
	output := make([]int, 2)

	// get search domain matrix dimensions
	maxRows, maxCols := searchDomain.Matrix.Dims()

	// generate and fix a bivariate normally distributed random vector
	// prohibit all zero cases and validate using the search domain
	for {

		// generate fixed random bivariate normally distributed numbers
		try := Newrnd(mu, sigma)

		for i := 0; i < 2; i++ {
			output[i] = currentSubscripts[i] + try[i]
		}

		// test if currentIndex inside search domain
		if searchDomain.Matrix.At(output[0], output[1]) == 0.0 {
			continue
		} else if output[0] > maxRows-1 || output[1] > maxCols-1 || output[0] < 0 || output[1] < 0 {
			continue
		} else {
			break
		}
	}

	// return final output
	return output
}
