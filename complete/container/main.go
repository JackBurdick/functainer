package main

/*
NOTE:
 - Could this be the "dataduit" wrapper?
 - Send this through a load balancer to call container w/compiled function?

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
	log.SetOutput(os.Stdout)
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)
	mux.HandleFunc("/cosineSim", cosineSim)

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

	// jsonData, err := json.Marshal(fNameToCosSim)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	fmt.Fprint(w, fNameToCosSim)
	//fmt.Fprint(w, jsonData)

}
