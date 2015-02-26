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
func AllMinDistance(aSubs, bSubs []int, searchDomain *mat64.Dense) (allDistMatrix *mat64.Dense) {

	// get matrix dimensions
	rows, cols := searchDomain.Dims()

	// initialize new output matrix
	output := mat64.NewDense(rows, cols, nil)

	// loop through all values and compute minimum distances
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			pSubs := make([]int, 2)
			pSubs[0] = i
			pSubs[1] = j
			curMinDist := MinDistance(pSubs, aSubs, bSubs)
			output.Set(pSubs[0], pSubs[1], curMinDist)
		}
	}

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

// fitness function to generate the total fitness and chromosome
// fitness values for a given input chromosome
func ChromosomeFitness(inputChromosome *Chromosome, inputObjective *Objective) (outputChromosome *Chromosome) {

	// get chromosome length
	chromLen := len(inputChromosome.Subs)

	// evaluate chromosome fitness according to input objective
	for i := 0; i < chromLen; i++ {
		curFit := inputObjective.Matrix.At(inputChromosome.Subs[i][0], inputChromosome.Subs[i][1])
		inputChromosome.Fitness[i] = curFit
		inputChromosome.TotalFitness = inputChromosome.TotalFitness + curFit
	}

	// return outputs
	return inputChromosome
}

// THERE IS SOMETHING WRONG HERE...THE CHANNEL SEEMS TO NOT
// BE OUTPUTTING CHROMOSOME VALUES PROPERLY....

func AccumFitness(chromosomes chan *Chromosome) (cumulativeFitness float64) {

	// initialize output
	var output float64

	// drain channel to accumulate fitness values
	for i := 0; i < cap(chromosomes); i++ {
		curChrom := <-chromosomes
		output = output + curChrom.TotalFitness
	}

	// return output
	return output
}

// fitness function generate the mean and standard deviation of
// fitness values for all of the chromosomes in a given population
func PopulationFitness(inputPopulation *Population, inputObjective *Objective) (outputPopulation *Population) {

	// generate cumulative fitness
	cumFit := AccumFitness(inputPopulation.Chromosomes)

	// initialize pop size
	popSize := cap(inputPopulation.Chromosomes)

	// compute mean from cumulative
	inputPopulation.MeanFitness = cumFit / float64(popSize)

	// return output
	return inputPopulation
}
