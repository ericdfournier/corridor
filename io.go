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
func CsvToSubs(inputFilepath string) (outputSubs []int) {

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

	// initialize output
	output := make([]int, 2)

	// loop through and extract values
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {

			// get string value and convert to integer
			strVal := rawCSVdata[i][j]
			intVal, err := strconv.Atoi(strVal)

			// shift value by one to account for buffer boundaries
			output[j] = intVal + 1

			// parse error
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	// return output
	return output
}

// function to write an input comma separated value
// file's contents to an output domain structure
func CsvToDomain(inputFilepath string) (outputDomain *Domain) {

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
	output := NewDomain(domMat)

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

// function to write a set of input comma separated value
// files' contents to an output multiobjective structure
func CsvToMultiObjective(inputFilepaths ...string) (outputMultiObjective *MultiObjective) {

	// get variadic input length
	objectiveCount := len(inputFilepaths)

	// initialize objective slice
	objectiveSlice := make([]*Objective, objectiveCount)

	// initialize objectives identifier
	var objectiveID int = 0

	// loop through and extract objectives
	for i := 0; i < objectiveCount; i++ {

		// read CSV data to objective
		objectiveSlice[i] = CsvToObjective(objectiveID, inputFilepaths[i])

		// increment objective identifier
		objectiveID += 1

	}

	// return multiObjective output
	return &MultiObjective{
		ObjectiveCount: objectiveCount,
		Objectives:     objectiveSlice,
	}
}

// function to write the values from an input
// chromosome structure to an output csv file
func ChromosomeToString(inputChromosome *Chromosome) (outputRawString [][]string) {

	// get input chromosome length
	chromLen := len(inputChromosome.Subs)

	// count input chromosome objectives
	objCount := len(inputChromosome.TotalFitness)

	// intitialize raw output string slice
	rawCSVdata := make([][]string, objCount+2)

	// loop through and format values as strings for output encoding
	for j := 0; j < objCount+2; j++ {

		// allocate inner slice
		rawCSVdata[j] = make([]string, chromLen)

		for i := 0; i < chromLen; i++ {

			// transpose subs by one to account for boundary buffer
			if j == 0 {
				rawCSVdata[j][i] = strconv.Itoa(inputChromosome.Subs[i][0]) - 1
			} else if j == 1 {
				rawCSVdata[j][i] = strconv.Itoa(inputChromosome.Subs[i][1]) - 1
			} else {
				rawCSVdata[j][i] = strconv.FormatFloat(inputChromosome.Fitness[j-2][i], 'f', 2, 64)
			}
		}
	}

	// return output
	return rawCSVdata
}

// function to write the values from an input elite set
// to an output csv file
func EliteSetToCsv(inputEliteSet []*Chromosome, outputFilepath string) {

	// open file
	csvfile, err := os.Create(outputFilepath)

	// parse file opening errors
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// close file on completion
	defer csvfile.Close()

	// get chromosome count
	chromCount := len(inputEliteSet)

	// initialize rawCSVdata and chromosome string structures
	var chromString, rawCSVdata [][]string

	// loop through chromsomes and generate composite string structure
	for i := 0; i < chromCount; i++ {
		chromString = ChromosomeToString(inputEliteSet[i])
		rawCSVdata = append(rawCSVdata, chromString...)
	}

	// initialize writer object
	writer := csv.NewWriter(csvfile)

	// write data or get error
	err = writer.WriteAll(rawCSVdata)

	// parse errors
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// flush writer object
	writer.Flush()
}
