package main

import (
	"fmt"

	"github.com/JackBurdick/dataduit/dataduit"
)

func main() {
	var err error

	// Initialize cosine battery.
	var cosineContainer dataduit.DdContainer

	// Build battery according to specification.
	err = cosineContainer.Build("./config/cosine_config.yml")
	if err != nil {
		fmt.Printf("Error with cosineContainer config: %v\n", err)
	}

	// Use cosine battery.
	cInputPath := "./input/input_cosine/"
	cosRes, err := cosineContainer.CompleteDD(cInputPath, "cosine")
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}

	// Print Result.
	fmt.Printf("Cosine Result: %v\n\n", cosRes)

	// Configure Sudoku container.
	var sudokuContainer dataduit.DdContainer

	sInputPath := "./input/input_sudoku/puzzle_01.txt"
	err = sudokuContainer.Build("./config/sudoku_config.yml")
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
