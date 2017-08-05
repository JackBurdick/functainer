package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/JackBurdick/dataduit/dataduit"
	"github.com/jdkato/prose/tokenize"
)

func main() {
	var err error

	// Initialize cosine battery.
	var cosineContainer dataduit.DdContainer

	// Build battery according to specification.
	err = cosineContainer.Configure("./config/cosine_config.yml")
	if err != nil {
		fmt.Printf("Error with cosineContainer config: %v\n", err)
	}

	// Use cosine battery.
	cInputPath := "./input/input_cosine/"
	fileMap, err := createMap(cInputPath)
	if err != nil {
		fmt.Println(err)
	}

	// Build json data from map.
	cjsonData, err := json.Marshal(fileMap)
	if err != nil {
		fmt.Println(err)
	}
	cosRes, err := cosineContainer.FullUse(cjsonData)
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}

	// Print Result.
	fmt.Printf("Cosine Result: %v\n\n", cosRes)

	// Configure Sudoku container.
	var sudokuContainer dataduit.DdContainer

	sInputPath := "./input/input_sudoku/puzzle_01.txt"
	err = sudokuContainer.Configure("./config/sudoku_config.yml")
	if err != nil {
		fmt.Printf("Error with sudokuContainer config: %v\n", err)
	}

	sfileMap, err := createSudokuInput(sInputPath)
	if err != nil {
		fmt.Println(err)
	}

	// Build json data from map.
	sjsonData, err := json.Marshal(sfileMap)
	if err != nil {
		fmt.Println(err)
	}

	// Use sudoku container.
	sudokuRes, err := sudokuContainer.FullUse(sjsonData)
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}
	fmt.Printf("Sudoku Result: %v\n", sudokuRes)

}

// ------------------ HELPER FUNCTIONS USED BY THE ABOVE WRAPPER FUNCTIONS [END]

// ---------------------------------------------- INPUT PREPROCESSING FUNCTIONS
// NOTE: This will likely fall on the user to create the specified input data.
// and will only be included as an example that the user can pull from.

// createMap is a helper that accepts a path to a directory and creates the
// input data for the model.
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

func createSudokuInput(fPath string) (map[string][]string, error) {
	sudokuMap := make(map[string][]string)

	data, err := ioutil.ReadFile(fPath)
	if err != nil {
		return sudokuMap, err
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
			return sudokuMap, fmt.Errorf("unexpected value (%v) in Sudoku input", c)
		}
	}

	return sudokuMap, nil
}

//------------------------------------------------ Sudoku [END]
// ----------------------------------------- INPUT PREPROCESSING FUNCTIONS [END]
