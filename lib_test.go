// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"math"
	"testing"

	"github.com/gonum/matrix/mat64"
)

// test Distance function
func TestDistance(t *testing.T) {

	// initialize test case
	t.Log("Distance Test: Expected Distance = 10")

	// initialize expected values
	var expValue float64 = 10.0

	// initialize test case variables
	var aSubs = []int{0, 10}
	var bSubs = []int{0, 20}

	// perform test case
	testCase := Distance(aSubs, bSubs)

	// log test result
	if testCase == expValue {
		t.Log("Distance Test: Computed Distance =", testCase)
	} else {
		t.Error("Distance Test: Computed Distance =", testCase)
	}
}

// test AllDistance function
func TestAllDistance(t *testing.T) {

	// initialize test case
	t.Log("AllDistance Test: Expected Matrix = {{3 3 3 [1.4142135623730951 1 1.4142135623730951 1 0 1 1.4142135623730951 1 1.4142135623730951]} 3 3}")

	// initialize expected value
	var expValueVector = []float64{
		math.Sqrt(2.0), 1.0, math.Sqrt(2.0),
		1.0, 0, 1.0,
		math.Sqrt(2.0), 1.0, math.Sqrt(2.0)}
	expValueMatrix := mat64.NewDense(3, 3, expValueVector)

	// initialize test case variables
	var aSubs = []int{1, 1}
	searchDomainMatrix := mat64.NewDense(3, 3, nil)

	// perform test case
	testCase := AllDistance(aSubs, searchDomainMatrix)

	// log test result
	if testCase.Equals(expValueMatrix) == true {
		t.Log("AllDistance Test: Computed Matrix =", *testCase)
	} else {
		t.Error("AllDistance Test: Computed Matrix =", *testCase)
	}
}

// test MinDistance function
func TestMinDistance(t *testing.T) {

	// initialize test case
	t.Log("MinDistance Test: Expected Value = 1.4142135623730951")

	// initialize expected value
	var expValue float64 = math.Sqrt(2.0)

	// initialize test case variables
	var pSubs = []int{0, 2}
	var aSubs = []int{0, 0}
	var bSubs = []int{2, 2}

	// perform test case
	testCase := MinDistance(pSubs, aSubs, bSubs)

	// log test result
	if testCase == expValue {
		t.Log("MinDistance Test: Computed Value =", testCase)
	} else {
		t.Error("AllDistance Test: Computed Value =", testCase)
	}
}

// test AllMinDistance function
func TestAllMinDistance(t *testing.T) {

	// initialize test case
	t.Log("AllMinDistance Test: Expected Matrix = {{3 3 3 [0 1 1.4142135623730951 1 0 1 1.4142135623730951 1 0]} 3 3}")

	// initialize expected value
	var expValueVector = []float64{
		0.0, (math.Sqrt(2.0) / 2.0), math.Sqrt(2.0),
		(math.Sqrt(2.0) / 2.0), 0.0, (math.Sqrt(2.0) / 2.0),
		math.Sqrt(2.0), (math.Sqrt(2.0) / 2.0), 0.0}
	expValueMatrix := mat64.NewDense(3, 3, expValueVector)

	// initialize test case variables
	var aSubs = []int{0, 0}
	var bSubs = []int{2, 2}
	searchDomainMatrix := mat64.NewDense(3, 3, nil)

	// peform test case
	testCase := AllMinDistance(aSubs, bSubs, searchDomainMatrix)

	// log test result
	if testCase.Equals(expValueMatrix) {
		t.Log("AllMinDistance Test: Computed Matrix =", testCase)
	} else {
		t.Error("AllMinDistance Test: Computed Matrix =", testCase)
	}
}
