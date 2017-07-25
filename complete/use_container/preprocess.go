package main

import (
	"fmt"
	"io/ioutil"

	"github.com/jdkato/prose/tokenize"
)

// ------------- TODO: not sure where to include these preprocessing functions,
// Ideally, they would sit outside of main.
// ------------ Helper function related to loading data into dd_container fmt

// TODO: this functionality is dd_container specific.

// ------------------------------------------- Cosine

// createMap is a helper that accepts a path to a directory and creates the
// input data for the model.
// NOTE: this may/may not be included in functionality.  It will likely fall on
// the user to create the specified input data.
func createMap(dPath string) (map[string][]string, error) {
	fileMap := make(map[string][]string)

	// Create map of filename to tokenized content.
	dFiles, _ := ioutil.ReadDir(dPath)
	for _, f := range dFiles {
		b, err := ioutil.ReadFile(dPath + f.Name())
		if err != nil {
			fmt.Println(err)
		}

		// Convert bytes to string, then use 3rd party to tokenize.
		fileMap[f.Name()] = tokenize.TextToWords(string(b))
	}

	return fileMap, nil
}

// ------------------------------------------------[END] cosine

// ------------------------------------------------ Sudoku

// // crossIndex 'crosses' two strings such that the two individual values from
// // each string join together to create a new value.  For example, if string one
// // is "ABC" and string two is "123", the resulting return value will be;
// // ["A1","A2","A3","B1","B2","B3","C1","C2","C3"].
// func crossIndex(A string, N string) []string {
// 	var ks []string
// 	for _, a := range A {
// 		for _, n := range N {
// 			ks = append(ks, (string(a) + string(n)))
// 		}
// 	}
// 	return ks
// }

// func createSudokuInput(fPath string) (map[string][]string, error) {
// 	sudokuMap := make(map[string][]string)

// 	data, err := ioutil.ReadFile(fPath)
// 	if err != nil {
// 		return sudokuMap, err
// 	}

// 	// Global board information.  The Sudoku board is assumed to be a standard
// 	// 9x9 (A-I)x(1-9) grid -- where the first index (upper left) would be `A1`
// 	// and the last index (lower right) would be `I9`.
// 	rows := "ABCDEFGHI"
// 	cols := "123456789"
// 	inds := crossIndex(rows, cols)

// 	// Convert the string representing the board into a grid(map) that maps a
// 	// key (index) to the values (label for the box, or possible label for the
// 	// box). for instance, if we know A1=7, map['A1'] = '7', but if the given
// 	// index is empty (B2, as an example), the corresponding value would be
// 	// '123456789' (map['B2'] = '123456789') NOTE: i acts as an increment for
// 	// every target character found.
// 	i := 0
// 	for _, c := range data {
// 		switch string(c) {
// 		case "_":
// 			sudokuMap[inds[i]] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
// 			i++
// 		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
// 			sudokuMap[inds[i]] = []string{string(c)}
// 			i++
// 		case "\n", " ", "\r":
// 			continue
// 		default:
// 			return sudokuMap, fmt.Errorf("unexpected value (%v) in Sudoku input", c)
// 		}
// 	}

// 	return sudokuMap, nil
// }
// ------------------------------------------------ [END]Sudoku
