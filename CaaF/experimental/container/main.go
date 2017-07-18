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
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)
	mux.HandleFunc("/cosineSim", cosineSim)

	http.ListenAndServe(":8080", mux)
}

// handler is the main handler and returns the current time.
// NOTE: This is included for demo purposes only.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "jack")
	curTime := time.Now()
	fmt.Fprintf(w, "%s", curTime)
}

func cosineSim(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to be looked at more closely
	// - does the default case work?
	// - if I'm gziping everything, do we need a default case?
	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(r.Body)
		defer r.Body.Close()
	default:
		reader = r.Body
	}

	// Read in data from gzip reader.
	JSONData, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("unable to read data: %v", err)
	}

	// Call main implementation function.
	fNameToCosSim, err := CosineSimilarity(JSONData)
	if err != nil {
		fmt.Fprintf(w, "Unable to calculate cosineSimilarity: %v", err)
	}
	fmt.Fprintf(w, "Success: Cosine map: %v", fNameToCosSim)

}
