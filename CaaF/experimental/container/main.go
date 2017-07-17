package main

/*
NOTE:
 - Could this be the "dataduit" wrapper?
 - Send this through a load balancer to call container w/compiled function?

*/

import (
	"encoding/json"
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

	decoder := json.NewDecoder(r.Body)
	var JSONdata []byte
	err := decoder.Decode(&JSONdata)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	fNameToCosSim, err := cosineSimilarity(JSONdata)
	if err != nil {
		fmt.Fprintf(w, "Unable to calculate cosineSimilarity: %v", err)
	}
	fmt.Fprintf(w, "Success: Cosine map: %v", fNameToCosSim)

}
