/* Copyright ©2015 The corridor Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file. */

package corridor

import (
	"math/rand"
	"sync"
	"time"
)

// walker method to initialize a parallel pseudo random walk
func (w Walker) Start(chr chan *Chromosome, walkQueue chan bool, wg sync.WaitGroup) {

	// add go routine to waitgroup
	wg.Add(1)

	// initialize goroutine
	go func() {

		// close on completion
		defer wg.Done()

		// enter unbounded for/select loop
		for {

			// select on walk queue token availability
			select {

			// tokens available
			case <-walkQueue:

				// initialize new empty chromosome
				newChrom := NewEmptyChromosome(w.SearchDomain, w.SearchObjectives)

				// start walk to generate new chromosome
				newChrom = NewChromosome(w.SearchDomain, w.SearchParameters, w.SearchObjectives)

				// compute chromosome fitness and return to channel
				chr <- ChromosomeFitness(newChrom, w.SearchObjectives)

			// tokens not available
			default:

				// terminate go routine
				return
			}
		}
	}()
}

// mutator method to initialize a parallel mutation procedure
func (m Mutator) Start(chr chan *Chromosome, mutationQueue chan bool, wg sync.WaitGroup) {

	// add go routine to waitgroup
	wg.Add(1)

	// initialize goroutine
	go func() {

		// close on completion
		defer wg.Done()

		// enter ubounded for/select loop
		for {

			// seed random number generator
			rand.Seed(time.Now().UnixNano())

			// generate random mutation selection binary integer
			mutTest := rand.Intn(2)

			// pull chromosomes
			curChrom := <-chr

			// select on mutation token availability
			select {

			// tokens available
			case mutate := <-mutationQueue:

				// muation desired
				if mutTest == 1 {

					// mutate current chromosome
					curChrom = ChromosomeMultiMutation(curChrom, m.SearchDomain, m.SearchParameters, m.SearchObjectives)

					// return mutant to channel
					chr <- curChrom

					// mutation not desired
				} else if mutTest != 1 {

					// return token to channel
					mutationQueue <- mutate

					// return chromosome to channel
					chr <- curChrom
				}

			// no tokens available
			default:

				// return chromosome to channel
				chr <- curChrom

				// terminate go routine
				return
			}
		}
	}()
}
