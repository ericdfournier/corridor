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

// test DistanceBands
func TestDistanceBands(t *testing.T) {

	// initialize test case
	t.Log("DistanceBands Test: Expected Matrix = {{3 3 3 [0 1 2 1 1 2 2 2 2]} 3 3}")

	// initialize expected value
	var expValueVector = []float64{
		0.0, 1.0, 2.0,
		1.0, 1.0, 2.0,
		2.0, 2.0, 2.0}
	expValueMatrix := mat64.NewDense(3, 3, expValueVector)

	// initialize test case variables
	var aSubs = []int{0, 0}
	var bandCount int = 2
	searchDomainMatrix := mat64.NewDense(3, 3, nil)

	// compute distance matrix !! dependent on AllDistance test result !!
	distanceMatrix := AllDistance(aSubs, searchDomainMatrix)

	// perform test case
	testCase := DistanceBands(bandCount, distanceMatrix)

	// log test result
	if testCase.Equals(expValueMatrix) {
		t.Log("DistanceBands Test: Computed Matrix =", *testCase)
	} else {
		t.Error("DistanceBands Test: Computed Matrix =", *testCase)
	}
}

// test BandMask
func TestBandMask(t *testing.T) {

	// initialize test case
	t.Log("BandMask Test: Expected Matrix = {{3 3 3 [0 0 0 0 1 0 0 0 0]} 3 3}")

	// initialize expected value
	var expValueVector = []float64{
		0.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 0.0}
	expValueMatrix := mat64.NewDense(3, 3, expValueVector)

	// initialize test case variables
	var aSubs = []int{0, 0}
	var bandCount int = 2
	var bandValue float64 = 1.0
	searchDomainMatrix := mat64.NewDense(3, 3, nil)

	// compute distance matrix !! dependent on AllDistance test result !!
	distanceMatrix := AllDistance(aSubs, searchDomainMatrix)

	// compute band matrix !! dependent on DistanceBands test result!!
	bandMatrix := DistanceBands(bandCount, distanceMatrix)

	// perform test case
	testCase := BandMask(bandValue, bandMatrix)

	// log test results
	if testCase.Equals(expValueMatrix) {
		t.Log("BandMask Test: Computed Matrix =", *testCase)
	} else {
		t.Error("BandMask Test: Computed Matrix =", *testCase)
	}
}

// test NonZeroSubs
func TestNonZeroSubs(t *testing.T) {

	// initialize test case
	t.Log("NonZeroSubs Test: Expected Vector = [[1 1]]")

	// initialize expected value
	expValueVector := make([][]int, 1)
	expValueVector[0] = []int{1, 1}

	// initialize test case variables
	var aSubs = []int{0, 0}
	var bandCount int = 2
	var bandValue float64 = 1.0
	searchDomainMatrix := mat64.NewDense(3, 3, nil)

	// compute distance matrix !! dependent on AllDistance test result !!
	distanceMatrix := AllDistance(aSubs, searchDomainMatrix)

	// compute band matrix !! dependent on DistanceBands test result!!
	bandMatrix := DistanceBands(bandCount, distanceMatrix)

	// compute band mask !! dependent on BandMask test result!!
	bandMask := BandMask(bandValue, bandMatrix)

	// perform test case
	testCase := NonZeroSubs(bandMask)

	// log test results
	if testCase[0][0] == expValueVector[0][0] && testCase[0][1] == expValueVector[0][1] {
		t.Log("NonZeroSubs Test: Computed Vector =", testCase)
	} else {
		t.Error("NonZeroSubs Test: Computed Vector =", testCase)
	}
}

// test FindSubs
func TestFindSubs(t *testing.T) {

	// initialize test case
	t.Log("FindSubs Test: Expected Vector = [[0 0]]")

	// initialize expected value
	expValueVector := make([][]int, 1)
	expValueVector[0] = []int{0, 0}

	// initialize test case variables
	var inputValue float64 = 0.0
	var aSubs = []int{0, 0}
	searchDomainMatrix := mat64.NewDense(3, 3, nil)

	// compute distance matrix !! dependent on AllDistance test result !!
	distanceMatrix := AllDistance(aSubs, searchDomainMatrix)

	// perform test case
	testCase := FindSubs(inputValue, distanceMatrix)

	// log test results
	if testCase[0][0] == expValueVector[0][0] && testCase[0][1] == expValueVector[0][1] {
		t.Log("FindSubs Test: Computed Vector =", testCase)
	} else {
		t.Error("FindSubs Test: Computed Vector =", testCase)
	}
}

// test Orientation
func TestOrientation(t *testing.T) {

	// initialize test case
	t.Log("Orientation Test: Expected Vector = [1 1]")

	// initialize expected value
	var expValueVector = []int{1, 1}

	// initialize test case variables
	var aSubs = []int{0, 0}
	var bSubs = []int{2, 2}

	// perform test case
	testCase := Orientation(aSubs, bSubs)

	// log test results
	if testCase[0] == expValueVector[0] && testCase[1] == expValueVector[1] {
		t.Log("Orientation Test: Computed Vector =", testCase)
	} else {
		t.Error("Orientation Test: Computed Vector =", testCase)
	}
}

// test OrientationMask
func TestOrientationMask(t *testing.T) {

	// initialize test case
	t.Log("OrientationMask Test: Expected Matrix = {{3 3 3 [0 0 0 0 1 0 0 0 0]} 3 3}")

	// initialize expected value
	var expValueVector = []float64{
		0.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 0.0}
	expValueMatrix := mat64.NewDense(3, 3, expValueVector)

	// initialize test case variables
	var aSubs = []int{0, 0}
	var bSubs = []int{2, 2}
	searchDomainMatrix := mat64.NewDense(3, 3, nil)

	// perform test case
	testCase := OrientationMask(aSubs, bSubs, searchDomainMatrix)

	// log test results
	if testCase.Equals(expValueMatrix) {
		t.Log("OrientationMask Test: Computed Matrix =", *testCase)
	} else {
		t.Error("OrientationMask Test: Computed Matrix =", *testCase)
	}
}

// test Bresenham
func TestBresenham(t *testing.T) {

	// initialize test case
	t.Log("Bresenham Test: Expected Vector = [[0 0] [1 1] [2 2]]")

	// initialize expected value
	expValueVector := make([][]int, 1)
	expValueVector[0] = []int{0, 0}
	expValueVector = append(expValueVector, []int{1, 1})
	expValueVector = append(expValueVector, []int{2, 2})

	// initialize test case variables
	var aSubs = []int{0, 0}
	var bSubs = []int{2, 2}
	var testBool bool = true

	// perform test case
	testCase := Bresenham(aSubs, bSubs)

	// examine test results
	for i := 0; i < len(testCase); i++ {
		if testCase[i][0] != expValueVector[i][0] || testCase[i][1] != expValueVector[i][1] {
			testBool = false
			break
		}
	}

	// log test results
	if testBool == true {
		t.Log("Bresenham Test: Computed Vector =", testCase)
	} else {
		t.Error("Bresenham Test: Computed Vector =", testCase)
	}
}

// test NeighborhoodSubs
func TestNeighborhoodSubs(t *testing.T) {

	// initialize test case
	t.Log("NeighborhoodSubs Test: Expected Vector = [[0 0] [0 1] [0 2] [1 0] [1 1] [1 2] [2 0] [2 1] [2 2]]")

	// initialize expected value
	expValueVector := make([][]int, 0)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			expValueVector = append(expValueVector, []int{i, j})
		}
	}

	// initialize test ccase variables
	var aSubs = []int{1, 1}
	var testBool bool = true

	// perform test case
	testCase := NeighborhoodSubs(aSubs)

	// examine test results
	for k := 0; k < len(testCase); k++ {
		if testCase[k][0] != expValueVector[k][0] || testCase[k][1] != expValueVector[k][1] {
			testBool = false
			break
		}
	}

	// log test results
	if testBool == true {
		t.Log("NeighborhoodSubs Test: Computed Vector =", testCase)
	} else {
		t.Error("NeighborhoodSubs Test: Computed Vector =", testCase)
	}
}

// test ValidateSubDomain
func TestValidateSubDomain(t *testing.T) {

	// initialize test case
	t.Log("ValidateSubDomain: Expecte Value = true")

	// initialize expected values
	var invalidVector = []float64{
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	invalidMatrix := mat64.NewDense(5, 5, invalidVector)
	var validVector = []float64{
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	validMatrix := mat64.NewDense(5, 5, validVector)

	// initialize test case variables
	var subSource = []int{1, 1}
	var subDestin = []int{2, 3}
	var testCase1 bool
	var testCase2 bool

	// perform test cases
	testCase1 = ValidateMutationSubDomain(subSource, subDestin, invalidMatrix)
	testCase2 = ValidateMutationSubDomain(subSource, subDestin, validMatrix)

	// log test results
	if testCase1 == false && testCase2 == true {
		t.Log("ValidateMutationSubDomain Test: Computed Value =", true)
	} else {
		t.Error("ValidateMutationSubDomain Test: Computed Value =", false)
	}
}

// test ValidateTabu
func TestValidateTabu(t *testing.T) {

	// initialize test case
	t.Log("ValidateTabu Test: Expected Value = true")

	// initialize expected values
	var invalidVector = []float64{
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	invalidMatrix := mat64.NewDense(5, 5, invalidVector)
	var validVector = []float64{
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 1.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	validMatrix := mat64.NewDense(5, 5, validVector)

	// initialize test case variables
	var currentSubs = []int{2, 2}
	var testCase1 bool
	var testCase2 bool

	// perform test cases
	testCase1 = ValidateTabu(currentSubs, invalidMatrix)
	testCase2 = ValidateTabu(currentSubs, validMatrix)
	testBool := testCase1 == false && testCase2 == true

	// log test results
	if testBool {
		t.Log("ValidateTabu Test: Computed Value =", testBool)
	} else {
		t.Error("ValidateTabu Test: Computed Value =", testBool)
	}
}

// test DigitCount
func TestDigitCount(t *testing.T) {

	// initialize test case
	t.Log("DigitCount Test: Expected Value = 10")

	// initialize expected values
	var expValue int = 10

	// initialize test case variables
	input := 10000000000

	// peform test case
	testCase := DigitCount(input)

	// log test results
	if testCase == expValue {
		t.Log("DigitCount Test: Computed Value =", testCase)
	} else {
		t.Error("DigitCount Test: Computed Value =", testCase)
	}
}
