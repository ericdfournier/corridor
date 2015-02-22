// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import "github.com/gonum/matrix/mat64"

// sub2ind converts a single pair of input row and column references to a 2-D
// matrix with a given stride and size into a corresponding linear index
func Sub2Ind(row, col, stride, size int) (index int) {

}

// ind2sub converts a linear index value to a 2-D matrix with a given stride
// and size into a corresponding pair of row and column subscripts
func Ind2Sub(index, stride, size int) (row, col int) {

}

// mat2bnd takes an input dense matrix and determines the location of the
// boundary indexs returning them as a vector slice
func Mat2Bnd(matrix *mat64.Dense) (boundaryIndices map[int]bool) {

}
