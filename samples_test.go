// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import "testing"

func BenchmarkConvexSmall(b *testing.B) {

	// initialize domain
	sampleDomain := NewSampleDomain(50, 50)

	// initialize objectives
	objectiveCount := 3
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)

}
