package main

/*
NOTE:
 - Could this be the "dataduit" wrapper?
 - maybe

*/

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// crossIndex 'crosses' two strings such that the two individual values from
// each string join together to create a new value.  For example, if string one
// is "ABC" and string two is "123", the resulting return value will be;
// ["A1","A2","A3","B1","B2","B3","C1","C2","C3"].
// func crossIndex(A string, N string) []string {
// 	var ks []string
// 	for _, a := range A {
// 		for _, n := range N {
// 			ks = append(ks, (string(a) + string(n)))
// 		}
// 	}
// 	return ks
// }

func main() {
	// TODO: logging should be included
	log.SetOutput(os.Stdout)
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)

	// TODO: The name of this endpoint should be generated by the config file
	// and should match the call in the file that calls the container-right now
	// declared as a const in `../use_container/main.go`
	mux.HandleFunc("/sudoku", sudoku)

	// TODO: The port of this endpoint should be generated by the config file
	// and should match the call in the file that calls the container-right now
	// declared as a const in `../use_container/main.go`
	http.ListenAndServe(":8080", mux)
}

// handler is the main handler and returns the current time.
// NOTE: This is included for demo purposes only.
func handler(w http.ResponseWriter, r *http.Request) {
	curTime := time.Now()
	fmt.Fprintf(w, "%s", curTime)
}

// cosineSim is the wrapper for the cosineSimilarity functionality.
// Format input from {raw} -> {expected}
func sudoku(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to be looked at more closely
	// - does the default case work?
	// - if I'm gziping everything, do we need a default case?
	var reader io.ReadCloser
	var err error
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			log.Printf("Error creating new gzip reader%v\n", err)
		}
		defer r.Body.Close()
	default:
		reader = r.Body
	}

	// Read in data from gzip reader.
	JSONData, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("unable to read data gzip data: %v\n", err)
	}

	// Call main implementation function passing the data as a JSON object.
	ssResults, err := solveSudoku(JSONData)
	if err != nil {
		log.Printf("Unable to calculate cosineSimilarity: %v", err)
	}

	fmt.Fprint(w, ssResults)

}