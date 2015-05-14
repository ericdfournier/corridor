// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"errors"
	"math"

	"github.com/gonum/matrix/mat64"
)

// compute euclidean distance for a pair of subscript indices
func Distance(aSubs, bSubs []int) (dist float64) {

	// initialize variables
	var x0 float64 = float64(aSubs[0])
	var x1 float64 = float64(bSubs[0])
	var y0 float64 = float64(aSubs[1])
	var y1 float64 = float64(bSubs[1])
	var pow float64 = 2.0
	var dx float64 = x1 - x0
	var dy float64 = y1 - y0

	// compute distance
	var output float64 = math.Sqrt(math.Pow(dx, pow) + math.Pow(dy, pow))

	// return final output
	return output
}

// alldistance computes the distance from each location with the input
// search domain and a given point defined by an input pair of row
// column subscripts
func AllDistance(aSubs []int, searchDomain *mat64.Dense) (allDistMatrix *mat64.Dense) {

	// get matrix dimensions
	rows, cols := searchDomain.Dims()

	// initialize new output matrix
	output := mat64.NewDense(rows, cols, nil)

	// initialize destination point subscript slice
	bSubs := make([]int, 2)

	// loop through all values and compute minimum distances
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			bSubs[0] = i
			bSubs[1] = j
			output.Set(bSubs[0], bSubs[1], Distance(aSubs, bSubs))
		}
	}

	// return output
	return output

}

// compute the minimum distance between a given input point and
// the subscripts comprised of a line segement joining two other
// input points
func MinDistance(pSubs, aSubs, bSubs []int) (minDist float64) {

	// initialize variables
	var x float64 = float64(pSubs[0])
	var y float64 = float64(pSubs[1])
	var x0 float64 = float64(aSubs[0])
	var y0 float64 = float64(aSubs[1])
	var x1 float64 = float64(bSubs[0])
	var y1 float64 = float64(bSubs[1])

	// compute difference components
	a := x - x0
	b := y - y0
	c := x1 - x0
	d := y1 - y0

	// compute dot product of difference components
	dot := a*c + b*d
	lenSq := c*c + d*d

	// initialize parameter
	var param float64 = -1.0

	// if zero length condition
	if lenSq != 0 {
		param = dot / lenSq
	}

	// initialize transform variables
	var xx, yy float64

	// switch transform mechanism on orientation
	if param < 0 {
		xx = x0
		yy = y0
	} else if param > 1 {
		xx = x1
		yy = y1
	} else {
		xx = x0 + param*c
		yy = y0 + param*d
	}

	// execute transform
	var dx float64 = x - xx
	var dy float64 = y - yy

	// generate output
	output := math.Sqrt(dx*dx + dy*dy)

	// return final output
	return output
}

// allmindistance computes the distance from each location within the
// input search domain and to the nearest subscript located along the
// line formed by the two input subscripts
func AllMinDistance(aSubs, bSubs []int, searchDomainMatrix *mat64.Dense) (allMinDistMatrix *mat64.Dense) {

	// get matrix dimensions
	rows, cols := searchDomainMatrix.Dims()

	// initialize new output matrix
	output := mat64.NewDense(rows, cols, nil)

	// initialize slice
	pSubs := make([]int, 2)

	// loop through all values and compute minimum distances
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			pSubs[0] = i
			pSubs[1] = j
			curMinDist := MinDistance(pSubs, aSubs, bSubs)
			output.Set(pSubs[0], pSubs[1], curMinDist)
		}
	}

	// return final output
	return output
}

// distancebands recodes a distance matrix computed from a single
// source location to ordinal set of bands of increasing distance
func DistanceBands(bandCount int, distanceMatrix *mat64.Dense) (bandMatrix *mat64.Dense) {

	// get matrix dimensions
	rows, cols := distanceMatrix.Dims()

	// initialize output
	output := mat64.NewDense(rows, cols, nil)

	// check band count against input distance matrix size
	if bandCount > rows || bandCount > cols {
		err := errors.New("Input band count too large for input distance matrix \n")
		panic(err)
	}

	// generate band range
	minDist := distanceMatrix.Min()
	maxDist := distanceMatrix.Max()

	// initialize band interval unit and slice
	bandUnit := (maxDist - minDist) / float64(bandCount+1)
	bandInt := make([]float64, bandCount+1)

	// generate band intervals
	for i := 0; i < bandCount+1; i++ {
		if i == 0 {
			bandInt[i] = 0
		} else {
			bandInt[i] = bandInt[i-1] + bandUnit
		}
	}

	// perform conversion to the appropriate band interval
	for i := 0; i < len(bandInt)-1; i++ {
		for j := 0; j < rows; j++ {
			for k := 0; k < rows; k++ {
				if distanceMatrix.At(j, k) > bandInt[i] && distanceMatrix.At(j, k) < bandInt[i+1] {
					output.Set(j, k, float64(i))
				} else if distanceMatrix.At(j, k) > bandInt[i+1] {
					output.Set(j, k, float64(i+1))
				}
			}
		}
	}

	// return output
	return output
}

// bandmask selects the elements in a distance band matrix
// corresponding to a specified input band identification number
// and outputs a binary matrix of the same dimensions as the distance
// band matrix with the values at those locations encoded as ones
// and all other locations encoded as zeros
func BandMask(bandValue float64, bandMatrix *mat64.Dense) (binaryBandMat *mat64.Dense) {

	// get row column dimensions of band matrix
	rows, cols := bandMatrix.Dims()

	// initialize output
	output := mat64.NewDense(rows, cols, nil)

	// loop through matrix values and perform binary encoding
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {

			// perform elementwise equality test
			if bandValue == bandMatrix.At(i, j) {
				output.Set(i, j, 1.0)
			} else {
				output.Set(i, j, 0.0)
			}
		}
	}

	// return output
	return output
}

// nonzerosubs returns a 2-D slice containing the row column indices
// of all nonzero elements contained wihtin a given input matrix
func NonZeroSubs(inputMatrix *mat64.Dense) (nonZeroSubs [][]int) {

	// get matrix dimensions
	rows, cols := inputMatrix.Dims()

	// initialize output
	output := make([][]int, 1)
	output[0] = make([]int, 2)

	// initialize iterator and current subscript slice
	var iter int = 0

	// loop through and check values
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {

			// test for non-zero values
			if inputMatrix.At(i, j) != 0.0 {
				if iter == 0 {
					output[iter] = []int{i, j}
					iter += 1
				} else if iter > 0 {
					output = append(output, []int{i, j})
					iter += 1
				}
			} else if inputMatrix.At(i, j) == 0.0 {
			}
		}
	}

	// return output
	return output
}

// orientation accepts as inputs a pair of point subscripts
// and returns a binary vector indicating the relative orientation
// of the first point to the second in binary terms
func Orientation(aSubs, bSubs []int) (orientationVector []int) {

	// initialize output
	output := make([]int, 2)

	// generate reference orientation row parameter
	if aSubs[0]-bSubs[0] < 0 {
		output[0] = 1
	} else if aSubs[0]-bSubs[0] == 0 {
		output[0] = 0
	} else if aSubs[0]-bSubs[0] > 0 {
		output[0] = -1
	}

	// generate reference orientation column parameter
	if aSubs[1]-bSubs[1] < 0 {
		output[1] = 1
	} else if aSubs[1]-bSubs[1] == 0 {
		output[1] = 0
	} else if aSubs[1]-bSubs[1] > 0 {
		output[1] = -1
	}

	// return output
	return output
}

// orientation mask returns a binary encoded matrix for
// a given point where all points orientated towards
// a given second point are encoded as 1 and all points
// orientated away from the given second point as 0
func OrientationMask(aSubs, bSubs []int, searchDomainMatrix *mat64.Dense) (orientationMask *mat64.Dense) {

	// generate matrix dimensions
	rows, cols := searchDomainMatrix.Dims()

	// initialize output matrix
	output := mat64.NewDense(rows, cols, nil)

	// generate reference orientation vectors
	sRefOrientVec := Orientation(aSubs, bSubs)
	dRefOrientVec := Orientation(bSubs, aSubs)

	// initialize current subs and orientation vectors
	curSubs := make([]int, 2)
	sOrientVec := make([]int, 2)
	dOrientVec := make([]int, 2)

	// loop through domain matrix and generate orientation matrix values
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {

			// compute current orientation
			curSubs[0] = i
			curSubs[1] = j
			sOrientVec = Orientation(curSubs, bSubs)
			dOrientVec = Orientation(curSubs, aSubs)

			// check for match and assign values
			if sOrientVec[0] == sRefOrientVec[0] && sOrientVec[1] == sRefOrientVec[1] {
				if dOrientVec[0] == dRefOrientVec[0] && dOrientVec[1] == dRefOrientVec[1] {
					output.Set(i, j, 1.0)
				}
			}
		}
	}

	// return output
	return output
}

// bresenham generates the list of subscript indices corresponding to the
// euclidean shortest paths connecting two subscript pairs in discrete space
func Bresenham(aSubs, bSubs []int) (lineSubs [][]int) {

	// initialize variables
	var x0 int = aSubs[0]
	var x1 int = bSubs[0]
	var y0 int = aSubs[1]
	var y1 int = bSubs[1]

	// check row differential
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}

	// check column differential
	dy := y1 - y0

	// if differential is negative flip
	if dy < 0 {
		dy = -dy
	}

	// initialize stride variables
	var sx, sy int

	// set row stride direction
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	// set column stride direction
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	// calculate error component
	err := dx - dy

	// initialize output 2D slice vector
	dist := math.Ceil(Distance(aSubs, bSubs))
	maxLen := int(dist)
	output := make([][]int, 0, maxLen)

	// loop through and generate subscripts
	for {
		var val = []int{x0, y0}
		output = append(output, val)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}

	// return final output
	return output
}

// function to return the subscript indices of the cells corresponding to the
// queens neighborhood for a given subscript pair
func NeighborhoodSubs(row, col int) (subs [9][2]int) {

	// initialize output slice
	var output [9][2]int

	// write neighborhood subscript values
	output[0][0] = row - 1
	output[0][1] = col - 1
	output[1][0] = row - 1
	output[1][1] = col
	output[2][0] = row - 1
	output[2][1] = col + 1
	output[3][0] = row
	output[3][1] = col - 1
	output[4][0] = row
	output[4][1] = col
	output[5][0] = row
	output[5][1] = col + 1
	output[6][0] = row + 1
	output[6][1] = col - 1
	output[7][0] = row + 1
	output[7][1] = col
	output[8][0] = row + 1
	output[8][1] = col + 1

	return output
}

// function to validate an input sub domain for use in generating
// a chromosomal mutation via the random walk procedure
func ValidateSubDomain(subSource, subDestin []int, subMat *mat64.Dense) bool {

	// initialize output
	var output bool

	// generate sub source neighborhood
	sNeigh := NeighborhoodSubs(subSource[0], subSource[1])

	// generate sub destination neighborhood
	dNeigh := NeighborhoodSubs(subDestin[0], subDestin[1])

	// generate center row
	centerRow := subMat.RowView(2)

	// generate center column
	centerCol := subMat.ColView(2)

	// initialize summation variables
	var sSum float64 = 0.0
	var dSum float64 = 0.0
	var rSum float64 = 0.0
	var cSum float64 = 0.0

	// enter for loop for start and destination sums
	for i := 0; i < 9; i++ {
		sSum = sSum + subMat.At(sNeigh[i][0], sNeigh[i][1])
		dSum = dSum + subMat.At(dNeigh[i][0], dNeigh[i][1])
	}

	// enter for loop for row column sums
	for j := 0; j < 5; j++ {
		rSum = rSum + centerRow.At(j, 0)
		cSum = cSum + centerCol.At(j, 0)
	}

	// check conditions to validate neighborhood
	if sSum <= 1.0 || dSum <= 1.0 || rSum == 0.0 || cSum == 0.0 {
		output = false
	} else {
		output = true
	}

	//return final output
	return output
}

// function validate the tabu neighborhood of an input pair of
// row column subscripts
func ValidateTabu(currentSubs []int, tabuMatrix *mat64.Dense) bool {

	// initialize output
	var output bool

	// initialize tabu neighborhood sum
	var tSum int = 0

	// generate neighborhood subscripts
	tNeigh := NeighborhoodSubs(currentSubs[0], currentSubs[1])

	// loop through and compute sum
	for i := 0; i < 9; i++ {
		if i != 4 {
			tSum += int(tabuMatrix.At(tNeigh[i][0], tNeigh[i][1]))
		}
	}

	// write output boolean
	if tSum == 0 {
		output = false
	} else {
		output = true
	}

	// return output
	return output
}

// function to count the number of digits in an input integer as
// its base ten logarithm
func DigitCount(input int) (digits int) {

	// compute digits as the log of the input
	output := int(math.Floor(math.Log10(float64(input))))

	// return output
	return output
}
