// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gonum/matrix/mat64"
)

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

// fitness function generate the mean and standard deviation of
// fitness values for all of the chromosomes in a given population
func PopulationFitness(inputPopulation *Population, inputObjective *Objective) (outputPopulation *Population) {

	// initialize pop size
	popSize := cap(inputPopulation.Chromosomes)

	// initialize output
	var output float64

	// drain channel to accumulate fitness values
	for i := 0; i < popSize; i++ {

		// read current chromosome from channel
		curChrom := <-inputPopulation.Chromosomes

		// launch go thread to compute cumulative fitness
		go func(curChrom *Chromosome) {
			output = curChrom.TotalFitness + output
		}(curChrom)

		// recieve from channel
		inputPopulation.Chromosomes <- curChrom
	}

	// compute mean from cumulative
	inputPopulation.MeanFitness = output / float64(popSize)

	// return output
	return inputPopulation
}

// selection operator selects between two chromosomes with a
// probability of the most fit chromosome being selected
// determined by the input selection probability ratio
func ChromosomeSelection(chrom1, chrom2 *Chromosome, selectionProb float64) (selectedChrom *Chromosome) {

	// initialize output
	output := chrom1

	// get current time for random number seed
	rand.Seed(time.Now().UnixNano())

	// generate random number to determine selection result
	dec := rand.Float64()

	// perform conditional selection
	if dec > selectionProb { // normal
		if chrom1.TotalFitness > chrom2.TotalFitness {
			output = chrom1
		} else {
			output = chrom2
		}
	} else { // inverted
		if chrom1.TotalFitness > chrom2.TotalFitness {
			output = chrom2
		} else {
			output = chrom1
		}
	}

	// return output
	return output
}

// population selection operator selects half of the input
// population for reproduction based upon comparative
// fitness and some randomized input selection fraction
func PopulationSelection(inputPopulation *Population, inputParameters *Parameters) (selection chan *Chromosome) {

	// initialize selection channel size
	selSize := int(math.Floor(float64(cap(inputPopulation.Chromosomes)) * inputParameters.SelFrac))

	// initialize selection channel
	output := make(chan *Chromosome, selSize)

	// initialize selection probability
	selProb := inputParameters.SelProb

	// initialize selection loop
	for i := 0; i < selSize; i++ {

		// write selection to output channel
		output <- ChromosomeSelection(<-inputPopulation.Chromosomes, <-inputPopulation.Chromosomes, selProb)
	}

	// return selection channel
	return output
}

// intersection determines whether or not the subscripts
// associated with two input chromosomes share any in
// values in common and reports their relative indices
func ChromosomeIntersection(subs1, subs2 [][]int) (subs1Indices, subs2Indices []int) {

	// initialize output index slice
	output1 := make([]int, 0)
	output2 := make([]int, 0)

	// initialize subscript lengths
	len1 := len(subs1)
	len2 := len(subs2)

	// enter intersection loop
	for i := 0; i < len1; i++ {
		for j := 0; j < len2; j++ {
			if subs1[i][0] == subs2[j][0] && subs1[i][1] == subs2[j][1] {
				output1 = append(output1, i)
				output2 = append(output2, j)
			} else {
			}
		}
	}

	// return intersection output
	return output1, output2

}

// crossover operator performs the single point crossover
// operation for two input chromosomes that have
// previously been selected from a source population
func ChromosomeCrossover(chrom1Ind, chrom2Ind []int, chrom1Subs, chrom2Subs [][]int) (crossoverChrom [][]int) {

	// initialize maximum length
	maxLen := len(chrom1Subs) + len(chrom2Subs)

	// initialize output
	output := make([][]int, 0, maxLen)

	// get current time for random number seed
	rand.Seed(time.Now().UnixNano())

	var r int

	// generate random number to determine selection result
	// while screening out initial source index match
	for {
		r = rand.Intn(len(chrom1Ind) - 1)
		if r == 0 {
			continue
		} else {
			break
		}
	}

	// generate subscript slice 1
	for i := 0; i < (chrom1Ind[r] + 1); i++ {
		output = append(output, chrom1Subs[i])
	}

	// generate subscript slice 2
	for j := (chrom2Ind[r] + 1); j < len(chrom2Subs); j++ {
		output = append(output, chrom2Subs[j])
	}

	// return output
	return output

}

// selection crossover operator performs a single part
// crossover on each of the individuals provided in an
// input selection channel of chromosomes
func SelectionCrossover(inputSelection chan *Chromosome, inputParameters *Parameters, inputObjective *Objective, inputDomain *Domain) (crossover chan *Chromosome) {

	// initialize crossover channel size
	popSize := int(inputParameters.PopSize)

	// initialize crossover channel
	output := make(chan *Chromosome, popSize)

	// initialize crossover loop
	for i := 0; i < popSize; i++ {

		for {

			// extract chromosomes
			chrom1 := <-inputSelection
			chrom2 := <-inputSelection

			// initialize empty index slices
			var chrom1Ind []int
			var chrom2Ind []int

			// initialize empty chromosome
			empChrom := NewEmptyChromosome(inputDomain)

			// check for valid crossover point
			chrom1Ind, chrom2Ind = ChromosomeIntersection(chrom1.Subs, chrom2.Subs)

			// resample chromosomes if no intersection
			if len(chrom1Ind) > 2 {
				empChrom.Subs = ChromosomeCrossover(chrom1Ind, chrom2Ind, chrom1.Subs, chrom2.Subs)
				empChrom = ChromosomeFitness(empChrom, inputObjective)
				output <- empChrom
				inputSelection <- chrom1
				inputSelection <- chrom2
				break
			} else {
				inputSelection <- chrom2
				inputSelection <- chrom1
				continue
			}
		}
	}

	// return output
	return output
}

// mutationLocus to randomly select a mutation locus and return the adjacent
// loci along the length of the chromosome
func MutationLoci(inputChromosome *Chromosome) (previousLocus, mutationLocus, nextLocus []int, mutationIndex int) {

	// compute chromosome length
	lenChrom := len(inputChromosome.Subs)

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// randomly select mutation index
	mutIndex := rand.Intn(lenChrom-4) + 2

	// get mutation locui subscripts from mutIndex
	mutLocus := inputChromosome.Subs[mutIndex]
	prvLocus := inputChromosome.Subs[mutIndex-1]
	nxtLocus := inputChromosome.Subs[mutIndex+1]

	// return output
	return prvLocus, mutLocus, nxtLocus, mutIndex
}

// mutation sub domain returns the subdomain to be used for the mutation
// specific directed walk procedure
func MutationSubDomain(previousLocus, mutationLocus, nextLocus []int, inputDomain *mat64.Dense) (outputSubDomain *mat64.Dense) {

	// generate mutation locus neighborhood indices
	nInd := NeighborhoodSubs(mutationLocus[0], mutationLocus[1])

	// initialize iterator
	var n int = 0

	// initialize sub domain matrix
	subMat := mat64.NewDense(5, 5, nil)

	// clean sub domain
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if i == 0 {
				subMat.Set(i, j, 0.0)
			} else if i == 4 {
				subMat.Set(i, j, 0.0)
			} else if j == 0 {
				subMat.Set(i, j, 0.0)
			} else if j == 4 {
				subMat.Set(i, j, 0.0)
			} else if nInd[n][0] == previousLocus[0] && nInd[n][1] == previousLocus[1] {
				subMat.Set(i, j, 1.0)
				// iterate counter
				n += 1
			} else if nInd[n][0] == nextLocus[0] && nInd[n][1] == nextLocus[1] {
				subMat.Set(i, j, 1.0)
				// iterate counter
				n += 1
			} else {
				subMat.Set(i, j, inputDomain.At(nInd[n][0], nInd[n][1]))
				// iterate counter
				n += 1
			}
		}
	}

	// return output
	return subMat
}

// function to generate a mutation within a given chromosome at a specified
// number of mutation loci
func ChromosomeMutation(inputChromosome *Chromosome, inputDomain *Domain, inputParameters *Parameters) (outputChromosome *Chromosome) {

	// compute chromosome length
	lenChrom := len(inputChromosome.Subs)

	// initialize output chromosome
	output := inputChromosome

	// copy input domain matrix
	refDomain := inputDomain.Matrix

	// block out cells on current chromosome
	for k := 0; k < lenChrom; k++ {
		refDomain.Set(inputChromosome.Subs[k][0], inputChromosome.Subs[k][1], 0.0)
	}

	// generate mutation loci
	prvLocus, mutLocus, nxtLocus, mutIndex := MutationLoci(inputChromosome)

	// first check if deletion is valid, else perform mutation
	if Distance(prvLocus, nxtLocus) < 1.5 {

		// perform simple deletion of mutation index
		output.Subs = append(output.Subs[:mutIndex], output.Subs[(mutIndex+1):]...)
	} else {

		fmt.Println("Mutation Locus")
		fmt.Println(mutLocus)

		// generate mutation subdomain
		subMat := MutationSubDomain(prvLocus, mutLocus, nxtLocus, refDomain)

		// generate sub source and sub destination
		subSource := make([]int, 2)
		subDestin := make([]int, 2)
		subSource[0] = prvLocus[0] - mutLocus[0] + 2
		subSource[1] = prvLocus[1] - mutLocus[1] + 2
		subDestin[0] = nxtLocus[0] - mutLocus[0] + 2
		subDestin[1] = nxtLocus[1] - mutLocus[1] + 2

		fmt.Println("Sub Source")
		fmt.Println(subSource)
		fmt.Println("Sub Destination")
		fmt.Println(subDestin)

		// generate subdomain from sub matrix and generate sub basis
		subDomain := NewDomain(1, subMat)
		subParams := NewParameters(subSource, subDestin, 1.0, 1, 1.0, 1.0, 1)

		// generate directed walk based mutation
		subWlk := RndWlk(subDomain, subParams)

		fmt.Println("Raw Sub Walk Subscripts")
		fmt.Println(subWlk)

		// THERE IS A PROBLEM HERE WITH THE TRANSLATION OF THE SUBSCRIPTS
		// I HAVE TO FIGURE IT OUT TO MAKE THIS WORK...

		// translate subscripts
		for i := 0; i < len(subWlk); i++ {
			subWlk[i][0] = subWlk[i][0] - 2 + mutLocus[0]
			subWlk[i][1] = subWlk[i][1] - 2 + mutLocus[1]
		}

		fmt.Println("Translated Sub Walk Subscripts")
		fmt.Println(subWlk)

		// delete mutation locus
		output.Subs = append(output.Subs[:mutIndex], output.Subs[(mutIndex+1):]...)

		// insert sub walk section into original chromosome
		output.Subs = append(output.Subs[:mutIndex-1], append(subWlk, output.Subs[mutIndex+1:]...)...)

	}

	// return output
	return output
}

// function generate mutations within a specified fraction of an input
// population with those chromosomes being selected at random
func PopulationMutation(inputPopulation *Population, inputDomain *Domain, inputParameters *Parameters) (outputPopulation *Population) {

	// return output
	return
}

// population evolution operator generates a new population
// from an input population using the selection and crossover operators
func PopulationEvolution(inputPopulation *Population, inputDomain *Domain, inputParameters *Parameters, inputObjective *Objective) (outputPopulation *Population) {

	// initialize new empty population
	output := NewEmptyPopulation()

	// perform population selection
	popSel := PopulationSelection(inputPopulation, inputParameters)

	// perform selection crossover
	selCrs := SelectionCrossover(popSel, inputParameters, inputObjective, inputDomain)

	// fill empty population
	output.Chromosomes = selCrs

	// return output
	return output
}
