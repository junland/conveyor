package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/junland/conveyor/queue"
	log "github.com/sirupsen/logrus"
)

// JobRequest describes the statement of work.
type JobRequest struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`
}

// CreateJob is a function that collects and parses incoming jobs.
func (c *Config) CreateJob(w http.ResponseWriter, r *http.Request) {

	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Something went wrong with parsing the json request: ", err)
		respondError(w, http.StatusBadRequest, "Could not parse json.")
		return
	}

	var newJob JobRequest

	json.Unmarshal(reqBody, &newJob)

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if newJob.Name == "" {
		respondError(w, http.StatusBadRequest, "No job name specified.")
		return
	}

	exectime := time.Now().UnixNano() / 1e6

	rand.Seed(exectime)

	log.Info("Queueing up job")

	work := queue.Job{ID: uint64(exectime), CmdList: newJob.Commands}

	queue.JobQueue <- work

	respondJSON(w, http.StatusOK, map[string]string{"message": "Job Submitted"})
	if err != nil {
		log.Error("Something went wrong with submitting work to queue: ", err)
		return
	}

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
