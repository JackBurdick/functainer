package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jdkato/prose/tokenize"
)

func main() {
	inputDir := "./input/"
	URL := "localhost:8000/cosineSim"

	// Create map from input directory.
	fileMap, err := createMap(inputDir)
	if err != nil {
		fmt.Println(err)
	}

	// Build form from map.
	jsonData, err := json.Marshal(fileMap)
	if err != nil {
		fmt.Println(err)
	}

	// Make http request to container. (using JSON)
	result, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	// Process results.
	fmt.Println(result)

}

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
