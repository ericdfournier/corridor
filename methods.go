/* Copyright ©2015 The corridor Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file. */

package corridor

// start walker method
func (w Walker) Start(searchDomain *Domain, searchParameters *Parameters, searchObjectives *MultiObjective, chr chan *Chromosome, walkQueue chan bool) {

	// Initialize goroutine
	go func(searchDomain *Domain, searchParameters *Parameters, searchObjectives *MultiObjective, chr chan *Chromosome, walkQueue chan bool) {

		// enter unbounded for/select loop
		for {
			// pull walk counter
			walk := <-walkQueue

			if walk == true {

				// initialize new empty chromosome
				newChrom := NewEmptyChromosome(searchDomain, searchObjectives)

				// start walk to generate new chromosome
				newChrom = NewChromosome(searchDomain, searchParameters, searchObjectives)

				// compute chromosome fitness and return to channel
				chr <- ChromosomeFitness(newChrom, searchObjectives)

			} else if walk == false {

				close(walkQueue)

				return
			}
		}
	}(searchDomain, searchParameters, searchObjectives, chr, walkQueue)
}
