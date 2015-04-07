// Copyright ©2015 The corridor Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package corridor

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/gonum/matrix/mat64"
)

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
