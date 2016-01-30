/* Copyright ©2015 The corridor Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file. */

package corridor

import (
	"math/rand"
	"time"
)

// walker method to initialize a parallel pseudo random walk
func (w Walker) Start(chr chan *Chromosome, walkQueue chan bool) {

	// initialize goroutine
	go func() {

		// enter unbounded for/select loop
		for {
			// pull walk token
			walk := <-walkQueue

			if walk == true {

				// initialize new empty chromosome
				newChrom := NewEmptyChromosome(w.SearchDomain, w.SearchObjectives)

				// start walk to generate new chromosome
				newChrom = NewChromosome(w.SearchDomain, w.SearchParameters, w.SearchObjectives)

				// compute chromosome fitness and return to channel
				chr <- ChromosomeFitness(newChrom, w.SearchObjectives)

			} else if walk == false {

				close(walkQueue)

				return
			}
		}
	}()
}

// mutator method to initialize a parallel mutation procedure
func (m Mutator) Start(chr chan *Chromosome, mutationQueue chan bool) {

	// initialize goroutine
	go func() {

		// initialize ubounded for loop
		for {

			// seed random number generator
			rand.Seed(time.Now().UnixNano())

			// generate random mutation selection binary integer
			mutTest := rand.Intn(2)

			// initialize token
			mutate := <-mutationQueue

			// mutation selected and token available
			if mutTest == 1 && mutate == true {

				// pull chromosomes and generate mutation
				curChrom := <-chr
				curChrom = ChromosomeMultiMutation(curChrom, m.SearchDomain, m.SearchParameters, m.SearchObjectives)
				chr <- curChrom

				// mutation not selected and token available
			} else if mutTest != 1 && mutate == true {

				// return token to channel and exit
				mutationQueue <- mutate

				// no tokens available
			} else if mutate == false {

				close(mutationQueue)
				return
			}
		}
	}()
}
