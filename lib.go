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

func MinDistance(pSubs, aSubs, bSubs []int) (minDist float64) {

	var x float64 = float64(pSubs[0])
	var y float64 = float64(pSubs[1])
	var x0 float64 = float64(aSubs[0])
	var y0 float64 = float64(aSubs[1])
	var x1 float64 = float64(bSubs[0])
	var y1 float64 = float64(bSubs[1])

	a := x - x0
	b := y - y0
	c := x1 - x0
	d := y1 - y0

	dot := a*c + b*d
	lenSq := c*c + d*d

	var param float64 = -1.0

	if lenSq != 0 {
		param = dot / lenSq
	}

	var xx, yy float64

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

	var dx float64 = x - xx
	var dy float64 = y - yy

	output := math.Sqrt(dx*dx + dy*dy)
	// return final output
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
