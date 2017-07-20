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
	log.Printf("inside handler\n")
	fmt.Fprintf(w, "%s", "jack")
	curTime := time.Now()
	fmt.Fprintf(w, "%s", curTime)
}

func cosineSim(w http.ResponseWriter, r *http.Request) {
	log.Printf("inside cosineSim\n")

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
	log.Printf("after gzip\n")

	// Read in data from gzip reader.
	JSONData, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("unable to read data gzip data: %v\n", err)
		fmt.Printf("unable to read data: %v", err)
		fmt.Fprintf(w, "unable to read data: %v", err)
	}
	log.Printf("after read JSON data\n")

	// Call main implementation function.
	fNameToCosSim, err := CosineSimilarity(JSONData)
	if err != nil {
		log.Printf("Unable to calculate cosineSimilarity: %v", err)
		fmt.Fprintf(w, "Unable to calculate cosineSimilarity: %v", err)
	}
	fmt.Fprintf(w, "Success: Cosine map: %v", fNameToCosSim)
	log.Printf("Success: Cosine map: %v", fNameToCosSim)

}
