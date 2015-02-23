// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import "github.com/gonum/matrix/mat64"

func Fitness(subs [][]int, obj *mat64.Dense) (fitnessValues []float64, totalFitness float64) {

	// get individual length
	indSize := len(subs)

	// initialize fitness values and total fitness
	fitVal := make([]float64, indSize)
	var totFit float64
	totFit = 0.0

	// evaluate individual fitness according to input objective
	for i := 0; i < indSize; i++ {
		curFit := obj.At(subs[i][0], subs[i][1])
		fitVal[i] = curFit
		totFit = totFit + curFit
	}

	// return outputs
	return fitVal, totFit

}
