// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/gonum/matrix/mat64"
)

// new sample parameters initialization function
func NewSampleParameters(searchDomain *Domain) *Parameters {

	// initialize variables
	sourceSubscripts := make([]int, 2)
	sourceSubscripts[0] = 3
	sourceSubscripts[1] = 3
	destinationSubscripts := make([]int, 2)
	destinationSubscripts[0] = searchDomain.Rows - 3
	destinationSubscripts[1] = searchDomain.Cols - 3
	randomnessCoefficient := 1.0
	populationSize := 1000
	selectionFraction := 0.5
	selectionProbability := 0.8
	mutationCount := 1
	mutationFraction := 0.2
	evolutionSize := 1000
	maxConcurrency := runtime.NumCPU()

	// return output
	return &Parameters{
		SrcSubs: sourceSubscripts,
		DstSubs: destinationSubscripts,
		RndCoef: randomnessCoefficient,
		PopSize: populationSize,
		SelFrac: selectionFraction,
		SelProb: selectionProbability,
		MutaCnt: mutationCount,
		MutaFrc: mutationFraction,
		EvoSize: evolutionSize,
		ConSize: maxConcurrency,
	}
}

// new sample domain initialization function
func NewSampleDomain(rows, cols int) *Domain {

	// initialize empty matrix
	domainSize := rows * cols
	mat := make([]float64, domainSize)
	domainMatrix := mat64.NewDense(rows, cols, mat)

	// loop through index values togo define domain
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if i == 0 {
				domainMatrix.Set(i, j, 0.0)
			} else if i == rows-1 {
				domainMatrix.Set(i, j, 0.0)
			} else if j == 0 {
				domainMatrix.Set(i, j, 0.0)
			} else if j == cols-1 {
				domainMatrix.Set(i, j, 0.0)
			} else {
				domainMatrix.Set(i, j, 1.0)
			}
		}
	}

	// compute band count
	bandCount := 2 + (int(math.Floor(math.Sqrt(math.Pow(float64(rows), 2.0)+math.Pow(float64(cols), 2.0)))) / 142)

	// return output
	return &Domain{
		Rows:   rows,
		Cols:   cols,
		Matrix: domainMatrix,
		BndCnt: bandCount,
	}
}

// new sample mutation domain initialization function
func NewSampleMutationDomain() *Domain {

	// fix domain size
	var rows int = 5
	var cols int = 5

	// initialize empty matrix
	domainSize := rows * cols
	mat := make([]float64, domainSize)
	domainMatrix := mat64.NewDense(rows, cols, mat)

	// loop through index values togo define domain
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if i == 0 {
				domainMatrix.Set(i, j, 0.0)
			} else if i == rows-1 {
				domainMatrix.Set(i, j, 0.0)
			} else if j == 0 {
				domainMatrix.Set(i, j, 0.0)
			} else if j == cols-1 {
				domainMatrix.Set(i, j, 0.0)
			} else {
				domainMatrix.Set(i, j, 1.0)
			}
		}
	}

	// eliminate center
	domainMatrix.Set(2, 2, 0.0)

	// set band count to nil
	var bandCount int = 2

	// return output
	return &Domain{
		Rows:   rows,
		Cols:   cols,
		Matrix: domainMatrix,
		BndCnt: bandCount,
	}
}

// new sample objective initialization function
func NewSampleObjectives(rows, cols, objectiveCount int) *MultiObjective {

	// initialize matrix dimensions
	objectiveSize := rows * cols
	var objectiveId int = 0

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// initialize empty objective slice
	objSlice := make([]*Objective, objectiveCount)

	// loop through matrix indices and assign random objective values
	for k := 0; k < objectiveCount; k++ {

		// initialize empty objective matrix
		mat := make([]float64, objectiveSize)
		objMat := mat64.NewDense(rows, cols, mat)

		// write random objective values
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				objMat.Set(i, j, math.Abs(rand.Float64()))
			}
		}

		// write to objective slice
		objSlice[k] = NewObjective(objectiveId, objMat)

		// iterate objective id
		objectiveId += 1
	}

	return &MultiObjective{
		ObjectiveCount: objectiveCount,
		Objectives:     objSlice,
	}
}
