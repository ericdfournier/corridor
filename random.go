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
func Newind(currentIndex int, mu, sigma *mat64.Dense, searchDomain *Domain) (newIndex int) {

	// initialize output
	var output int

	// generate and fix a bivariate normally distributed random vector
	// prohibit all zero cases and validate using the search domain
	for {

		// generate fixed random bivariate normally distributed numbers
		nR := Newrnd(mu, sigma)

		// check cases and assign values
		if nR[0] == 0 && nR[1] == -1 {
			output = currentIndex - 1
		} else if nR[0] == 0 && nR[1] == 1 {
			output = currentIndex + 1
		} else if nR[0] == -1 && nR[1] == -1 {
			output = currentIndex - searchDomain.Stride - 1
		} else if nR[0] == -1 && nR[1] == 1 {
			output = currentIndex - searchDomain.Stride + 1
		} else if nR[0] == 1 && nR[1] == 1 {
			output = currentIndex + searchDomain.Stride + 1
		} else if nR[0] == 1 && nR[1] == -1 {
			output = currentIndex + searchDomain.Stride - 1
		} else if nR[0] == -1 && nR[1] == 0 {
			output = currentIndex - searchDomain.Stride
		} else if nR[0] == 1 && nR[1] == 0 {
			output = currentIndex + searchDomain.Stride
		}

		// NEED TO DEVELOP A TEST FOR INDEX OUT OF RANGE ERROR

		// test if currentIndex inside search domain
		if searchDomain.Vals[output] == false {
			continue
		} else {
			break
		}
	}

	// return final output
	return output

}
