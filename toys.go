// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

// new test domain initialization function
func NewToyDomain(identifier, domainSize, domainStride int) *Domain {

	// initialize value slice
	domainValues := make([]bool, domainSize)

	// loop through index values togo define domain
	for i := 0; i < domainSize; i++ {
		if i >= 0 && i < domainStride {
			domainValues[i] = false
		} else if i > (domainSize - domainStride) {
			domainValues[i] = false
		} else {
			domainValues[i] = true
		}
	}
	for i := 0; i < domainSize; i = (i + domainStride) {
		domainValues[i] = false
	}
	for i := 1; i < domainSize; i = (i + domainStride) {
		domainValues[i] = false
	}

	// return output
	return &Domain{
		id:     identifier,
		size:   domainSize,
		stride: domainStride,
		vals:   domainValues,
	}
}

// func NewToyObjective(identifier Int, problemParameters Parameters, problemDomain, )
