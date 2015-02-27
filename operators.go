// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// NEED TO CHECK IN SELECTION ROUTINE THAT THE INITIAL
// ALLOCATION FOR THE OUTPUT DOES, IN SOME CASES, GET
// OVERWRITTEN...OTHERWISE CHROM1 WILL ALWAYS BE SELECTED

// selection operator selects between two chromosomes with a
// probability of the most fit chromosome being selected
// determined by the input selection probability ratio
func Selection(chrom1, chrom2 *Chromosome, selectionProb float64) (selectedChrom *Chromosome) {

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

		// drain chromosome input channel
		chrom1 := <-inputPopulation.Chromosomes
		chrom2 := <-inputPopulation.Chromosomes

		// launch selection go routines
		go func(chrom1, chrom2 *Chromosome, selProb float64) {
			output <- Selection(chrom1, chrom2, selProb)
		}(chrom1, chrom2, selProb)
	}

	// return selection channel
	return output
}

// intersection determines whether or not the subscripts
// associated with two input chromosomes share any in
// values in common and reports their relative indices
func Intersection(subs1, subs2 [][]int) (subs1Indices, subs2Indices []int) {

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
func Crossover(chrom1Ind, chrom2Ind []int, chrom1Subs, chrom2Subs [][]int) (crossoverChrom [][]int) {

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

	fmt.Println(len(chrom1Ind))
	fmt.Println(r)

	for i := 0; i < (chrom1Ind[r] + 1); i++ {
		output = append(output, chrom1Subs[i])
	}

	for j := (chrom2Ind[r] + 1); j < len(chrom2Subs); j++ {
		output = append(output, chrom2Subs[j])
	}

	// return output
	return output

}

//// selection crossover operator performs a single part
//// crossover on each of the individuals provided in an
//// input selection channel of chromosomes
//func SelectionCrossover(inputSelection chan *Chromosome, inputParameters *Parameters) (crossover chan *Chromosome) {

//	// initialize crossover channel size
//	croSize := int(inputParameters.PopSize)

//	// initialize crossover channel
//	output := make(chan *Chromosome, croSize)

//	// initialize crossover loop
//	for i := 0; i < croSize; i++ {

//		for {

//			// extract chromosomes
//			chrom1 := <-inputSelection
//			chrom2 := <-inputSelection

//			// check for valid crossover point
//			chrom1Ind, chrom2Ind := Intersection(chrom1, chrom2)

//			// resample chromosomes if no intersection
//			if len(chrom1Ind) > 0 {
//				break
//			}
//			else {
//				inputSelection <- chrom1
//				inputSelction <- chrom2
//				continue
//			}
//		}

//		// launch crossover routines
//		go func(chrom1Ind, chrom2Ind []int, chrom1, chrom2 *Chromosome) {
//			output <- Crossover(chrom1Ind, chrom2Ind, chrom1, chrom2)
//		}(chrom1Ind, chrom2Ind, chrom1, chrom2)

//		// repopulate selection channel
//		inputSelection <- chrom1
//		inputSelection <- chrom2
//	}

//	// return output
//	return output
//}
