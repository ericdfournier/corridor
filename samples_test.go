// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import "testing"

// small convex problem benchmark
func BenchmarkConvexSmall(b *testing.B) {

	// initialize domain
	sampleDomain := NewSampleDomain(20, 20)

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

// medium convex problem benchmark
func BenchmarkConvexMedium(b *testing.B) {

	// initialize domain
	sampleDomain := NewSampleDomain(20, 20)

	// initialize objectives
	objectiveCount := 3
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)
	sampleParameters.PopSize = 10000

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)
}

// large convex problem benchmark
func BenchmarkConvexLarge(b *testing.B) {

	// initialize domain
	sampleDomain := NewSampleDomain(20, 20)

	// initialize objectives
	objectiveCount := 3
	sampleObjectives := NewSampleObjectives(sampleDomain.Rows, sampleDomain.Cols, objectiveCount)

	// initialize parameters
	sampleParameters := NewSampleParameters(sampleDomain)
	sampleParameters.PopSize = 100000

	// evolve populations
	toyEvolution := NewEvolution(sampleParameters, sampleDomain, sampleObjectives)

	// extract output population
	finalPop := <-toyEvolution.Populations

	// view output population
	ViewPopulation(sampleDomain, sampleParameters, finalPop)
}
