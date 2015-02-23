// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"math/rand"
	"time"

	"github.com/gonum/matrix/mat64"
)

// new test domain initialization function
func NewToyDomain(identifier, rows, cols int) *Domain {

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

	// return output
	return &Domain{
		Id:     identifier,
		Matrix: domainMatrix,
	}
}

func NewToyObjective(identifier, rows, cols int) *Objective {

	// initialize empty matrix
	objectiveSize := rows * cols
	mat := make([]float64, objectiveSize)
	objectiveMatrix := mat64.NewDense(rows, cols, mat)

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// loop through matrix indices and assign random objective values
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			objectiveMatrix.Set(i, j, rand.Float64())
		}
	}

	return &Objective{
		Id:     identifier,
		Matrix: objectiveMatrix,
	}

}
