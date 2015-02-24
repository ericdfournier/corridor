// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

// compute euclidean distance for a pair of subscript indices
func Distance(aSubs, bSubs []int) (dist float64) {

	// initialize power variable
	var pow float64 = 2.0

	// initialize output variable
	var output float64

	// compute distance
	output = math.Sqrt(math.Pow(float64(aSubs[1]-aSubs[0]), pow) + math.Pow(float64(bSubs[0]-bSubs[1]), pow))

	// return final output

	return output
}

func Bresenham(aSubs, bSubs []int, searchDomain *mat64.Dense) (lineSubs [][]int) {

	var x0 = aSubs[0]
	var x1 = bSubs[0]
	var y0 = aSubs[1]
	var y1 = bSubs[1]

	dx := x1 - x0

	if dx < 0 {
		dx = -dx
	}

	dy := y1 - y0

	if dy < 0 {
		dy = -dy
	}

	var sx, sy int

	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy
	rows, cols := searchDomain.Dims()
	maxLen := rows * cols
	output := make([][]int, 1, maxLen)
	output[0] = make([]int, 2)
	val := make([]int, 2)

	for {
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
		val[0] = x0
		val[1] = y0
		output = append(output, val)
	}

	return output
}

// fitness function to generate the total fitness and individual
// fitness values for a given input set of subscripts
// corresponding to a single individual
func Fitness(subs [][]int, obj *mat64.Dense) (fitnessValues []float64, totalFitness float64) {

	// get individual length
	indSize := len(subs)

	// initialize fitness values and total fitness
	fitVal := make([]float64, indSize)
	var totFit float64 = 0.0

	// evaluate individual fitness according to input objective
	for i := 0; i < indSize; i++ {
		curFit := obj.At(subs[i][0], subs[i][1])
		fitVal[i] = curFit
		totFit = totFit + curFit
	}

	// return outputs
	return fitVal, totFit

}
