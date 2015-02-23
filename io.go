// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package corridor

import "fmt"

func ViewSearchDomain(searchDomain *Domain) {

	// get search domain matrix size
	rows, _ := searchDomain.Matrix.Dims()

	// print domain id
	fmt.Printf("Search Domain ID = %d\n", searchDomain.Id)

	// print domain values to command line
	fmt.Printf("Search Domain Values = \n")
	for i := 0; i < rows; i++ {
		rawRowVals := searchDomain.Matrix.RawRowView(i)
		fmt.Println(rawRowVals)
	}
}
