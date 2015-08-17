// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"math"
	"math/rand"
	"time"

	"github.com/gonum/matrix/mat64"
)

// fitness function to generate the total fitness and chromosome
// fitness values for a given input chromosome
func ChromosomeFitness(inputChromosome *Chromosome, inputObjectives *MultiObjective) (outputChromosome *Chromosome) {

	// get chromosome length
	chromLen := len(inputChromosome.Subs)

	// clear current chromosome fitness values
	inputChromosome.Fitness = make([][]float64, inputObjectives.ObjectiveCount)
	for i := 0; i < inputObjectives.ObjectiveCount; i++ {
		inputChromosome.Fitness[i] = make([]float64, len(inputChromosome.Subs))
	}

	// initialize current & aggregate fitness
	var aggFit float64 = 0.0
	var curFit float64 = 0.0

	// evaluate chromosome length and objectives to compute fitnesses
	for i := 0; i < inputObjectives.ObjectiveCount; i++ {
		for j := 0; j < chromLen; j++ {
			curFit = inputObjectives.Objectives[i].Matrix.At(inputChromosome.Subs[j][0], inputChromosome.Subs[j][1])
			inputChromosome.Fitness[i][j] = curFit
			inputChromosome.TotalFitness[i] = inputChromosome.TotalFitness[i] + curFit
		}

		// compute aggregate fitness
		aggFit = aggFit + inputChromosome.TotalFitness[i]
	}

	// calculate aggregate fitness
	inputChromosome.AggregateFitness = aggFit

	// return outputs
	return inputChromosome
}

// fitness function generate the mean fitness values for all of the chromosomes
// in a given population
func PopulationFitness(inputPopulation *Population, inputParameters *Parameters, inputObjectives *MultiObjective) (outputPopulation *Population) {

	// initialize output
	var cumFit float64 = 0.0
	var aggMeanFit float64 = 0.0

	// iterate over the different objectives and drain channel to compute fitness
	for i := 0; i < inputObjectives.ObjectiveCount; i++ {
		for j := 0; j < inputParameters.PopSize; j++ {

			// read current chromosome from channel
			curChrom := <-inputPopulation.Chromosomes

			// compute cumulative fitness
			cumFit = cumFit + curChrom.TotalFitness[i]

			// recieve from channel
			inputPopulation.Chromosomes <- curChrom
		}

		// compute mean from cumulative
		inputPopulation.MeanFitness[i] = cumFit / float64(inputParameters.PopSize)

		// compute aggreage mean fitness
		aggMeanFit = aggMeanFit + inputPopulation.MeanFitness[i]
	}

	// write aggregate mean fitness to output
	inputPopulation.AggregateMeanFitness = aggMeanFit

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
		if chrom1.AggregateFitness > chrom2.AggregateFitness {
			output = chrom1
		} else {
			output = chrom2
		}
	} else { // inverted
		if chrom1.AggregateFitness > chrom2.AggregateFitness {
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
		chrom1 := <-inputPopulation.Chromosomes
		chrom2 := <-inputPopulation.Chromosomes

		go func(chrom1, chrom2 *Chromosome) {
			// write selection to output channel
			output <- ChromosomeSelection(chrom1, chrom2, selProb)
		}(chrom1, chrom2)
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
func SelectionCrossover(inputSelection chan *Chromosome, inputParameters *Parameters, inputObjectives *MultiObjective, inputDomain *Domain) (crossover chan *Chromosome) {

	// initialize crossover channel
	output := make(chan *Chromosome, inputParameters.PopSize)

	// initialize crossover loop
	for i := 0; i < inputParameters.PopSize; i++ {
		for {
			// extract chromosomes
			chrom1 := <-inputSelection
			chrom2 := <-inputSelection

			// initialize empty index slices
			var chrom1Ind []int
			var chrom2Ind []int

			// initialize empty chromosome
			empChrom := NewEmptyChromosome(inputDomain, inputObjectives)

			// check for valid crossover point
			chrom1Ind, chrom2Ind = ChromosomeIntersection(chrom1.Subs, chrom2.Subs)

			// resample chromosomes if no intersection
			if len(chrom1Ind) > 2 {
				empChrom.Subs = ChromosomeCrossover(chrom1Ind, chrom2Ind, chrom1.Subs, chrom2.Subs)
				empChrom = ChromosomeFitness(empChrom, inputObjectives)
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
	nInd := NeighborhoodSubs(mutationLocus)

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

// function to generate a generic subDomain for an arbitrary set of node
// subscripts contained within a given input search domain
func SubDomain(sourceLocus, destinationLocus []int, inputDomain *mat64.Dense) (subDomain *Domain, subSourceLocus, subDestinationLocus []int) {

	// compute row index value ranges
	minRow := math.Min(float64(sourceLocus[0]), float64(destinationLocus[0]))
	maxRow := math.Max(float64(sourceLocus[0]), float64(destinationLocus[0]))

	// compute column index value ranges
	minCol := math.Min(float64(sourceLocus[1]), float64(destinationLocus[1]))
	maxCol := math.Max(float64(sourceLocus[1]), float64(destinationLocus[1]))

	// generate ranges
	rowRng := []int{int(minRow - 1.0), int(maxRow + 1.0)}
	colRng := []int{int(minCol - 1.0), int(maxCol + 1.0)}

	// extract raw domain values
	rowSpread := rowRng[1] - rowRng[0]
	colSpread := colRng[1] - colRng[0]

	// initialize subdomain values
	rawDomMat := mat64.DenseCopyOf(inputDomain.View(rowRng[0], colRng[0], rowSpread, colSpread))

	// overwrite matrix if singleton dimension
	if rowSpread == 2 {
		rawDomMat = mat64.DenseCopyOf(inputDomain.View(rowRng[0], colRng[0], rowSpread+1, colSpread))
	}
	if colSpread == 2 {
		rawDomMat = mat64.DenseCopyOf(inputDomain.View(rowRng[0], colRng[0], rowSpread, colSpread+1))
	}

	// get subdomain matrix dimensions
	rows, cols := rawDomMat.Dims()

	// mask edge values
	rawDomMat.SetRow(0, make([]float64, rows+cols))
	rawDomMat.SetRow(rows-1, make([]float64, rows+cols))
	rawDomMat.SetCol(0, make([]float64, rows+cols))
	rawDomMat.SetCol(cols-1, make([]float64, rows+cols))

	// generate sub domain structure
	subDom := NewDomain(rawDomMat)

	// compute sub source and sub destination subscript indexes
	orient := Orientation(sourceLocus, destinationLocus)

	// allocate output variables
	subSrc := make([]int, 2)
	subDst := make([]int, 2)

	// pivot output subscript values on orientation vector
	if orient[0] == -1 && orient[1] == -1 {
		subSrc[0] = int(rows - 2.0)
		subSrc[1] = int(cols - 2.0)
		subDst[0] = int(1.0)
		subDst[1] = int(1.0)
	} else if orient[0] == -1 && orient[1] == 0 {
		subSrc[0] = int(rows - 2.0)
		subSrc[1] = int(cols - 2.0)
		subDst[0] = int(1.0)
		subDst[1] = int(1.0)
	} else if orient[0] == -1 && orient[1] == 1 {
		subSrc[0] = int(rows - 2.0)
		subSrc[1] = int(1.0)
		subDst[0] = int(1.0)
		subDst[1] = int(cols - 2.0)
	} else if orient[0] == 0 && orient[1] == -1 {
		subSrc[0] = int(rows - 2.0)
		subSrc[1] = int(cols - 2.0)
		subDst[0] = int(1.0)
		subDst[1] = int(1.0)
	} else if orient[0] == 0 && orient[1] == 1 {
		subSrc[0] = int(1.0)
		subSrc[1] = int(1.0)
		subDst[0] = int(rows - 2.0)
		subDst[1] = int(cols - 2.0)
	} else if orient[0] == 1 && orient[1] == -1 {
		subSrc[0] = int(1.0)
		subSrc[1] = int(cols - 2.0)
		subDst[0] = int(rows - 2.0)
		subDst[1] = int(1.0)
	} else if orient[0] == 1 && orient[1] == 0 {
		subSrc[0] = int(1.0)
		subSrc[1] = int(cols - 2.0)
		subDst[0] = int(rows - 2.0)
		subDst[1] = int(1.0)
	} else if orient[0] == 1 && orient[1] == 1 {
		subSrc[0] = int(1.0)
		subSrc[1] = int(1.0)
		subDst[0] = int(rows - 2.0)
		subDst[1] = int(cols - 2.0)
	}

	// return output
	return subDom, subSrc, subDst
}

//// function to translate the subscript index values for a given slice of input
//// loci relative to a given offset vector
func TranslateWalkSubs(sourceSubs []int, inputWalkSubs [][]int) (outputWalkSubs [][]int) {

	// initialize output
	wlkLen := len(inputWalkSubs)
	outWlkSubs := make([][]int, wlkLen)
	outWlkSubs[0] = make([]int, 2)
	outWlkSubs[0] = sourceSubs

	// loop through and translate subscript values
	for i := 1; i < wlkLen; i++ {
		nSubs := make([]int, 2)
		nSubs[0] = outWlkSubs[i-1][0] + (inputWalkSubs[i][0] - inputWalkSubs[i-1][0])
		nSubs[1] = outWlkSubs[i-1][1] + (inputWalkSubs[i][1] - inputWalkSubs[i-1][1])
		outWlkSubs[i] = nSubs
	}

	// return output
	return outWlkSubs
}

// function to generate a mutation within a given chromosome at a specified
// number of mutation loci
func ChromosomeMutation(inputChromosome *Chromosome, inputDomain *Domain, inputParameters *Parameters, inputObjectives *MultiObjective) (outputChromosome *Chromosome) {

	// compute chromosome len.gth
	lenChrom := len(inputChromosome.Subs)

	// initialize output chromosome
	output := inputChromosome

	// initialize reference domain matrix
	refDomain := mat64.NewDense(inputDomain.Rows, inputDomain.Cols, nil)

	// clone input domain matrix
	refDomain.Clone(inputDomain.Matrix)

	// block out cells on current chromosome
	for k := 0; k < lenChrom; k++ {
		refDomain.Set(inputChromosome.Subs[k][0], inputChromosome.Subs[k][1], 0.0)
	}

	// enter unbounded mutation search loop
	for {
		// generate mutation loci
		prvLocus, mutLocus, nxtLocus, mutIndex := MutationLoci(inputChromosome)

		// first check if deletion is valid, else perform mutation
		if Distance(prvLocus, nxtLocus) < 1.5 {

			// perform simple deletion of mutation index
			output.Subs = append(output.Subs[:mutIndex], output.Subs[(mutIndex+1):]...)

			// loop over objective and remove fitness values
			for r := 0; r < inputObjectives.ObjectiveCount; r++ {
				output.Fitness[r] = append(output.Fitness[r][:mutIndex], output.Fitness[r][(mutIndex+1):]...)
			}
			break
		} else {

			// generate mutation subdomain
			subMat := MutationSubDomain(prvLocus, mutLocus, nxtLocus, refDomain)

			// generate sub source and sub destination
			subSource := make([]int, 2)
			subDestin := make([]int, 2)
			subSource[0] = prvLocus[0] - mutLocus[0] + 2
			subSource[1] = prvLocus[1] - mutLocus[1] + 2
			subDestin[0] = nxtLocus[0] - mutLocus[0] + 2
			subDestin[1] = nxtLocus[1] - mutLocus[1] + 2

			// generate subdomain from sub matrix and generate sub basis
			subDomain := NewDomain(subMat)
			subParams := NewParameters(subSource, subDestin, 1, 1, inputParameters.RndCoef)
			subBasis := NewBasis(subSource, subDestin, subDomain)

			// check validity of sub domain
			subDomainTest := ValidateMutationSubDomain(subSource, subDestin, subMat)

			// resample if subdomain is invalid
			if subDomainTest == false {
				continue
			} else {

				// generate directed walk based mutation
				subWlk, tabuTest := MutationWalk(subParams.SrcSubs, subParams.DstSubs, subDomain, subParams, subBasis)

				// if tabu test fails abort mutation and restart
				if tabuTest == false {
					subWlk = make([][]int, 1)
					continue
				} else {

					subLen := len(subWlk)
					subFit := make([][]float64, inputObjectives.ObjectiveCount)

					// translate subscripts and evaluate fitnesses
					for i := 0; i < inputObjectives.ObjectiveCount; i++ {

						// initialize subfit section
						subFit[i] = make([]float64, subLen)

						// translate subscripts and compute sub walk fitness
						for j := 0; j < subLen; j++ {
							if i == 0 {
								subWlk[j][0] = subWlk[j][0] - 2 + mutLocus[0]
								subWlk[j][1] = subWlk[j][1] - 2 + mutLocus[1]
							}

							subFit[i][j] = inputObjectives.Objectives[i].Matrix.At(subWlk[j][0], subWlk[j][1])
						}

						// delete mutation locus from fitnesses
						output.Fitness[i] = append(output.Fitness[i][:mutIndex], output.Fitness[i][(mutIndex+1):]...)

						// insert sub walk section into original chromosome fitnesses
						output.Fitness[i] = append(output.Fitness[i][:mutIndex-1], append(subFit[i], output.Fitness[i][mutIndex+1:]...)...)
					}

					// delete mutation locus from subs
					output.Subs = append(output.Subs[:mutIndex], output.Subs[(mutIndex+1):]...)

					// insert new sub walk subscripts into subs
					output.Subs = append(output.Subs[:mutIndex-1], append(subWlk, output.Subs[mutIndex+1:]...)...)
					break
				}
			}
		}

	}

	// return output
	return output
}

// function to generate multiple mutations on multiple separate loci on the same
// input chromosome
func ChromosomeMultiMutation(inputChromosome *Chromosome, inputDomain *Domain, inputParameters *Parameters, inputObjectives *MultiObjective) (outputChromosome *Chromosome) {

	// loop through mutation count
	for i := 0; i < inputParameters.MutaCnt; i++ {
		inputChromosome = ChromosomeMutation(inputChromosome, inputDomain, inputParameters, inputObjectives)
	}

	// return output
	return inputChromosome
}

// function to generate mutations within a specified fraction of an input
// population with those chromosomes being selected at random
func PopulationMutation(inputChromosomes chan *Chromosome, inputParameters *Parameters, inputObjectives *MultiObjective, inputDomain *Domain) (outputChromosomes chan *Chromosome) {

	// calculate the total number of chromosomes that are to receive mutations
	mutations := int(math.Floor(float64(inputParameters.PopSize) * float64(inputParameters.MutaFrc)))

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// initialize mutation selection test variable and iteration counter variable
	var iter int
	var mutTest int

	// initialize throttle size
	conc := make(chan bool, inputParameters.ConSize)

	// initialize selection loop
	for j := 0; j < inputParameters.PopSize; j++ {

		// get current chromosome from channel
		curChrom := <-inputChromosomes

		// generate random mutation selection binary integer
		mutTest = rand.Intn(2)

		// screen on mutation indices
		if mutTest == 1 {

			// write to control channel
			conc <- true

			// launch go routines
			go func(curChrom *Chromosome, inputDomain *Domain, inputParameters *Parameters, inputObjectives *MultiObjective) {
				defer func() { <-conc }()
				curChrom = ChromosomeMultiMutation(curChrom, inputDomain, inputParameters, inputObjectives)
			}(curChrom, inputDomain, inputParameters, inputObjectives)

			// update iterator
			iter += 1

			// return current chromosome back to channel
			inputChromosomes <- curChrom

		} else {

			// return current chromosome back to channel
			inputChromosomes <- curChrom
		}

		// break once the desired number of mutants has been generated
		if iter == mutations {
			break
		}
	}

	// cap parallelism at concurrency limit
	for j := 0; j < cap(conc); j++ {
		conc <- true
	}

	// return selection channel
	return inputChromosomes
}

// population evolution operator generates a new population
// from an input population using the selection and crossover operators
func PopulationEvolution(inputPopulation *Population, inputDomain *Domain, inputParameters *Parameters, inputObjectives *MultiObjective) (outputPopulation *Population) {

	// initialize new empty population
	output := NewEmptyPopulation(inputPopulation.Id+1, inputObjectives)

	// perform population selection
	popSel := PopulationSelection(inputPopulation, inputParameters)

	// perform selection crossover
	selCrs := SelectionCrossover(popSel, inputParameters, inputObjectives, inputDomain)

	// fill empty population
	popMut := PopulationMutation(selCrs, inputParameters, inputObjectives, inputDomain)

	// assign channel to output population
	output.Chromosomes = popMut

	// return output
	return output
}
