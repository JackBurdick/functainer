package main

import (
	"fmt"
	"io/ioutil"

	"github.com/jdkato/prose/tokenize"
)

func main() {

	// get input ready

	// make http request to container

	// process results

}

func createMap(dPath string) (map[string][]string, error) {
	fileMap := make(map[string][]string)
	// Create map of filename to tokenized content.
	dFiles, _ := ioutil.ReadDir(dPath)
	for _, f := range dFiles {
		b, err := ioutil.ReadFile(dPath + f.Name())
		if err != nil {
			fmt.Print(err)
		}

		// Convert bytes to string, then use 3rd party to tokenize.
		fileMap[f.Name()] = tokenize.TextToWords(string(b))
	}

	return fileMap, nil
}
