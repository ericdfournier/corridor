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
		fmt.Printf("%4.0f\n", rawRowVals)
	}
}

func ViewChromosome(searchDomain *Domain, searchParameters *Parameters, newChromosome *Chromosome) {

	// get search domain matrix dimensions and empty value slice
	rows, cols := searchDomain.Matrix.Dims()
	domainSize := rows * cols
	v := make([]float64, domainSize)

	// allocate new empty matrix
	blankMat := mat64.NewDense(rows, cols, v)

	// assign chromosome values to the empty matrix
	for i := 0; i < len(newChromosome.Subs); i++ {
		blankMat.Set(newChromosome.Subs[i][0], newChromosome.Subs[i][1], 1.0)
	}

	// print chromosome values to command line
	fmt.Printf("Chromosome = \n")
	for i := 0; i < rows; i++ {
		rawRowVals := blankMat.RawRowView(i)
		fmt.Printf("%4.0f\n", rawRowVals)
	}
	fmt.Printf("Chromosome Length = %d\n", len(newChromosome.Subs))
	fmt.Printf("Chromosome Total Fitness = %1.5f\n", newChromosome.TotalFitness)
}

func ViewPopulation(searchDomain *Domain, searchParameters *Parameters, inputPopulation *Population) {

	// get search domain matrix dimensions and empty value slice
	popSize := searchParameters.PopSize

	// get search domain dimensions
	rows, cols := searchDomain.Matrix.Dims()

	// allocate new empty matrix
	mat := mat64.NewDense(rows, cols, nil)

	// accumulated visited subscripts in new empty matrix
	for i := 0; i < popSize; i++ {
		curChrom := <-inputPopulation.Chromosomes
		curInd := curChrom.Subs
		lenCurInd := len(curInd)
		for j := 0; j < lenCurInd; j++ {
			curSubs := curInd[j]
			curVal := mat.At(curSubs[0], curSubs[1])
			newVal := curVal + 1
			mat.Set(curSubs[0], curSubs[1], newVal)
		}

		// repopulate channel
		//inputPopulation.Chromosomes <- curChrom
	}

	// print matrix values to command line
	fmt.Printf("Population Size = %d\n", searchParameters.PopSize)
	fmt.Printf("Population Frequency = \n")
	for q := 0; q < rows; q++ {
		rawRowVals := mat.RawRowView(q)
		fmt.Printf("%4.0f\n", rawRowVals)
	}

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
		fmt.Printf("%4.0f\n", rawRowVals)
	}
}
