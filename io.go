// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strconv"

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
				domMat.Set(i, j, 1.0)
			}
		}
	}

	// initialize new domain
	output := NewDomain(identifier, domMat)

	// return output
	return output
}

// function to import an input RGBA formatted jpeg and
// convert it into an output domain structure
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

// function to import an input RGBA formatted geotiff and
// convert it into an output domain structure
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

// function to write an input comma separated value
// file's contents to an output domain structure
func CsvToDomain(identifier int, inputFilepath string) (outputDomain *Domain) {

	// open file
	data, err := os.Open(inputFilepath)

	// parse error if file not found
	if err != nil {
		fmt.Println(err)
		return
	}

	// close file on completion
	defer data.Close()

	// generate new reader from open file
	reader := csv.NewReader(data)

	// set reader structure field
	reader.FieldsPerRecord = -1

	// use reader to read raw csv data
	rawCSVdata, err := reader.ReadAll()

	// parse csv file formatting errors
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// initialize empty row and column counts
	rows := len(rawCSVdata)
	cols := len(rawCSVdata[0])

	// initialize domain matrix
	domMat := mat64.NewDense(rows+2, cols+2, nil)

	// write values from rawCSVdata to domain matrix
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

				// get string value and convert to integer
				strVal := rawCSVdata[i-1][j-1]
				fltVal, err := strconv.ParseFloat(strVal, 64)

				// parse error
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// write value to matrix
				domMat.Set(i, j, fltVal)
			}
		}
	}

	// initialize new domain
	output := NewDomain(identifier, domMat)

	// return output
	return output
}

// function to convert an input image structure
// into and output objective structure
func ImageToObjective(identifier int, inputImage image.Image) (outputObjective *Objective) {

	// get input image boundaries
	rows := inputImage.Bounds().Max.X - 1
	cols := inputImage.Bounds().Max.Y - 1

	// initialize domain matrix
	objMat := mat64.NewDense(rows+2, cols+2, nil)

	// write values from image object to domain matrix
	for i := 0; i < rows+2; i++ {
		for j := 0; j < cols+2; j++ {

			// create a 1 pixel boundary buffer of zeros
			if i == 0 {
				objMat.Set(i, j, 0.0)
			} else if i == rows+1 {
				objMat.Set(i, j, 0.0)
			} else if j == 0 {
				objMat.Set(i, j, 0.0)
			} else if j == cols+1 {
				objMat.Set(i, j, 0.0)
			} else {

				//// NEED TO IMPOSE THAT THE INPUT JPEG IS FORMATTED AS A GRAYSCALE
				//r, g, b, a := inputImage.At(i-1, j-1).RGBA()
				//objMat.Set(i, j, r)

			}
		}
	}

	// initialize new domain
	output := NewObjective(identifier, objMat)

	// return output
	return output
}

// function to import an input RGBA formatted jpeg and
// convert it into an output objective structure
func JpegToObjective(identifier int, inputFilepath string) (outputObjective *Objective) {

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
	output := ImageToObjective(identifier, img)

	// return output
	return output
}

// function to import an input RGBA formatted tiff and
// convert it into an output objective structure
func TiffToObjective(identifier int, inputFilepath string) (outputObjective *Objective) {

	// read local jpeg image file
	data, err := ioutil.ReadFile(inputFilepath)

	// parse file not found error
	if err != nil {
		fmt.Println("File not found!")
		os.Exit(1)
	}

	// decode input jpeg image
	img, err := tiff.Decode(bytes.NewReader(data))

	// parse error
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// convert image to domain
	output := ImageToObjective(identifier, img)

	// return output
	return output
}

// function to write an input comma separated value
// file's contents to an output objective structure
func CsvToObjective(identifier int, inputFilepath string) (outputObjective *Objective) {

	// open file
	dataFile, err := os.Open(inputFilepath)

	// parse error if file not found
	if err != nil {
		fmt.Println(err)
		return
	}

	// close file on completion
	defer dataFile.Close()

	// generate new reader from open file
	reader := csv.NewReader(dataFile)

	// set reader structure field
	reader.FieldsPerRecord = -1

	// use reader to read raw csv data
	rawCSVdata, err := reader.ReadAll()

	// parse csv file formatting errors
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// initialize empty row and column counts
	rows := len(rawCSVdata)
	cols := len(rawCSVdata[0])

	// initialize domain matrix
	objMat := mat64.NewDense(rows+2, cols+2, nil)

	// write values from rawCSVdata to domain matrix
	for i := 0; i < rows+2; i++ {
		for j := 0; j < cols+2; j++ {

			// create a 1 pixel boundary buffer of zeros
			if i == 0 {
				objMat.Set(i, j, 0.0)
			} else if i == rows+1 {
				objMat.Set(i, j, 0.0)
			} else if j == 0 {
				objMat.Set(i, j, 0.0)
			} else if j == cols+1 {
				objMat.Set(i, j, 0.0)
			} else {

				// get string value and convert to float
				strVal := rawCSVdata[i-1][j-1]
				fltVal, err := strconv.ParseFloat(strVal, 64)

				// parse error
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// write matrix value
				objMat.Set(i, j, fltVal)
			}
		}
	}

	// initialize new domain
	output := NewObjective(identifier, objMat)

	// return output
	return output
}

//// function to write an the values from an input
//// chromosome structure to an output csv file
//func ChromosomeToCsv(inputChromosome *Chromosome, outputFilepath string) {

//	// open file
//	csvfile, err := os.Create(outputFilepath)

//	// parse file opening errors
//	if err != nil {
//    	fmt.Println("Error:", err)
//        return
//	}

//	// close file on completion
//	defer csvfile.Close()

//	// get input chromosome length
//	chromLen := len(inputChromosome.Subs)

//	// count input chromosome objectives
//	objCount := len(inputChromosome.Fitness)

//	// intitialize raw output string slice
//	rawCSVdata := make([][]string, 2+objCount)
//	rawCSVdata[0] := make([]string, chromLen)

//	// extract and encode the subs
//	for i := 0; i < chromLen; i++ {
//			rawCSVdata[0][i] = strconv.Itoa(inputChromosome.Subs[i][0])
//			rawCSVdata[1][i] = strconv.Itoa(inputChromosome.Subs[i][1])
//			for j := 0; j < objCount; j++ {
//				rawCSVdata[j+2][i] = strconv.Itoa(inputChromosome.Fitness[j][i])
//			}
//	}

//	//records := [][]string{{"item1", "value1"}, {"item2", "value2"}, {"item3", "value3"}}

//  //        writer := csv.NewWriter(csvfile)
//  //        for _, record := range records {
//  //                err := writer.Write(record)
//  //                if err != nil {
//  //                        fmt.Println("Error:", err)
//  //                        return
//  //                }
//  //        }
//  //        writer.Flush()

//}
