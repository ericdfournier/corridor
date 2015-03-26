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
func MvnRnd(mu *mat64.Dense, sigma *mat64.SymDense) (rndsmp *mat64.Dense) {

	// initialize vector slices
	o := make([]float64, 2)
	n := make([]float64, 2)

	// generate random numbers from normal distribution, prohibit [0,0]
	// combinations
	rand.Seed(time.Now().UnixNano())

	// enter loop
	for i := 0; i <= 1; i++ {
		n[i] = rand.NormFloat64()
	}

	// convert to matrix type
	rnd := mat64.NewDense(2, 1, n)
	output := mat64.NewDense(2, 1, o)

	// perform cholesky decomposition on covariance matrix
	lower := mat64.NewTriangular(2, false, nil)
	lower.Cholesky(sigma, true)

	// compute output
	output.Mul(lower, rnd)
	output.Add(output, mu)

	//return final output
	return output
}

// fixrnd converts an input vector of bivariate normally distributed random
// numbers into a version where the values have been fixed to a [-1, 0 ,1]
// range
func FixRnd(rndsmp *mat64.Dense) (fixsmp *mat64.Dense) {

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

// newpsdrnd repeatedly generates a new random sample from mvrnd and then fixes
// it using fixrnd until the sample is comprised of a non [0, 0] case
func NewPsdRnd(mu *mat64.Dense, sigma *mat64.SymDense) (newRand []int) {

	// initialize rndsmp and fixsmp and output variables
	rndsmp := mat64.NewDense(2, 1, nil)
	fixsmp := mat64.NewDense(1, 2, nil)

	// generate random vectors prohibiting zero-zero cases
	for {
		rndsmp = MvnRnd(mu, sigma)
		fixsmp = FixRnd(rndsmp)
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
func NewSig(iterations int, randomness, distance float64) (sigma *mat64.SymDense) {

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

// newpsdind generates a feasible new index value within the input search
// domain
func NewPsdInd(curSubs []int, curDist float64, searchParameters *Parameters, searchDomain *Domain) (newSubscripts []int) {

	// initialize iteration counter
	var iterations int = 1

	// initialize output
	output := make([]int, 2)

	// generate and fix a bivariate normally distributed random vector
	// prohibit all zero cases and validate using the search domain
	for {

		// generate mu and sigma values
		mu := NewMu(curSubs, searchParameters.DstSubs)
		sigma := NewSig(iterations, searchParameters.RndCoef, curDist)

		// generate fixed random bivariate normally distributed numbers
		try := NewPsdRnd(mu, sigma)

		// write output
		output[0] = curSubs[0] + try[0]
		output[1] = curSubs[1] + try[1]

		// test if currentIndex inside search domain
		if searchDomain.Matrix.At(output[0], output[1]) == 0.0 {
			iterations += 1
			continue
		} else if output[0] > searchDomain.Rows-1 || output[1] > searchDomain.Cols-1 || output[0] < 0 || output[1] < 0 {
			iterations += 1
			continue
		} else {
			break
		}
	}

	// return final output
	return output
}

//func NewRndInd(curSubs []int, searchParameters *Parameters, searchDomain *Domain, tabuMatrix *mat64.Dense) (newSubscripts []int) {

//	// initialize output
//	output := make([]int, 2)

//	// enter unbounded for loop
//	for {

//		// seed random number generator
//		rand.Seed(time.Now().UnixNano())

//		// randomly generate compute sign
//		sign1 := rand.Intn(2) - 1
//		sign2 := rand.Intn(2) - 1

//		// randomnly generate values
//		value1 := rand.Intn(2) - 1
//		value2 := rand.Intn(2) - 1

//		output[0] = curSubs[0]
//		output[1] = curSubs[1]

//		// assign signs to values
//		if sign1 == 0 && sign2 == 0 {
//			output[0] = output[0] - value1
//			output[1] = output[1] - value2
//		} else if sign1 == 0 && sign2 == 1 {
//			output[0] = output[0] - value1
//			output[1] = output[1] + value2
//		} else if sign1 == 1 && sign2 == 0 {
//			output[0] = output[0] + value1
//			output[1] = output[1] - value2
//		} else {
//			output[0] = output[0] + value1
//			output[1] = output[1] + value2
//		}

//		// validate current index relative to the search domain
//		if tabuMatrix.At(output[0], output[1]) == 1.0 && output[0] < searchDomain.Rows-1 && output[1] < searchDomain.Cols-1 && output[0] > 0 && output[1] > 0 {
//			break
//		} else {
//			continue
//		}
//	}

//	// return output
//	return output
//}

// dirwlk generates a new directed walk connecting a source subscript to a
// destination subscript within the context of an input search domain
func DirWlk(searchDomain *Domain, searchParameters *Parameters, basisSolution *Basis) (subscripts [][]int) {

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

	// enter bounded for loop
	for i := 0; i < searchDomain.MaxLen; i++ {

		// get current subscripts
		curSubs = output[len(output)-1]

		// compute current distance
		curDist = basisSolution.Matrix.At(curSubs[0], curSubs[1])

		// generate new try
		try = NewPsdInd(curSubs, curDist, searchParameters, searchDomain)

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

// dirwlk generates a new directed walk connecting a source subscript to a
// destination subscript within the context of an input search domain
func MutWlk(searchDomain *Domain, searchParameters *Parameters, basisSolution *Basis) (subscripts [][]int, tabuTest bool) {

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
		try = NewPsdInd(curSubs, curDist, searchParameters, searchDomain)

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

//// rndwlk is a purely random walk procedure that connects a source subscript
//// to a destination subscript with a uniformly randomly generated
//// non-reversing walk
//func RndWlk(searchDomain *Domain, searchParameters *Parameters) (subscripts [][]int) {

//	// initialize chromosome as empty 2D slice with source subs as lead
//	output := make([][]int, 1, searchDomain.MaxLen)
//	output[0] = make([]int, 2)
//	output[0][0] = searchParameters.SrcSubs[0]
//	output[0][1] = searchParameters.SrcSubs[1]

//	// initialize new tabu matrix
//	tabu := mat64.NewDense(searchDomain.Rows, searchDomain.Cols, nil)
//	tabu.Clone(searchDomain.Matrix)
//	tabu.Set(searchParameters.SrcSubs[0], searchParameters.SrcSubs[1], 0.0)

//	// initialize current subscripts and current try
//	curSubs := make([]int, 2)
//	var try []int

//	// DEBUG
//	fmt.Println("Try")

//	// enter unbounded for loop
//	for i := 0; i < 100; i++ {

//		// get current subscripts
//		curSubs = output[len(output)-1]

//		if Distance(curSubs, searchParameters.DstSubs) < 1.5 {

//			// if one unit away assign destination
//			try = searchParameters.DstSubs
//			output = append(output, try)
//			break
//		} else {

//			// generate new try
//			try = NewRndInd(curSubs, searchParameters, searchDomain, tabu)

//			//DEBUG
//			fmt.Println("Current Try")
//			fmt.Println(try)

//			// test if destination found
//			if try[0] == searchParameters.DstSubs[0] && try[1] == searchParameters.DstSubs[1] {
//				output = append(output, try)
//				break
//			} else if tabu.At(try[0], try[1]) == 0.0 {
//				continue
//			} else {
//				output = append(output, try)
//				//tabu.Set(try[0], try[1], 0.0)
//			}

//			// DEBUG
//			fmt.Println("Tabu Matrix")
//			for i := 0; i < 5; i++ {
//				fmt.Println(tabu.RawRowView(i))
//			}

//			// check validity of current tabu matrix
//			test := ValidateSubDomain(try, searchParameters.DstSubs, tabu)

//			// DEBUG
//			fmt.Println("Tabu Matrix Validity")
//			fmt.Println(test)
//			fmt.Println("Iteration Count")

//			// if tabu matrix is invalid reset and restart
//			if test == false {

//				// DEBUG
//				fmt.Println("Reset Tabu Matrix")

//				// initialize chromosome as empty 2D slice with source subs as lead
//				output = make([][]int, 1, searchDomain.MaxLen)
//				output[0] = make([]int, 2)
//				output[0][0] = searchParameters.SrcSubs[0]
//				output[0][1] = searchParameters.SrcSubs[1]

//				// initialize new tabu matrix
//				tabu = mat64.NewDense(searchDomain.Rows, searchDomain.Cols, nil)
//				tabu.Clone(searchDomain.Matrix)
//				//tabu.Set(searchParameters.SrcSubs[0], searchParameters.SrcSubs[1], 0.0)
//			} else {
//				tabu.Set(try[0], try[1], 0.0)
//			}
//		}
//	}

//	// return output
//	return output

//}

//func BiRndWlk(searchDomain *Domain, searchParameters *Parameters) (subscripts [][]int) {

//	// initialize output
//	output := make([][]int, 1, searchDomain.MaxLen)

//	// initialize chromosome as empty 2D slice with source subs as lead
//	sourceList := make([][]int, 1, searchDomain.MaxLen)
//	sourceList[0] = make([]int, 2)
//	sourceList[0][0] = searchParameters.SrcSubs[0]
//	sourceList[0][1] = searchParameters.SrcSubs[1]

//	// initialize chromosome as empty 2D slice with source subs as lead
//	destinList := make([][]int, 1, searchDomain.MaxLen)
//	destinList[0] = make([]int, 2)
//	destinList[0][0] = searchParameters.DstSubs[0]
//	destinList[0][1] = searchParameters.DstSubs[1]

//	// initialize new source tabu matrix
//	sourceTabu := mat64.NewDense(searchDomain.Rows, searchDomain.Cols, nil)
//	sourceTabu.Clone(searchDomain.Matrix)
//	sourceTabu.Set(searchParameters.SrcSubs[0], searchParameters.SrcSubs[1], 0.0)

//	// initialize new destination tabu matrix
//	destinTabu := mat64.NewDense(searchDomain.Rows, searchDomain.Cols, nil)
//	destinTabu.Clone(searchDomain.Matrix)
//	destinTabu.Set(searchParameters.DstSubs[0], searchParameters.DstSubs[1], 0.0)

//	// initialize current subscripts and current try
//	curSrcSubs := make([]int, 2)
//	curDstSubs := make([]int, 2)
//	var sTry []int
//	var dTry []int

//	// DEBUG
//	fmt.Println("Try")

//	// enter unbounded for loop
//	for i := 0; i < 100; i++ {

//		// get current subscripts
//		curSrcSubs = sourceList[len(sourceList)-1]
//		curDstSubs = destinList[0]

//		if Distance(curSrcSubs, curDstSubs) < 1.5 {

//			// assign initial output component
//			output = sourceList[:len(sourceList)-1]

//			// if one unit away assign destination list to output
//			for i := 0; i < len(destinList); i++ {
//				output = append(output, destinList[i])
//			}
//			break
//		} else {

//			// generate new try
//			sTry = NewRndInd(curSrcSubs, searchParameters, searchDomain, sourceTabu)
//			dTry = NewRndInd(curDstSubs, searchParameters, searchDomain, destinTabu)

//			//DEBUG
//			fmt.Println("Current Try")
//			fmt.Println(sTry)
//			fmt.Println(dTry)

//			// test if destination found
//			if sTry[0] == dTry[0] && sTry[1] == dTry[1] {
//				// assign initial output component
//				output = sourceList[:len(sourceList)-1]

//				// assign destination list to output
//				for i := 0; i < len(destinList); i++ {
//					output = append(output, destinList[i])
//				}
//				break
//			} else if sourceTabu.At(sTry[0], sTry[1]) == 0.0 && destinTabu.At(dTry[0], dTry[1]) == 0.0 {
//				continue
//			} else {
//				// append to last element of source list
//				sourceList = append(sourceList, sTry)

//				// insert to first element of destination list
//				destinList = append(destinList, nil)
//				copy(destinList[0+1:], destinList[0:])
//				destinList[0] = dTry
//			}

//			// DEBUG
//			fmt.Println("Source & Distance Tabu Matrices")
//			for i := 0; i < 5; i++ {
//				fmt.Println(sourceTabu.RawRowView(i))
//			}
//			fmt.Println(" ")
//			for i := 0; i < 5; i++ {
//				fmt.Println(destinTabu.RawRowView(i))
//			}

//			// check validity of current tabu matrix
//			testSrcTabu := ValidateSubDomain(sTry, searchParameters.DstSubs, sourceTabu)
//			testDstTabu := ValidateSubDomain(dTry, searchParameters.SrcSubs, destinTabu)

//			// DEBUG
//			fmt.Println("Tabu Matrix Validity")
//			fmt.Println(testSrcTabu)
//			fmt.Println(testDstTabu)
//			fmt.Println("Iteration Count")
//			//fmt.Println(iter)

//			// if tabu matrix is invalid reset and restart
//			if testSrcTabu == false {

//				// DEBUG
//				fmt.Println("Reset Tabu Matrix")

//				// initialize chromosome as empty 2D slice with source subs as lead
//				sourceList = make([][]int, 1, searchDomain.MaxLen)
//				sourceList[0] = make([]int, 2)
//				sourceList[0][0] = searchParameters.SrcSubs[0]
//				sourceList[0][1] = searchParameters.SrcSubs[1]

//				// initialize new tabu matrix
//				sourceTabu = mat64.NewDense(searchDomain.Rows, searchDomain.Cols, nil)
//				sourceTabu.Clone(searchDomain.Matrix)
//				//tabu.Set(searchParameters.SrcSubs[0], searchParameters.SrcSubs[1], 0.0)
//			} else {
//				sourceTabu.Set(sTry[0], sTry[1], 0.0)
//			}

//			// if tabu matrix is invalid reset and restart
//			if testDstTabu == false {

//				// DEBUG
//				fmt.Println("Reset Tabu Matrix")

//				// initialize chromosome as empty 2D slice with source subs as lead
//				destinList = make([][]int, 1, searchDomain.MaxLen)
//				destinList[0] = make([]int, 2)
//				destinList[0][0] = searchParameters.DstSubs[0]
//				destinList[0][1] = searchParameters.DstSubs[1]

//				// initialize new tabu matrix
//				destinTabu = mat64.NewDense(searchDomain.Rows, searchDomain.Cols, nil)
//				destinTabu.Clone(searchDomain.Matrix)
//				//tabu.Set(searchParameters.SrcSubs[0], searchParameters.SrcSubs[1], 0.0)
//			} else {
//				destinTabu.Set(dTry[0], dTry[1], 0.0)
//			}
//		}
//	}

//	// return output
//	return output

//}
