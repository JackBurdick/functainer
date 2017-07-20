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

func main() {
	// TODO: logging should be included
	log.SetOutput(os.Stdout)
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)

	// TODO: The name of this endpoint should be generated by the config file
	// and should match the call in the file that calls the container-right now
	// declared as a const in `../use_container/main.go`
	mux.HandleFunc("/cosineSim", cosineSim)

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
func cosineSim(w http.ResponseWriter, r *http.Request) {

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
	fNameToCosSim, err := CosineSimilarity(JSONData)
	if err != nil {
		log.Printf("Unable to calculate cosineSimilarity: %v", err)
	}

	fmt.Fprint(w, fNameToCosSim)

	// jsonData, err := json.Marshal(fNameToCosSim)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Fprint(w, jsonData)
}
