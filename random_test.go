// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"math"
	"testing"

	"github.com/gonum/matrix/mat64"
)

// test multivariaterandomnormal
func TestMultiVariateRandomNormal(t *testing.T) {

	// initialize test case
	t.Log("MultiVariateRandomNormal Test: Expected Value = [1 +- 0.1 1 +- 0.1]")

	// initialize expected mean
	var muVec = []float64{1, 1}
	mu := mat64.NewDense(2, 1, muVec)
	var sigmaVec = []float64{1, 0, 0, 1}
	sigma := mat64.NewSymDense(2, sigmaVec)

	// initialize test case variables
	testMat := mat64.NewDense(10000, 2, nil)
	testCase := make([]float64, 2)
	var testBool bool
	var aSum float64 = 0.0
	var bSum float64 = 0.0

	// perform test case
	for i := 0; i < 10000; i++ {
		curVal := MultiVariateNormalRandom(mu, sigma)
		testMat.Set(i, 0, curVal.At(0, 0))
		testMat.Set(i, 1, curVal.At(1, 0))
	}

	// compute test result
	for i := 0; i < 10000; i++ {
		aSum += testMat.At(i, 0)
		bSum += testMat.At(i, 1)
	}
	testCase[0] = aSum / 10000.0
	testCase[1] = bSum / 10000.0
	testBool = (math.Abs(1-testCase[0]) < 0.1) && (math.Abs(1-testCase[1]) < 0.1)

	// log test results
	if testBool {
		t.Log("MultiVariateRandomNormal Test: Computed Mean =", testCase)
	} else {
		t.Error("MultiVariateRandomNormal Test: Computed Mean =", testCase)
	}
}
