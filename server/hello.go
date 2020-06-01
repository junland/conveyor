package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
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

	//queue.GenQScript(uint64(exectime), c.WorkerPool, newJob.Commands)

	c.WorkerPool.RunScript(queue.JobScriptCmd{
		JobId:     uint64(exectime),
		Cmds:      newJob.Commands,
		ScriptDir: c.JobDir,
	})

	respondJSON(w, http.StatusOK, map[string]string{"message": "Job Submitted"})
	if err != nil {
		log.Error("Something went wrong with submitting work to queue: ", err)
		return
	}

	return
}

// StopJob function will stop the specified job.
func (c *Config) StopJob(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	i, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		log.Error("Something went wrong with converting the string: ", err)
		respondError(w, http.StatusInternalServerError, "Something went wrong.")
		return
	}

	c.WorkerPool.StopJob(i)

	respondJSON(w, http.StatusOK, map[string]string{"message": "Issueing stop to job."})
	if err != nil {
		log.Error("Something went wrong with submitting stop command: ", err)
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
