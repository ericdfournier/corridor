// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"fmt"

	"github.com/gonum/matrix/mat64"
)

// function to print the properties of a search domain to the command line
func ViewDomain(searchDomain *Domain) {

	// get search domain matrix dimensions
	rows, _ := searchDomain.Matrix.Dims()

	// print domain id
	fmt.Printf("Search Domain ID = %d\n", searchDomain.Id)

	// print domain values to command line
	fmt.Printf("Search Domain Values = \n")
	for i := 0; i < rows; i++ {
		rawRowVals := searchDomain.Matrix.RawRowView(i)
		fmt.Printf("%1.0f\n", rawRowVals)
	}
}

// function to print the properties of a basis solution to the command line
func ViewBasis(basisSolution *Basis) {

	// get basis solution matrix dimensions
	rows, _ := basisSolution.Matrix.Dims()

	// print domain id
	fmt.Printf("Basis ID = %d\n", basisSolution.Id)

	// print domain values to command line
	fmt.Printf("Basis Solution Values = \n")
	for i := 0; i < rows; i++ {
		rawRowVals := basisSolution.Matrix.RawRowView(i)
		fmt.Printf("%1.0f\n", rawRowVals)
	}
}

// function to print the properties of a chromosome to the command line
func ViewChromosome(searchDomain *Domain, searchParameters *Parameters, inputChromosome *Chromosome) {

	// get search domain matrix dimensions and empty value slice
	rows, cols := searchDomain.Matrix.Dims()
	domainSize := rows * cols
	v := make([]float64, domainSize)

	// allocate new empty matrix
	blankMat := mat64.NewDense(rows, cols, v)

	// assign chromosome values to the empty matrix
	for i := 0; i < len(inputChromosome.Subs); i++ {
		blankMat.Set(inputChromosome.Subs[i][0], inputChromosome.Subs[i][1], 1.0)
	}

	// print chromosome values to command line
	fmt.Printf("Chromosome = \n")
	for i := 0; i < rows; i++ {
		rawRowVals := blankMat.RawRowView(i)
		fmt.Printf("%1.0f\n", rawRowVals)
	}

	// print output to the command line
	fmt.Printf("Chromosome Length = %d\n", len(inputChromosome.Subs))
	fmt.Printf("Chromosome Total Fitness = %1.5f\n", inputChromosome.TotalFitness)
}

// functions to print the frequency of chromosomes in a search domain to the command line
func ViewPopulation(searchDomain *Domain, searchParameters *Parameters, inputPopulation *Population) {

	// get search domain dimensions
	rows, cols := searchDomain.Matrix.Dims()

	// allocate new empty matrix
	mat := mat64.NewDense(rows, cols, nil)

	// accumulated visited subscripts in new empty matrix
	for i := 0; i < searchParameters.PopSize; i++ {

		// extract current chromosome from channel
		curChrom := <-inputPopulation.Chromosomes
		curInd := curChrom.Subs
		lenCurInd := len(curInd)

		// iterate over subscript indices
		for j := 0; j < lenCurInd; j++ {
			curSubs := curInd[j]
			curVal := mat.At(curSubs[0], curSubs[1])
			newVal := curVal + 1
			mat.Set(curSubs[0], curSubs[1], newVal)
		}

		// repopulate channel
		inputPopulation.Chromosomes <- curChrom
	}

	// print matrix values to command line
	fmt.Printf("Population Size = %d\n", searchParameters.PopSize)
	fmt.Printf("Population Frequency = \n")
	for q := 0; q < rows; q++ {
		rawRowVals := mat.RawRowView(q)
		fmt.Printf("%*.0f\n", DigitCount(searchParameters.PopSize)+1, rawRowVals)
	}
}
