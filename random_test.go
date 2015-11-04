/* Copyright ©2015 The corridor Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file. */

package corridor

import (
	"math"
	"testing"

	"github.com/gonum/matrix/mat64"
)

// test multivariaterandomnormal
func TestMultiVariateRandomNormal(t *testing.T) {

	// initialize test case
	t.Log("MultiVariateRandomNormal Test: Expected Value = [(1 +- 0.1) (1 +- 0.1)]")

	// initialize expected values
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

// test fixmultivariatenormalrandom
func TestFixMultiVariateNormalRandom(t *testing.T) {

	// initialize test case
	t.Log("FixMultiVariateNormalRandom Test: Expected Value = [(0.47 +- 0.1) (0.47 +- 0.1)]")

	// initialize expected values
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

	// generate fixed random samples
	for i := 0; i < 10000; i++ {
		curRnd := MultiVariateNormalRandom(mu, sigma)
		curFix := FixMultiVariateNormalRandom(curRnd)
		testMat.Set(i, 0, curFix.At(0, 0))
		testMat.Set(i, 1, curFix.At(0, 1))
	}

	// compute test result
	for i := 0; i < 10000; i++ {
		aSum += testMat.At(i, 0)
		bSum += testMat.At(i, 1)
	}
	testCase[0] = aSum / 10000.0
	testCase[1] = bSum / 10000.0
	testBool = (math.Abs(0.47-testCase[0]) < 0.1) && (math.Abs(0.47-testCase[1]) < 0.1)

	// log test results
	if testBool {
		t.Log("FixMultiVariateNormalRandom Test: Computed Mean =", testCase)
	} else {
		t.Error("FixMultiVariateNormalRandom Test: Computed Value =", testCase)
	}
}

// test newrandom
func TestNewRandom(t *testing.T) {

	// initialize test case
	t.Log("NewRandom Test: Expected Value = [(0.615 +- 0.1) (0.615 +- 0.1)]")

	// initialize expected values
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

	// generate random samples
	for i := 0; i < 10000; i++ {
		curVal := NewRandom(mu, sigma)
		testMat.Set(i, 0, float64(curVal[0]))
		testMat.Set(i, 1, float64(curVal[1]))
	}

	// compute test result
	for i := 0; i < 10000; i++ {
		aSum += testMat.At(i, 0)
		bSum += testMat.At(i, 1)
	}
	testCase[0] = aSum / 10000.0
	testCase[1] = bSum / 10000.0
	testBool = (math.Abs(0.615-testCase[0]) < 0.1) && (math.Abs(0.615-testCase[1]) < 0.1)

	// log test results
	if testBool {
		t.Log("NewRandom Test: Computed Value =", testCase)
	} else {
		t.Error("NewRandom Test: Computed Value =", testCase)
	}
}

// test newmu
func TestNewMu(t *testing.T) {

	// initialize test case
	t.Log("NewMu Test: Expected Value = {{2 1 1 [1 1]} 2 1}")

	// initialize expected values
	var expVal = []float64{1.0, 1.0}

	// initialize test case variables
	var curSubs = []int{10, 10}
	var dstSubs = []int{100, 100}

	// perform test case
	testCase := NewMu(curSubs, dstSubs)

	// log test results
	if testCase.At(0, 0) == expVal[0] && testCase.At(1, 0) == expVal[1] {
		t.Log("NewMu Test: Computed Value =", *testCase)
	} else {
		t.Error("NewMu Test: Computed Value =", *testCase)
	}
}

// test newsigma
func TestNewSigma(t *testing.T) {

	// initialize test case
	t.Log("NewSigma Test: Expected Matrix = {{2 1 [10 0 0 10]} 121}")

	// initialize expected values
	var expValVec = []float64{10.0, 0.0, 0.0, 10.0}
	expValMat := mat64.NewDense(2, 2, expValVec)

	// initialize test case variables
	var iterations int = 1
	var randomness float64 = 1.0
	var distance float64 = 10.0

	// perform test case
	testCase := NewSigma(iterations, randomness, distance)

	// log test results
	if testCase.At(0, 0) == expValMat.At(0, 0) &&
		testCase.At(0, 1) == expValMat.At(0, 1) &&
		testCase.At(1, 0) == expValMat.At(1, 0) &&
		testCase.At(1, 1) == expValMat.At(1, 1) {
		t.Log("NewSigma Test: Computed Matrix =", *testCase)
	} else {
		t.Error("NewSigma Test: Compute Matrix =", *testCase)
	}
}

// test newsubs
func TestNewSubs(t *testing.T) {

	// initialize test case
	t.Log("NewSubs Test: Expected Value = [1 2] or [2 1]")

	// initialize expected values
	var expVal1 = []int{1, 2}
	var expVal2 = []int{2, 1}

	// initialize test case variables
	var curSubs = []int{1, 1}
	var dstSubs = []int{3, 3}
	var curDist float64 = 0.0
	testParams := NewParameters(curSubs, dstSubs, 10, 10, 1.0)
	var domainVec = []float64{
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	domainMat := mat64.NewDense(5, 5, domainVec)
	testDomain := NewDomain(domainMat)

	// perform test case
	testCase := NewSubs(curSubs, dstSubs, curDist, testParams, testDomain)

	// log test results
	if (testCase[0] == expVal1[0] && testCase[1] == expVal1[1]) ||
		(testCase[0] == expVal2[0] && testCase[1] == expVal2[1]) {
		t.Log("NewSubs Test: Computed Value =", testCase)
	} else {
		t.Error("NewSubs Test: Computed Value =", testCase)
	}
}

// test directedwalk
func TestDirectedWalk(t *testing.T) {

	// initialize test case
	t.Log("DirectedWalk Test: Expected Value = [[1 1] [1 2] [2 3] [3 3]]")

	// initialize expected values
	expVal := make([][]int, 5)
	expVal[0] = []int{1, 1}
	expVal[1] = []int{1, 2}
	expVal[2] = []int{2, 3}
	expVal[3] = []int{3, 3}

	// initialize test case variables
	var sourceSubs = []int{1, 1}
	var destinationSubs = []int{3, 3}
	testParams := NewParameters(sourceSubs, destinationSubs, 10, 10, 1.0)
	var domainVec = []float64{
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	domainMat := mat64.NewDense(5, 5, domainVec)
	testDomain := NewDomain(domainMat)
	testBasis := NewBasis(sourceSubs, destinationSubs, testDomain)
	var testBool bool

	// perform test case
	testCase := DirectedWalk(sourceSubs, destinationSubs, testDomain, testParams, testBasis)

	// evaluate test results
	for i := 0; i < 4; i++ {
		if testCase[i][0] == expVal[i][0] && testCase[i][1] == expVal[i][1] {
			testBool = true
		} else {
			testBool = false
			break
		}
	}

	// log test results
	if testBool {
		t.Log("DirectedWalk Test: Computed Value =", testCase)
	} else {
		t.Error("DirectedWalk Test: Computed Value =", testCase)
	}
}

// test mutationwalk
func TestMutationWalk(t *testing.T) {

	// initialize test case
	t.Log("MutationWalk Test: Expected Value = [[1 1] [2 1] [3 2] [3 3]] or [[1 1] [1 2] [2 3] [3 3]]")

	// initialize test case variables
	var sourceSubs = []int{1, 1}
	var destinationSubs = []int{3, 3}
	testParams := NewParameters(sourceSubs, destinationSubs, 10, 10, 1.0)
	var domainVec = []float64{
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	domainMat := mat64.NewDense(5, 5, domainVec)
	testDomain := NewDomain(domainMat)
	testBasis := NewBasis(sourceSubs, destinationSubs, testDomain)
	var testBool bool
	var testCase [][]int

	// perform test case
	for {
		testCase, testBool = MutationWalk(sourceSubs, destinationSubs, testDomain, testParams, testBasis)
		if testBool == true {
			break
		} else {
			continue
		}
	}

	// log test results
	if testBool {
		t.Log("MutationWalk Test: Computed Value =", testCase)
	} else {
		t.Error("MutationWalk Test: Computed Value =", testCase)
	}
}

// test newnodesubs
func TestNewNodeSubs(t *testing.T) {

	// initialize test case
	t.Log("NewNodeSubs Test: Expected Value = [[1 1] [2 2] [3 3]]")

	// initialize expected values
	expVal := make([][]int, 3)
	expVal[0] = []int{1, 1}
	expVal[1] = []int{2, 2}
	expVal[2] = []int{3, 3}

	// initialize test case variables
	var sourceSubs = []int{1, 1}
	var destinationSubs = []int{3, 3}
	testParams := NewParameters(sourceSubs, destinationSubs, 10, 10, 1.0)
	var domainVec = []float64{
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	domainMat := mat64.NewDense(5, 5, domainVec)
	testDomain := NewDomain(domainMat)
	testDomain.BndCnt = 3

	// perform test case
	testCase := NewNodeSubs(testDomain, testParams)
	testBool := (testCase[0][0] == expVal[0][0] &&
		testCase[0][1] == expVal[0][1] &&
		testCase[1][0] == expVal[1][0] &&
		testCase[1][1] == expVal[1][1] &&
		testCase[2][0] == expVal[2][0] &&
		testCase[2][1] == expVal[2][1])

	// log test results
	if testBool {
		t.Log("NewNodeSubs Test: Computed Value =", testCase)
	} else {
		t.Error("NewNodeSubs Test: Computed Value =", testCase)
	}
}

// test multipartdirectedwalk
func TestMultiPartDirectedWalk(t *testing.T) {

	// initialize test case
	t.Log("MultiPartDirectedWalk: Expected Value = [[1 1]...[3 3]]")

	// initialize test case variables
	var sourceSubs = []int{1, 1}
	var destinationSubs = []int{3, 3}
	testParams := NewParameters(sourceSubs, destinationSubs, 10, 10, 1.0)
	var domainVec = []float64{
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	domainMat := mat64.NewDense(5, 5, domainVec)
	testDomain := NewDomain(domainMat)
	testDomain.BndCnt = 3
	nodeSubs := make([][]int, 3)
	nodeSubs[0] = sourceSubs
	nodeSubs[1] = []int{1, 3}
	nodeSubs[2] = destinationSubs

	// perform test case
	testCase := MultiPartDirectedWalk(nodeSubs, testDomain, testParams)

	// evaluate test case
	t.Log("MultiPartDirectedWalk: Computed Value =", testCase)
}
