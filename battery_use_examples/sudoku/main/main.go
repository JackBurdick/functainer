package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/JackBurdick/dataduit/dataduit"
)

func main() {

	// Initialize cosine battery.
	var sudokuContainer dataduit.DdContainer

	// Configure Sudoku container.
	err := sudokuContainer.Configure("./sudoku_config.yml")
	if err != nil {
		fmt.Printf("Error with sudokuContainer config: %v\n", err)
	}

	// Preprocess input (according to required specified input type)
	sjsonData, err := createSudokuInput("../input/puzzle_01.txt")
	if err != nil {
		fmt.Println(err)
	}

	// Use sudoku container.
	sudokuRes, err := sudokuContainer.FullUse(sjsonData)
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}

	// Display result
	fmt.Printf("Sudoku Result: %v\n", sudokuRes)

}

// ---------------------------------------------- INPUT PREPROCESSING FUNCTIONS
// NOTE: This will likely fall on the user to create the specified input data.
// and will only be included as an example that the user can pull from.

// Sudoku
// crossIndex 'crosses' two strings such that the two individual values from
// each string join together to create a new value.  For example, if string one
// is "ABC" and string two is "123", the resulting return value will be;
// ["A1","A2","A3","B1","B2","B3","C1","C2","C3"].
func crossIndex(A string, N string) []string {
	var ks []string
	for _, a := range A {
		for _, n := range N {
			ks = append(ks, (string(a) + string(n)))
		}
	}
	return ks
}

func createSudokuInput(fPath string) ([]byte, error) {
	sudokuMap := make(map[string][]string)

	data, err := ioutil.ReadFile(fPath)
	if err != nil {
		// Build json data from map.
		sjsonData, err := json.Marshal(sudokuMap)
		if err != nil {
			fmt.Println(err)
		}
		return sjsonData, err
	}

	// Global board information.  The Sudoku board is assumed to be a standard
	// 9x9 (A-I)x(1-9) grid -- where the first index (upper left) would be `A1`
	// and the last index (lower right) would be `I9`.
	rows := "ABCDEFGHI"
	cols := "123456789"
	inds := crossIndex(rows, cols)

	// Convert the string representing the board into a grid(map) that maps a
	// key (index) to the values (label for the box, or possible label for the
	// box). for instance, if we know A1=7, map['A1'] = '7', but if the given
	// index is empty (B2, as an example), the corresponding value would be
	// '123456789' (map['B2'] = '123456789') NOTE: i acts as an increment for
	// every target character found.
	i := 0
	for _, c := range data {
		switch string(c) {
		case "_":
			sudokuMap[inds[i]] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
			i++
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			sudokuMap[inds[i]] = []string{string(c)}
			i++
		case "\n", " ", "\r":
			continue
		default:
			// Build json data from map.
			sjsonData, err := json.Marshal(sudokuMap)
			if err != nil {
				fmt.Println(err)
			}
			return sjsonData, fmt.Errorf("unexpected value (%v) in Sudoku input", c)
		}
	}

	// Build json data from map.
	sjsonData, err := json.Marshal(sudokuMap)
	if err != nil {
		fmt.Println(err)
	}

	return sjsonData, nil
}

//------------------------------------------------ Sudoku [END]
// ----------------------------------------- INPUT PREPROCESSING FUNCTIONS [END]
