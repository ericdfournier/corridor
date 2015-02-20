// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"math/rand"
	"time"

	"github.com/gonum/matrix/mat64"
)

func Mvnrnd(mu, sigma *mat64.Dense) (out *mat64.Dense) {

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

func Fixrnd(rnd *mat64.Dense) (out *mat64.Dense) {
	
	// initialize vector slice
	o := make([]float64, 2)

	// write up down movement direction
	if rnd.At(0, 0) > 1 {
		o[0] = 1
	} else if rnd.At(0, 0) >= -1 && rnd.At(0, 0) <= 1 {
		o[0] = 0
	} else if rnd.At(0, 0) < -1 {
		o[0] = -1
	}

	// write left right movement direction
	if rnd.At(0, 1) > 1 {
		o[1] = 1
	} else if rnd.At(0, 1) >= -1 && rnd.At(0, 0) <= 1 {
		o[1] = 0
	} else if rnd.At(0, 1) < -1 {
		o[1] = -1
	}

	// convert to matrix type
	output := mat64.NewDense(1, 2, o)

	// return final output
	return output
}
