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
	fmt.Printf("Individual Length = %d\n", len(newIndividual.Subs))
	fmt.Printf("Individual Total Fitness = %1.5f\n", newIndividual.TotalFitness)
}

func ViewPopulation(searchDomain *Domain, searchParameters *Parameters, newPopulation *Population) {

	// get search domain matrix dimensions and empty value slice
	popSize := searchParameters.PopSize

	// get search domain dimensions
	rows, cols := searchDomain.Matrix.Dims()

	// allocate new empty matrix
	mat := mat64.NewDense(rows, cols, nil)

	// extract all individuals
	allIndiv := *newPopulation.Individuals

	// accumulated visited subscripts in new empty matrix
	for i := 0; i < popSize; i++ {
		curInd := allIndiv[i].Subs
		lenCurInd := len(curInd)
		for j := 0; j < lenCurInd; j++ {
			curSubs := curInd[j]
			curVal := mat.At(curSubs[0], curSubs[1])
			newVal := curVal + 1
			mat.Set(curSubs[0], curSubs[1], newVal)
		}
	}

	// print matrix values to command line
	fmt.Printf("Population Frequency = \n")
	for q := 0; q < rows; q++ {
		rawRowVals := mat.RawRowView(q)
		fmt.Println(rawRowVals)
	}
	fmt.Printf("Population Size = %d\n", searchParameters.PopSize)

}

func ViewBasis(basisSolution *Basis) {

	// get basis solution matrix dimensions
	rows, _ := basisSolution.Matrix.Dims()

	// print domain id
	fmt.Printf("Basis ID = %d\n", basisSolution.Id)

	// print domain values to command line
	fmt.Printf("Basis Solution Values = \n")
	for i := 0; i < rows; i++ {
		rawRowVals := basisSolution.Matrix.RawRowView(i)
		fmt.Println(rawRowVals)
	}
}
