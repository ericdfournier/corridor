// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import (
	"math"
	"math/rand"
	"time"
)

// tournament selection operator selects half of the input
// population for reproduction based upon comparative fitness
// and some randomized input selection fraction
func TournamentSelection(inputPopulation *Population, inputParameters *Parameters) (selection chan *Chromosome) {

	// initialize selection channel size
	selSize := int(math.Floor(float64(cap(inputPopulation.Chromosomes)) * inputParameters.SelFrac))

	// initialize selection channel
	output := make(chan *Chromosome, selSize)

	// get current time for random number seed
	rand.Seed(time.Now().UnixNano())

	// generate random number to determine selection result
	dec := rand.Float64()

	// initialize selection channel
	for i := 0; i < selSize; i++ {
		chrom1 := <-inputPopulation.Chromosomes
		chrom2 := <-inputPopulation.Chromosomes

		if dec > inputParameters.SelProb { // normal
			if chrom1.TotalFitness > chrom2.TotalFitness {
				output <- chrom1
			} else {
				output <- chrom2
			}
		} else { // inverted
			if chrom1.TotalFitness > chrom2.TotalFitness {
				output <- chrom2
			} else {
				output <- chrom1
			}
		}
	}

	// return selection channel
	return output
}

//func Crossover(inputSelection chan *Chromosome) {

//}
