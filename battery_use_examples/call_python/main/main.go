package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/JackBurdick/dataduit/dataduit"
)

func main() {
	var tfModel dataduit.DdContainer

	// Configure Sudoku container.
	err := tfModel.Configure("./tf_config.yml")
	if err != nil {
		fmt.Printf("Error with tf config: %v\n", err)
	}

	// Preprocess input (according to required specified input type)
	sjsonData, err := create_tf_input("../input/var_01.txt")
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

func create_tf_input(fPath string) ([]byte, error) {
	tfMap := make(map[int][]string)

	data, err := ioutil.ReadFile(fPath)
	if err != nil {
		// Build json data from map.
		jsonData, err := json.Marshal(tfMap)
		if err != nil {
			fmt.Println(err)
		}
		return data, err
	}

	for i, c := range data {
		tfMap[i] = []string{string(c)}
	}

	// Build json data from map.
	jsonData, err := json.Marshal(tfMap)
	if err != nil {
		fmt.Println(err)
	}

	return jsonData, nil
}
