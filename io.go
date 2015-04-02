// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/chai2010/tiff"
	"github.com/gonum/matrix/mat64"
)

// function to convert and image structure to a domain structure
func ImageToDomain(identifier int, inputImage image.Image) (outputDomain *Domain) {

	// get input image boundaries
	rows := inputImage.Bounds().Max.X - 1
	cols := inputImage.Bounds().Max.Y - 1

	// initialize domain matrix
	domMat := mat64.NewDense(rows+2, cols+2, nil)

	// write values from image object to domain matrix
	for i := 0; i < rows+2; i++ {
		for j := 0; j < cols+2; j++ {

			// create a 1 pixel boundary buffer of zeros
			if i == 0 {
				domMat.Set(i, j, 0.0)
			} else if i == rows+1 {
				domMat.Set(i, j, 0.0)
			} else if j == 0 {
				domMat.Set(i, j, 0.0)
			} else if j == cols+1 {
				domMat.Set(i, j, 0.0)
			} else {
				r, g, b, _ := inputImage.At(i-1, j-1).RGBA()

				if r > 0 || g > 0 || b > 0 {
					domMat.Set(i, j, 1.0)
				} else {
					domMat.Set(i, j, 0.0)
				}
			}
		}
	}

	// initialize new domain
	output := NewDomain(identifier, domMat)

	// return output
	return output
}

// function to convert an input RGBA formatted jpeg into
// an output domain structure
func JpegToDomain(identifier int, inputFilepath string) (outputDomain *Domain) {

	// read local jpeg image file
	data, err := ioutil.ReadFile(inputFilepath)

	// parse file not found error
	if err != nil {
		fmt.Println("File not found!")
		os.Exit(1)
	}

	// decode input jpeg image
	img, err := jpeg.Decode(bytes.NewReader(data))

	// parse error
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// convert image to domain
	output := ImageToDomain(identifier, img)

	// return output
	return output
}

// function to convert an input RGBA formatted geotiff into
// an output domain structure
func TiffToDomain(identifier int, inputFilepath string) (outputDomain *Domain) {

	// read byte data
	data, err := ioutil.ReadFile(inputFilepath)

	// parse file not found error
	if err != nil {
		fmt.Println("File not found!")
		os.Exit(1)
	}

	// decode input diff image
	img, err := tiff.Decode(bytes.NewReader(data))

	// parse error
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// convert image to domain
	output := ImageToDomain(identifier, img)

	// return output
	return output
}
