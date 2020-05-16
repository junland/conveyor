package server

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"

	"github.com/julienschmidt/httprouter"
)

// JobRequest describes the statement of work.
type JobRequest []struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`
}

// Collector is a fucntion that collects and parses incoming jobs.
func (c *Config) Collector(w http.ResponseWriter, r *http.Request) {
	var newJobRequest JobRequest

	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
    
    reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)

	// Now, we retrieve the person's name from the request.
	name := r.FormValue("name")

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	fmt.Println("Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)

	return
}

// helloRootHandle is a handle.
func helloRootHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(202)
	w.Write([]byte("I am root."))
}

// helloGlobalHandle is a example handler.
func helloGlobalHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello everyone!"))
}

// helloNameHandle is a example parameter handler.
func helloNameHandle(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())
	name := ps.ByName("name")
	w.WriteHeader(201)
	fmt.Fprintf(w, "Hello, %s\n", name)
}
