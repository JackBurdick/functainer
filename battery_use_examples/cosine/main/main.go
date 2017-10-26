package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/JackBurdick/dataduit/dataduit"
	"github.com/jdkato/prose/tokenize"
)

func main() {
	var cosineContainer dataduit.DdContainer

	// Build battery according to specification.
	err := cosineContainer.Configure("./cosine_config.yml")
	if err != nil {
		fmt.Printf("Error with cosineContainer config: %v\n", err)
	}

	// Create input for cosine battery
	cjsonInput, err := createMap("../input/")
	if err != nil {
		fmt.Println(err)
	}

	// Use cosine battery.
	cosRes, err := cosineContainer.FullUse(cjsonInput)
	if err != nil {
		fmt.Printf("Error using container: %v\n", err)
	}

	// Print Result.
	fmt.Printf("Cosine Result: %v\n\n", cosRes)

}

// INPUT PREPROCESSING FUNCTIONS
// NOTE: This will likely fall on the user to create the specified input data.
// and will only be included as an example that the user can pull from.

// createMap is a helper that accepts a path to a directory and creates the
// input data for the model.
func createMap(dPath string) ([]byte, error) {
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

	// Build json data from map.
	cjsonData, err := json.Marshal(fileMap)
	if err != nil {
		fmt.Println(err)
	}

	return cjsonData, nil
}
