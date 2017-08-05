package main

import (
	"fmt"

	"github.com/JackBurdick/dataduit/CaaF/dataduit"
)

func main() {
	var err error

	// Configure cosine container.
	var cosineContainer dataduit.DdContainer
	cosineConfig := "./config/cosine_config.yml"
	cInputPath := "./input/input_cosine/"
	err = cosineContainer.ConfigDD(cosineConfig)
	if err != nil {
		fmt.Printf("Error with cosineContainer config: %v\n", err)
	}

	// Use cosine container.
	cosRes, err := cosineContainer.CompleteDD(cInputPath, "cosine")
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}
	fmt.Printf("Cosine Result: %v\n\n", cosRes)

	// Configure Sudoku container.
	var sudokuContainer dataduit.DdContainer
	sudokuConfig := "./config/sudoku_config.yml"
	sInputPath := "./input/input_sudoku/puzzle_01.txt"
	err = sudokuContainer.ConfigDD(sudokuConfig)
	if err != nil {
		fmt.Printf("Error with sudokuContainer config: %v\n", err)
	}

	// Use sudoku container.
	sudokuRes, err := sudokuContainer.CompleteDD(sInputPath, "sudoku")
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}
	fmt.Printf("Sudoku Result: %v\n", sudokuRes)

}
