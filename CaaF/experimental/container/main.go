package main

/*
NOTE:
 - Could this be the "dataduit" wrapper?
 - Send this through a load balancer to call container w/compiled function?

*/

import (
	"fmt"
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
	curTime := time.Now()
	fmt.Fprintf(w, "%s", curTime)
}

func cosineSim(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form["data"]) > 0 {

		// Extract data from form and process data
		data := r.Form["data"][0]
		fNameToCosSim, err := cosineSimilarity(data)
		if err != nil {
			fmt.Fprintf(w, "Unable to calculate cosineSimilarity: %v", err)
		}
		fmt.Fprintf(w, "Success: Cosine map: %v", fNameToCosSim)

	} else {
		fmt.Fprintln(w, "Nothing to see here - did you specify the data?")
	}
}
