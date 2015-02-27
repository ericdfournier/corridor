// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"math"
	"math/rand"
	"time"
)

func Selection(chrom1, chrom2 *Chromosome, selectionProb float64) (selChrom *Chromosome) {

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

// tournament selection operator selects half of the input
// population for reproduction based upon comparative fitness
// and some randomized input selection fraction
func TournamentSelection(inputPopulation *Population, inputParameters *Parameters) (selection chan *Chromosome) {

	// initialize selection channel size
	selSize := int(math.Floor(float64(cap(inputPopulation.Chromosomes)) * inputParameters.SelFrac))

	// initialize selection channel
	output := make(chan *Chromosome, selSize)

	// initialize selection probability
	selProb := inputParameters.SelProb

	// initialize selection channel
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

//func Crossover(inputSelection chan *Chromosome) {

//}
