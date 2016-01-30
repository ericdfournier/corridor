/* Copyright ©2015 The corridor Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file. */

package corridor

// walker method to initialize a parallel pseudo random walk
func (w Walker) Start(chr chan *Chromosome, walkQueue chan bool) {

	// Initialize goroutine
	go func(chr chan *Chromosome, walkQueue chan bool) {

		// enter unbounded for/select loop
		for {
			// pull walk counter
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
	}(chr, walkQueue)
}
