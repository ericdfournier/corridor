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

	// generate random numbers from normal distribution
	rand.Seed(time.Now().Unix())
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
