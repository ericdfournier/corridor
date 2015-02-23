// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"fmt"

	"github.com/gonum/matrix/mat64"
)

func ViewSearchDomain(searchDomain *Domain) {

	// get search domain matrix dimensions
	rows, _ := searchDomain.Matrix.Dims()

	// print domain id
	fmt.Printf("Search Domain ID = %d\n", searchDomain.Id)

	// print domain values to command line
	fmt.Printf("Search Domain Values = \n")
	for i := 0; i < rows; i++ {
		rawRowVals := searchDomain.Matrix.RawRowView(i)
		fmt.Println(rawRowVals)
	}
}

func ViewIndividual(searchDomain *Domain, searchParameters *Parameters, newIndividual *Individual) {

	// get search domain matrix dimensions and empty value slice
	rows, cols := searchDomain.Matrix.Dims()
	domainSize := rows * cols
	v := make([]float64, domainSize)

	// allocate new empty matrix
	blankMat := mat64.NewDense(rows, cols, v)

	// assign individual values to the empty matrix
	for i := 0; i < len(newIndividual.Subs); i++ {
		blankMat.Set(newIndividual.Subs[i][0], newIndividual.Subs[i][1], 1.0)
	}

	// print individual values to command line
	fmt.Printf("Individual = \n")
	for i := 0; i < rows; i++ {
		rawRowVals := blankMat.RawRowView(i)
		fmt.Println(rawRowVals)
	}
}
